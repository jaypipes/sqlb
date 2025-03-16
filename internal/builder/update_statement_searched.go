//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
	"github.com/jaypipes/sqlb/core/grammar/symbol"
)

func (b *Builder) doUpdateStatementSearched(
	el *grammar.UpdateStatementSearched,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.Update)
	b.WriteString(symbol.Space)
	// We don't add any table alias when outputting the table identifier
	b.WriteString(el.TableName)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.Set)
	b.WriteString(symbol.Space)

	for x, c := range el.Columns {
		if x > 0 {
			b.WriteString(symbol.Comma)
			b.WriteString(symbol.Space)
		}
		// We don't add the table identifier or use an alias when outputting
		// the column names in the <column_value_lists> element of the UPDATE
		// statement
		b.WriteString(c)
		b.WriteString(symbol.Space)
		b.WriteString(symbol.EqualsOperator)
		b.WriteString(symbol.Space)
		b.WriteString(InterpolationMarker(b.opts, *curarg))
		qargs[*curarg] = el.Values[x]
		*curarg++
	}

	if el.Where != nil {
		b.doWhereClause(el.Where, qargs, curarg)
	}
}
