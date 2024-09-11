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
	if el.Common != nil {
		b.doCommonValueExpression(el.Common, qargs, curarg)
	} else if el.Boolean != nil {
		b.doBooleanValueExpression(el.Boolean, qargs, curarg)
	} else if el.Row != nil {
		b.doNonParenthesizedValueExpressionPrimary(el.Row.Primary, qargs, curarg)
	}
}

func (b *Builder) doCommonValueExpression(
	el *grammar.CommonValueExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.Numeric != nil {
		b.doNumericValueExpression(el.Numeric, qargs, curarg)
	} else if el.String != nil {
		b.doStringValueExpression(el.String, qargs, curarg)
	} else if el.Datetime != nil {
		b.doDatetimeValueExpression(el.Datetime, qargs, curarg)
	} else if el.Interval != nil {
		b.doIntervalValueExpression(el.Interval, qargs, curarg)
	}
}
