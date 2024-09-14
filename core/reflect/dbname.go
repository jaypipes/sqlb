//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package reflect

import (
	"database/sql"

	"github.com/jaypipes/sqlb/core/types"
)

const (
	selDBNameMySQL      = "SELECT DATABASE()"
	selDBNamePostgreSQL = "SELECT CURRENT_DATABASE()"
)

// DatabaseName returns the database schema name given a sql.DB handle
func DatabaseName(
	db *sql.DB,
	mods ...types.Option,
) string {
	opts := types.MergeOptions(mods)
	var d types.Dialect
	if !opts.HasDialect() {
		d = Dialect(db)
	} else {
		d = opts.Dialect()
	}
	var qs string
	switch d {
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
