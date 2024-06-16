//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package types

type Dialect int

const (
	DialectUnknown = iota
	DialectMySQL
	DialectPostgreSQL
	DialectTSQL
	DialectSQLite
)
