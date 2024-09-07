//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doValueExpressionPrimary(
	el *grammar.ValueExpressionPrimary,
	qargs []interface{},
	curarg *int,
) {
	if el.ParenthesizedValueExpression != nil {
		b.Write(grammar.Symbols[grammar.SYM_LPAREN])
		b.doValueExpression(el.ParenthesizedValueExpression, qargs, curarg)
		b.Write(grammar.Symbols[grammar.SYM_RPAREN])
	} else if el.NonParenthesizedValueExpressionPrimary != nil {
		b.doNonParenthesizedValueExpressionPrimary(el.NonParenthesizedValueExpressionPrimary, qargs, curarg)
	}
}

func (b *Builder) doNonParenthesizedValueExpressionPrimary(
	el *grammar.NonParenthesizedValueExpressionPrimary,
	qargs []interface{},
	curarg *int,
) {
	if el.UnsignedValueSpecification != nil {
		b.doUnsignedValueSpecification(el.UnsignedValueSpecification, qargs, curarg)
	} else if el.ColumnReference != nil {
		b.doColumnReference(el.ColumnReference, qargs, curarg)
	} else if el.SetFunctionSpecification != nil {
		b.doSetFunctionSpecification(el.SetFunctionSpecification, qargs, curarg)
	} else if el.ScalarSubquery != nil {
		b.doSubquery(el.ScalarSubquery, qargs, curarg)
	}
}
