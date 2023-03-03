//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

type SQLType int

const (
	SQL_TYPE_CHAR SQLType = iota
	SQL_TYPE_INT
	SQL_TYPE_FLOAT
	SQL_TYPE_DECIMAL
	SQL_TYPE_VARCHAR
	SQL_TYPE_TEXT
	SQL_TYPE_BINARY
)

var (
	sqlTypeToSymbol = map[SQLType]Symbol{
		SQL_TYPE_CHAR:    SYM_TYPE_CHAR,
		SQL_TYPE_VARCHAR: SYM_TYPE_VARCHAR,
		SQL_TYPE_INT:     SYM_TYPE_INT,
		SQL_TYPE_BINARY:  SYM_TYPE_BINARY,
		SQL_TYPE_FLOAT:   SYM_TYPE_FLOAT,
		SQL_TYPE_DECIMAL: SYM_TYPE_DECIMAL,
		SQL_TYPE_TEXT:    SYM_TYPE_TEXT,
	}
)

func SQLTypeToSymbol(st SQLType) Symbol {
	return sqlTypeToSymbol[st]
}
