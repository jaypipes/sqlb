//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
)

func (b *Builder) doRowValuePredicand(
	el *grammar.RowValuePredicand,
	qargs []interface{},
	curarg *int,
) {
	if el.Primary != nil {
		b.doNonParenthesizedValueExpressionPrimary(el.Primary, qargs, curarg)
	} else if el.Common != nil {
		b.doCommonValueExpression(el.Common, qargs, curarg)
	} else if el.Boolean != nil {
		b.doBooleanPredicand(el.Boolean, qargs, curarg)
	}
}
