//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package meta

import (
	"database/sql"

	"github.com/jaypipes/sqlb/types"
)

const (
	selDBNameMySQL      = "SELECT DATABASE()"
	selDBNamePostgreSQL = "SELECT CURRENT_DATABASE()"
)

// DatabaseName returns the database schema name given a sql.DB handle
func DatabaseName(
	db *sql.DB,
	mods ...MetaOptionModifier,
) string {
	opts := mergeOpts(mods)
	if opts.Dialect == types.DialectUnknown {
		opts.Dialect = Dialect(db)
	}
	var qs string
	switch opts.Dialect {
	case types.DialectMySQL:
		qs = selDBNameMySQL
	case types.DialectPostgreSQL:
		qs = selDBNamePostgreSQL
	}
	var dbName string
	if err := db.QueryRow(qs).Scan(&dbName); err != nil {
		return ""
	}
	return dbName
}
