//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doValueExpression(
	el *grammar.ValueExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.CommonValueExpression != nil {
		b.doCommonValueExpression(el.CommonValueExpression, qargs, curarg)
	} else if el.BooleanValueExpression != nil {
		b.doBooleanValueExpression(el.BooleanValueExpression, qargs, curarg)
	} else if el.RowValueExpression != nil {
		b.doNonParenthesizedValueExpressionPrimary(el.RowValueExpression.NonParenthesizedValueExpressionPrimary, qargs, curarg)
	}
}

func (b *Builder) doCommonValueExpression(
	el *grammar.CommonValueExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.NumericValueExpression != nil {
		b.doNumericValueExpression(el.NumericValueExpression, qargs, curarg)
	} else if el.StringValueExpression != nil {
		b.doStringValueExpression(el.StringValueExpression, qargs, curarg)
	} else if el.DatetimeValueExpression != nil {
		b.doDatetimeValueExpression(el.DatetimeValueExpression, qargs, curarg)
	} else if el.IntervalValueExpression != nil {
		b.doIntervalValueExpression(el.IntervalValueExpression, qargs, curarg)
	}
}
