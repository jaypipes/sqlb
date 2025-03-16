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

func (b *Builder) doDerivedColumn(
	el *grammar.DerivedColumn,
	qargs []interface{},
	curarg *int,
) {
	b.doValueExpression(&el.Value, qargs, curarg)
	if el.As != nil {
		b.WriteString(symbol.Space)
		b.WriteString(symbol.As)
		b.WriteString(symbol.Space)
		b.WriteString(*el.As)
	}
}
