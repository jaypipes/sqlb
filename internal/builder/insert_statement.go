//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
)

func (b *Builder) doInsertStatement(
	el *grammar.InsertStatement,
	qargs []interface{},
	curarg *int,
) {
	b.Write(grammar.Symbols[grammar.SYM_INSERT])
	// We don't add any table alias when outputting the table identifier
	b.WriteString(el.TableName)
	b.WriteRune(' ')
	b.Write(grammar.Symbols[grammar.SYM_LPAREN])

	for x, c := range el.Columns {
		if x > 0 {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
		// We don't add the table identifier or use an alias when outputting
		// the column names in the <columns> element of the INSERT statement
		b.WriteString(c)
	}
	b.Write(grammar.Symbols[grammar.SYM_VALUES])
	for x, v := range el.Values {
		if x > 0 {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
		b.WriteString(InterpolationMarker(b.opts, *curarg))
		qargs[*curarg] = v
		*curarg++
	}
	b.Write(grammar.Symbols[grammar.SYM_RPAREN])
}
