package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

type MigrationError struct {
	Version string
	Err     error
}

func (e *MigrationError) Error() string {
	return fmt.Sprintf("migrator: apply %s: %v", e.Version, e.Err)
}

func (e *MigrationError) Unwrap() error { return e.Err }

// Migrator applies embedded SQL migrations. Down migrations are unsupported;
// rollbacks are new forward migrations.
type Migrator struct {
	db     *sql.DB
	driver string
	fs     fs.FS
}

// New builds a Migrator; driver is the sql driver name (e.g. "sqlite3", "mysql").
func New(db *sql.DB, driver string, migrations fs.FS) *Migrator {
	return &Migrator{db: db, driver: driver, fs: migrations}
}

// Run applies pending migrations in lexicographic order. Safe to call multiple
// times; already-applied versions are skipped.
func (m *Migrator) Run(ctx context.Context) error {
	if err := m.bootstrap(ctx); err != nil {
		return fmt.Errorf("migrator: bootstrap: %w", err)
	}

	applied, err := m.appliedVersions(ctx)
	if err != nil {
		return fmt.Errorf("migrator: list applied: %w", err)
	}

	files, err := m.sqlFiles()
	if err != nil {
		return fmt.Errorf("migrator: list files: %w", err)
	}

	for _, name := range files {
		version := strings.TrimSuffix(name, ".sql")
		if applied[version] {
			continue
		}
		if err = m.apply(ctx, name, version); err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) bootstrap(ctx context.Context) error {
	_, err := m.db.ExecContext(ctx, bootstrapSQL(m.driver))
	return err
}

func (m *Migrator) appliedVersions(ctx context.Context) (map[string]bool, error) {
	rows, err := m.db.QueryContext(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out := make(map[string]bool)
	for rows.Next() {
		var v string
		if err = rows.Scan(&v); err != nil {
			return nil, err
		}
		out[v] = true
	}
	return out, rows.Err()
}

func (m *Migrator) sqlFiles() ([]string, error) {
	entries, err := fs.ReadDir(m.fs, ".")
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

func (m *Migrator) apply(ctx context.Context, filename, version string) error {
	raw, err := fs.ReadFile(m.fs, filename)
	if err != nil {
		return &MigrationError{Version: version, Err: err}
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return &MigrationError{Version: version, Err: err}
	}

	if _, err = tx.ExecContext(ctx, string(raw)); err != nil {
		_ = tx.Rollback()
		return &MigrationError{Version: version, Err: err}
	}

	_, err = tx.ExecContext(ctx,
		"INSERT INTO schema_migrations (version) VALUES (?)", version)
	if err != nil {
		_ = tx.Rollback()
		return &MigrationError{Version: version, Err: err}
	}

	if err = tx.Commit(); err != nil {
		return &MigrationError{Version: version, Err: err}
	}
	return nil
}
