//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doUpdateStatementSearched(
	el *grammar.UpdateStatementSearched,
	qargs []interface{},
	curarg *int,
) {
	b.Write(grammar.Symbols[grammar.SYM_UPDATE])
	// We don't add any table alias when outputting the table identifier
	b.WriteString(el.TableName)
	b.Write(grammar.Symbols[grammar.SYM_SET])

	for x, c := range el.Columns {
		if x > 0 {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
		// We don't add the table identifier or use an alias when outputting
		// the column names in the <column_value_lists> element of the UPDATE
		// statement
		b.WriteString(c)
		b.Write(grammar.Symbols[grammar.SYM_EQUAL])
		b.WriteString(InterpolationMarker(b.opts, *curarg))
		qargs[*curarg] = el.Values[x]
		*curarg++
	}

	if el.WhereClause != nil {
		b.doWhereClause(el.WhereClause, qargs, curarg)
	}
}
