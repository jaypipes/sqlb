//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doDatetimeValueExpression(
	el *grammar.DatetimeValueExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.Unary != nil {
		b.doDatetimeTerm(el.Unary, qargs, curarg)
	}
}

func (b *Builder) doDatetimeTerm(
	el *grammar.DatetimeTerm,
	qargs []interface{},
	curarg *int,
) {
	b.doDatetimeFactor(&el.Factor, qargs, curarg)
}

func (b *Builder) doDatetimeFactor(
	el *grammar.DatetimeFactor,
	qargs []interface{},
	curarg *int,
) {
	b.doDatetimePrimary(&el.Primary, qargs, curarg)
}

func (b *Builder) doDatetimePrimary(
	el *grammar.DatetimePrimary,
	qargs []interface{},
	curarg *int,
) {
	if el.Primary != nil {
		b.doValueExpressionPrimary(el.Primary, qargs, curarg)
	} else if el.Function != nil {
		b.doDatetimeValueFunction(el.Function, qargs, curarg)
	}
}
