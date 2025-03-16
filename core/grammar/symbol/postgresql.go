//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package symbol

// Special character symbols for PostgreSQL variants
const (
	SymbolPostgreSQLSpecialCharacterStart Symbol = 30000
	SymbolDollar
	SymbolPostgreSQLSpecialCharacterEnd = 30200
)

const (
	Dollar = "$"
)

// Reserved words in lexicographical order
const (
	SymbolPostgreSQLReservedStart Symbol = SymbolPostgreSQLSpecialCharacterEnd + 1
	SymbolLimit
)

const (
	Limit  = "LIMIT"
	Offset = "OFFSET"
)
