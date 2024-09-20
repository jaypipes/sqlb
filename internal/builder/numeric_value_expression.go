//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
)

func (b *Builder) doNumericValueExpression(
	el *grammar.NumericValueExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.Unary != nil {
		b.doTerm(el.Unary, qargs, curarg)
	}
}

func (b *Builder) doTerm(
	el *grammar.Term,
	qargs []interface{},
	curarg *int,
) {
	if el.Unary != nil {
		b.doFactor(el.Unary, qargs, curarg)
	}
}

func (b *Builder) doFactor(
	el *grammar.Factor,
	qargs []interface{},
	curarg *int,
) {
	if el.Sign != grammar.SignPlus {
		b.WriteString(grammar.SignSymbol[el.Sign])
	}
	b.doNumericPrimary(&el.Primary, qargs, curarg)
}

func (b *Builder) doNumericPrimary(
	el *grammar.NumericPrimary,
	qargs []interface{},
	curarg *int,
) {
	if el.Primary != nil {
		b.doValueExpressionPrimary(el.Primary, qargs, curarg)
	} else if el.Function != nil {
		b.doNumericValueFunction(el.Function, qargs, curarg)
	}
}
