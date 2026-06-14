package migrator

import (
	"context"
	"database/sql"
	"io/fs"
	"testing"
	"testing/fstest"

	_ "modernc.org/sqlite"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open in-memory sqlite: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func newTestFS(files map[string]string) fs.FS {
	m := make(fstest.MapFS)
	for name, content := range files {
		m[name] = &fstest.MapFile{Data: []byte(content)}
	}
	return m
}

func TestMigrator_AppliesInOrder(t *testing.T) {
	migrations := newTestFS(map[string]string{
		"001_create_foo.sql": "CREATE TABLE foo (id INTEGER PRIMARY KEY);",
		"002_create_bar.sql": "CREATE TABLE bar (id INTEGER PRIMARY KEY, foo_id INTEGER);",
	})
	db := newTestDB(t)
	m := New(db, "sqlite", migrations)

	if err := m.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}

	ctx := t.Context()
	var count int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations").
		Scan(&count); err != nil {
		t.Fatalf("query schema_migrations: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 applied migrations, got %d", count)
	}

	for _, table := range []string{"foo", "bar"} {
		if _, err := db.ExecContext(ctx, "SELECT 1 FROM "+table+" LIMIT 1"); err != nil {
			t.Fatalf("table %q not created: %v", table, err)
		}
	}
}

func TestMigrator_SkipsAlreadyApplied(t *testing.T) {
	migrations := newTestFS(map[string]string{
		"001_create_foo.sql": "CREATE TABLE foo (id INTEGER PRIMARY KEY);",
	})
	db := newTestDB(t)
	m := New(db, "sqlite", migrations)

	if err := m.Run(context.Background()); err != nil {
		t.Fatalf("first Run: %v", err)
	}
	// Second run must be a no-op; re-running CREATE TABLE would fail.
	if err := m.Run(context.Background()); err != nil {
		t.Fatalf("second Run: %v", err)
	}
}

func TestMigrator_ReturnsErrorOnBadSQL(t *testing.T) {
	migrations := newTestFS(map[string]string{
		"001_bad.sql": "THIS IS NOT SQL;",
	})
	db := newTestDB(t)
	m := New(db, "sqlite", migrations)

	err := m.Run(context.Background())
	if err == nil {
		t.Fatal("expected error for bad SQL, got nil")
	}
	var me *MigrationError
	if !errorAs(err, &me) {
		t.Fatalf("expected *MigrationError, got %T", err)
	}
	if me.Version != "001_bad" {
		t.Fatalf("version = %q, want %q", me.Version, "001_bad")
	}
}

// errorAs is errors.As without importing errors (avoids cycle in test).
func errorAs(err error, target **MigrationError) bool {
	for err != nil {
		if e, ok := err.(*MigrationError); ok {
			*target = e
			return true
		}
		type unwrapper interface{ Unwrap() error }
		if u, ok := err.(unwrapper); ok {
			err = u.Unwrap()
		} else {
			break
		}
	}
	return false
}
