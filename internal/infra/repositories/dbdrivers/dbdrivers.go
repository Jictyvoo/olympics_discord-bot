// Package dbdrivers registers the supported database/sql drivers as an import
// side effect, so commands blank-import it instead of each driver.
package dbdrivers

import (
	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)
