//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
)

func (b *Builder) doBooleanValueExpression(
	el *grammar.BooleanValueExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.Unary != nil {
		b.doBooleanTerm(el.Unary, qargs, curarg)
	} else if el.OrLeft != nil {
		if el.OrRight == nil {
			// This should not happen, so if it does, panic since it's a fault
			// in sqlb's parsing logic.
			panic("got nil OrRight but non-nil OrLeft")
		}
		b.Write(grammar.Symbols[grammar.SYM_LPAREN])
		b.doBooleanValueExpression(el.OrLeft, qargs, curarg)
		b.Write(grammar.Symbols[grammar.SYM_OR])
		b.doBooleanTerm(el.OrRight, qargs, curarg)
		b.Write(grammar.Symbols[grammar.SYM_RPAREN])
	}
}

func (b *Builder) doBooleanTerm(
	el *grammar.BooleanTerm,
	qargs []interface{},
	curarg *int,
) {
	if el.Unary != nil {
		b.doBooleanFactor(el.Unary, qargs, curarg)
	} else if el.AndLeft != nil {
		if el.AndRight == nil {
			// This should not happen, so if it does, panic since it's a fault
			// in sqlb's parsing logic.
			panic("got nil AndRight but non-nil AndLeft")
		}
		b.doBooleanTerm(el.AndLeft, qargs, curarg)
		b.Write(grammar.Symbols[grammar.SYM_AND])
		b.doBooleanFactor(el.AndRight, qargs, curarg)
	}
}

func (b *Builder) doBooleanFactor(
	el *grammar.BooleanFactor,
	qargs []interface{},
	curarg *int,
) {
	if el.Not {
		b.Write(grammar.Symbols[grammar.SYM_NOT])
	}
	b.doBooleanPrimary(&el.Test.Primary, qargs, curarg)
}

func (b *Builder) doBooleanPrimary(
	el *grammar.BooleanPrimary,
	qargs []interface{},
	curarg *int,
) {
	if el.Predicate != nil {
		b.doPredicate(el.Predicate, qargs, curarg)
	} else if el.Predicand != nil {
		b.doBooleanPredicand(el.Predicand, qargs, curarg)
	}
}

func (b *Builder) doBooleanPredicand(
	el *grammar.BooleanPredicand,
	qargs []interface{},
	curarg *int,
) {
	if el.Parenthesized != nil {
		b.Write(grammar.Symbols[grammar.SYM_LPAREN])
		b.doBooleanValueExpression(el.Parenthesized, qargs, curarg)
		b.Write(grammar.Symbols[grammar.SYM_RPAREN])
	} else if el.Primary != nil {
		b.doNonParenthesizedValueExpressionPrimary(el.Primary, qargs, curarg)
	}
}
