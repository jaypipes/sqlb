//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

// Dialect is the SQL variant of the underlying RDBMS
type Dialect int

const (
	DialectUnknown Dialect = iota
	DialectMySQL
	DialectPostgreSQL
	DialectTSQL
	DialectSQLite
)
