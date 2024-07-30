package bootstrap

import (
	"database/sql"
	"log/slog"
	"os"
)

func OpenDatabase(filename string) *sql.DB {
	db, dbErr := sql.Open("sqlite", filename)
	if dbErr != nil {
		slog.Error("failed to open database", slog.String("error", dbErr.Error()))
		os.Exit(1)
	}

	return db
}
