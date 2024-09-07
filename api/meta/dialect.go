//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package meta

import (
	"database/sql"
	"reflect"

	"github.com/jaypipes/sqlb/api"
)

var driverNameToDialect = map[string]api.Dialect{
	"*mssql.MssqlDriver": api.DialectTSQL,
	"*pq.Driver":         api.DialectPostgreSQL,
	"*stdlib.Driver":     api.DialectPostgreSQL,
}

// Dialect returns the SQL Dialect after examining the supplied database
// connection
func Dialect(
	db *sql.DB,
) api.Dialect {
	drv := db.Driver()
	dv := reflect.ValueOf(drv)
	d, found := driverNameToDialect[dv.Type().String()]
	if !found {
		return api.DialectUnknown
	}
	return d
}
