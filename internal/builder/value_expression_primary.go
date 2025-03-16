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

func (b *Builder) doValueExpressionPrimary(
	el *grammar.ValueExpressionPrimary,
	qargs []interface{},
	curarg *int,
) {
	if el.Parenthesized != nil {
		b.WriteString(symbol.LeftParen)
		b.doValueExpression(el.Parenthesized, qargs, curarg)
		b.WriteString(symbol.RightParen)
	} else if el.Primary != nil {
		b.doNonParenthesizedValueExpressionPrimary(el.Primary, qargs, curarg)
	}
}

func (b *Builder) doNonParenthesizedValueExpressionPrimary(
	el *grammar.NonParenthesizedValueExpressionPrimary,
	qargs []interface{},
	curarg *int,
) {
	if el.UnsignedValue != nil {
		b.doUnsignedValueSpecification(el.UnsignedValue, qargs, curarg)
	} else if el.ColumnReference != nil {
		b.doColumnReference(el.ColumnReference, qargs, curarg)
	} else if el.SetFunction != nil {
		b.doSetFunctionSpecification(el.SetFunction, qargs, curarg)
	} else if el.ScalarSubquery != nil {
		b.doSubquery(el.ScalarSubquery, qargs, curarg)
	}
}
