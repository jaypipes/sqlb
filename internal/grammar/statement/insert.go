//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package statement

import (
	"strings"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/identifier"
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

func (st *Insert) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
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
		b.WriteString(builder.InterpolationMarker(opts, *curarg))
		qargs[*curarg] = v
		*curarg++
		if x != (ncols - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
	b.Write(grammar.Symbols[grammar.SYM_RPAREN])
	return b.String()
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
