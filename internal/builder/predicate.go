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

func (b *Builder) doPredicate(
	el *grammar.Predicate,
	qargs []interface{},
	curarg *int,
) {
	if el.Comparison != nil {
		b.doComparisonPredicate(el.Comparison, qargs, curarg)
	} else if el.In != nil {
		b.doInPredicate(el.In, qargs, curarg)
	} else if el.Between != nil {
		b.doBetweenPredicate(el.Between, qargs, curarg)
	} else if el.Null != nil {
		b.doNullPredicate(el.Null, qargs, curarg)
	}
}

func (b *Builder) doComparisonPredicate(
	el *grammar.ComparisonPredicate,
	qargs []interface{},
	curarg *int,
) {
	b.doRowValuePredicand(&el.A, qargs, curarg)
	switch el.Operator {
	case grammar.ComparisonOperatorEquals:
		b.WriteString(symbol.Space)
		b.WriteString(symbol.EqualsOperator)
		b.WriteString(symbol.Space)
	case grammar.ComparisonOperatorNotEquals:
		b.WriteString(symbol.Space)
		b.WriteString(symbol.LessThanOperator)
		b.WriteString(symbol.GreaterThanOperator)
		b.WriteString(symbol.Space)
	case grammar.ComparisonOperatorGreaterThan:
		b.WriteString(symbol.Space)
		b.WriteString(symbol.GreaterThanOperator)
		b.WriteString(symbol.Space)
	case grammar.ComparisonOperatorGreaterThanEquals:
		b.WriteString(symbol.Space)
		b.WriteString(symbol.GreaterThanOperator)
		b.WriteString(symbol.EqualsOperator)
		b.WriteString(symbol.Space)
	case grammar.ComparisonOperatorLessThan:
		b.WriteString(symbol.Space)
		b.WriteString(symbol.LessThanOperator)
		b.WriteString(symbol.Space)
	case grammar.ComparisonOperatorLessThanEquals:
		b.WriteString(symbol.Space)
		b.WriteString(symbol.LessThanOperator)
		b.WriteString(symbol.EqualsOperator)
		b.WriteString(symbol.Space)
	}
	b.doRowValuePredicand(&el.B, qargs, curarg)
}

func (b *Builder) doInPredicate(
	el *grammar.InPredicate,
	qargs []interface{},
	curarg *int,
) {
	b.doRowValuePredicand(&el.Target, qargs, curarg)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.In)
	b.WriteString(symbol.LeftParen)
	for x, rve := range el.Values {
		if x > 0 {
			b.WriteString(symbol.Comma)
			b.WriteString(symbol.Space)
		}
		b.doNonParenthesizedValueExpressionPrimary(rve.Primary, qargs, curarg)
	}
	b.WriteString(symbol.RightParen)
}

func (b *Builder) doBetweenPredicate(
	el *grammar.BetweenPredicate,
	qargs []interface{},
	curarg *int,
) {
	b.doRowValuePredicand(&el.Target, qargs, curarg)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.Between)
	b.WriteString(symbol.Space)
	b.doRowValuePredicand(&el.Start, qargs, curarg)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.And)
	b.WriteString(symbol.Space)
	b.doRowValuePredicand(&el.End, qargs, curarg)
}

func (b *Builder) doNullPredicate(
	el *grammar.NullPredicate,
	qargs []interface{},
	curarg *int,
) {
	b.doRowValuePredicand(&el.Target, qargs, curarg)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.Is)
	if el.Not {
		b.WriteString(symbol.Space)
		b.WriteString(symbol.Not)
	}
	b.WriteString(symbol.Space)
	b.WriteString(symbol.Null)
}
