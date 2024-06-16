//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package meta

import (
	"database/sql"
	"reflect"

	"github.com/jaypipes/sqlb/types"
)

var driverNameToDialect = map[string]types.Dialect{
	"*mssql.MssqlDriver": types.DialectTSQL,
	"*pq.Driver":         types.DialectPostgreSQL,
	"*stdlib.Driver":     types.DialectPostgreSQL,
}

// Dialect returns the SQL Dialect after examining the supplied database
// connection
func Dialect(
	db *sql.DB,
) types.Dialect {
	drv := db.Driver()
	dv := reflect.ValueOf(drv)
	d, found := driverNameToDialect[dv.Type().String()]
	if !found {
		return types.DialectUnknown
	}
	return d
}
