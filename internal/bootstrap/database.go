package bootstrap

import (
	"database/sql"
	"log/slog"
	"os"
)

func OpenDatabase() *sql.DB {
	db, dbErr := sql.Open("sqlite", "olympics-2024_PARIS.db")
	if dbErr != nil {
		slog.Error("failed to open database", slog.String("error", dbErr.Error()))
		os.Exit(1)
	}

	return db
}
