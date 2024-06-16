//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package statement

import (
	"strings"

	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/identifier"
	"github.com/jaypipes/sqlb/internal/scanner"
)

// INSERT INTO <table> (<columns>) VALUES (<values>)

type Insert struct {
	table   *identifier.Table
	columns []*identifier.Column
	values  []interface{}
}

func (st *Insert) ArgCount() int {
	return len(st.values)
}

func (st *Insert) Size(s *scanner.Scanner) int {
	size := len(grammar.Symbols[grammar.SYM_INSERT]) + len(st.table.Name) + 1 // space after table name
	ncols := len(st.columns)
	for _, c := range st.columns {
		// We don't add the table identifier or use an alias when outputting
		// the column names in the <columns> element of the INSERT statement
		size += len(c.Name)
	}
	// We don't include interpolation marks in our sizing, since the length
	// differs with SQL dialects. This is accounted for by callers of scan().
	size += len(grammar.Symbols[grammar.SYM_LPAREN]) + len(grammar.Symbols[grammar.SYM_VALUES])
	// Two comma-delimited lists of same number of elements (columns and
	// values)
	size += 2 * (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (ncols - 1)) // the commas...
	size += len(grammar.Symbols[grammar.SYM_RPAREN])
	return size
}

func (st *Insert) Scan(s *scanner.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	b.Write(grammar.Symbols[grammar.SYM_INSERT])
	// We don't add any table alias when outputting the table identifier
	b.WriteString(st.table.Name)
	b.WriteRune(' ')
	b.Write(grammar.Symbols[grammar.SYM_LPAREN])

	ncols := len(st.columns)
	for x, c := range st.columns {
		// We don't add the table identifier or use an alias when outputting
		// the column names in the <columns> element of the INSERT statement
		b.WriteString(c.Name)
		if x != (ncols - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
	b.Write(grammar.Symbols[grammar.SYM_VALUES])
	for x, v := range st.values {
		scanner.ScanInterpolationMarker(s.Dialect, b, *curArg)
		args[*curArg] = v
		*curArg++
		if x != (ncols - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
	b.Write(grammar.Symbols[grammar.SYM_RPAREN])
}

// NewInsert returns a new InsertStatement struct that scans into an
// INSERT SQL statement
func NewInsert(
	table *identifier.Table,
	columns []*identifier.Column,
	values []interface{},
) *Insert {
	return &Insert{
		table:   table,
		columns: columns,
		values:  values,
	}
}
