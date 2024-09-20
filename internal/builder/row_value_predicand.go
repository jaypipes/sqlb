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
	if el.NonParenthesizedValueExpressionPrimary != nil {
		b.doNonParenthesizedValueExpressionPrimary(el.NonParenthesizedValueExpressionPrimary, qargs, curarg)
	} else if el.CommonValueExpression != nil {
		b.doCommonValueExpression(el.CommonValueExpression, qargs, curarg)
	} else {
		b.doBooleanPredicand(el.BooleanPredicand, qargs, curarg)
	}
}
