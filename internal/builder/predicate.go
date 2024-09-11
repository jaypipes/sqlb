//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
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
		b.Write(grammar.Symbols[grammar.SYM_EQUAL])
	case grammar.ComparisonOperatorNotEquals:
		b.Write(grammar.Symbols[grammar.SYM_NEQUAL])
	case grammar.ComparisonOperatorGreaterThan:
		b.Write(grammar.Symbols[grammar.SYM_GREATER])
	case grammar.ComparisonOperatorGreaterThanEquals:
		b.Write(grammar.Symbols[grammar.SYM_GREATER_EQUAL])
	case grammar.ComparisonOperatorLessThan:
		b.Write(grammar.Symbols[grammar.SYM_LESS])
	case grammar.ComparisonOperatorLessThanEquals:
		b.Write(grammar.Symbols[grammar.SYM_LESS_EQUAL])
	}
	b.doRowValuePredicand(&el.B, qargs, curarg)
}

func (b *Builder) doInPredicate(
	el *grammar.InPredicate,
	qargs []interface{},
	curarg *int,
) {
	b.doRowValuePredicand(&el.Target, qargs, curarg)
	b.Write(grammar.Symbols[grammar.SYM_IN])
	b.Write(grammar.Symbols[grammar.SYM_LPAREN])
	for x, rve := range el.Values {
		if x > 0 {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
		b.doNonParenthesizedValueExpressionPrimary(rve.Primary, qargs, curarg)
	}
	b.Write(grammar.Symbols[grammar.SYM_RPAREN])
}

func (b *Builder) doBetweenPredicate(
	el *grammar.BetweenPredicate,
	qargs []interface{},
	curarg *int,
) {
	b.doRowValuePredicand(&el.Target, qargs, curarg)
	b.Write(grammar.Symbols[grammar.SYM_BETWEEN])
	b.doRowValuePredicand(&el.Start, qargs, curarg)
	b.Write(grammar.Symbols[grammar.SYM_AND])
	b.doRowValuePredicand(&el.End, qargs, curarg)
}

func (b *Builder) doNullPredicate(
	el *grammar.NullPredicate,
	qargs []interface{},
	curarg *int,
) {
	b.doRowValuePredicand(&el.Target, qargs, curarg)
	if el.Not {
		b.Write(grammar.Symbols[grammar.SYM_IS_NOT_NULL])
	} else {
		b.Write(grammar.Symbols[grammar.SYM_IS_NULL])
	}
}
