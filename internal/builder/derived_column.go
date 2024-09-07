//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doDerivedColumn(
	el *grammar.DerivedColumn,
	qargs []interface{},
	curarg *int,
) {
	b.doValueExpression(&el.ValueExpression, qargs, curarg)
	if el.As != nil {
		b.Write(grammar.Symbols[grammar.SYM_AS])
		b.WriteString(*el.As)
	}
}
