package migrator

func bootstrapSQL(driver string) string {
	switch driver {
	case "mysql":
		return `CREATE TABLE IF NOT EXISTS schema_migrations (
			version     VARCHAR(255) NOT NULL PRIMARY KEY,
			applied_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	default: // sqlite, sqlite3, modernc
		return `CREATE TABLE IF NOT EXISTS schema_migrations (
			version    TEXT    NOT NULL PRIMARY KEY,
			applied_at DATETIME NOT NULL DEFAULT (DATETIME('now'))
		);`
	}
}
