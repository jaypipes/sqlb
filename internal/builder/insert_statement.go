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

func (b *Builder) doInsertStatement(
	el *grammar.InsertStatement,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.Insert)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.Into)
	b.WriteString(symbol.Space)
	// We don't add any table alias when outputting the table identifier
	b.WriteString(el.TableName)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.LeftParen)

	for x, c := range el.Columns {
		if x > 0 {
			b.WriteString(symbol.Comma)
			b.WriteString(symbol.Space)
		}
		// We don't add the table identifier or use an alias when outputting
		// the column names in the <columns> element of the INSERT statement
		b.WriteString(c)
	}
	b.WriteString(symbol.RightParen)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.Values)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.LeftParen)
	for x, v := range el.Values {
		if x > 0 {
			b.WriteString(symbol.Comma)
			b.WriteString(symbol.Space)
		}
		b.WriteString(InterpolationMarker(b.opts, *curarg))
		qargs[*curarg] = v
		*curarg++
	}
	b.WriteString(symbol.RightParen)
}
