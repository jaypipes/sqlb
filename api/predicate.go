//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import (
	"fmt"

	"github.com/jaypipes/sqlb/grammar"
)

// PredicateFromAny evaluates the supplied interface argument and
// returns a *Predicate if the supplied argument can be converted
// into a Predicate, or nil if the conversion cannot be done.
func PredicateFromAny(v interface{}) *grammar.Predicate {
	switch v := v.(type) {
	case *grammar.Predicate:
		return v
	case grammar.Predicate:
		return &v
	case *grammar.ComparisonPredicate:
		return &grammar.Predicate{
			Comparison: v,
		}
	case grammar.ComparisonPredicate:
		return &grammar.Predicate{
			Comparison: &v,
		}
	case *grammar.BetweenPredicate:
		return &grammar.Predicate{
			Between: v,
		}
	case grammar.BetweenPredicate:
		return &grammar.Predicate{
			Between: &v,
		}
	case *grammar.InPredicate:
		return &grammar.Predicate{
			In: v,
		}
	case grammar.InPredicate:
		return &grammar.Predicate{
			In: &v,
		}
	case *grammar.NullPredicate:
		return &grammar.Predicate{
			Null: v,
		}
	case grammar.NullPredicate:
		return &grammar.Predicate{
			Null: &v,
		}
	}
	return nil
}

// Equal accepts two things and returns a ComparisonPredicate representing an
// equality expression that can be passed to a Join or Where clause.
func Equal(
	leftAny interface{},
	rightAny interface{},
) *grammar.ComparisonPredicate {
	left := RowValuePredicandFromAny(leftAny)
	if left == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			leftAny, leftAny,
		)
		panic(msg)
	}
	right := RowValuePredicandFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			rightAny, rightAny,
		)
		panic(msg)
	}
	return &grammar.ComparisonPredicate{
		Operator: grammar.ComparisonOperatorEquals,
		A:        *left,
		B:        *right,
	}
}

// NotEqual accepts two things and returns a ComparisonPredicate representing
// an inequality expression that can be passed to a Join or Where clause.
func NotEqual(
	leftAny interface{},
	rightAny interface{},
) *grammar.ComparisonPredicate {
	left := RowValuePredicandFromAny(leftAny)
	if left == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			leftAny, leftAny,
		)
		panic(msg)
	}
	right := RowValuePredicandFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			rightAny, rightAny,
		)
		panic(msg)
	}
	return &grammar.ComparisonPredicate{
		Operator: grammar.ComparisonOperatorNotEquals,
		A:        *left,
		B:        *right,
	}
}

// GreaterThan accepts two things and returns a ComparisonPredicate
// representing greater than expression that can be passed to a Join or Where
// clause.
func GreaterThan(
	leftAny interface{},
	rightAny interface{},
) *grammar.ComparisonPredicate {
	left := RowValuePredicandFromAny(leftAny)
	if left == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			leftAny, leftAny,
		)
		panic(msg)
	}
	right := RowValuePredicandFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			rightAny, rightAny,
		)
		panic(msg)
	}
	return &grammar.ComparisonPredicate{
		Operator: grammar.ComparisonOperatorGreaterThan,
		A:        *left,
		B:        *right,
	}
}

// GreaterThanOrEqual accepts two things and returns a ComparisonPredicate
// representing greater than or equal expression that can be passed to a Join
// or Where clause.
func GreaterThanOrEqual(
	leftAny interface{},
	rightAny interface{},
) *grammar.ComparisonPredicate {
	left := RowValuePredicandFromAny(leftAny)
	if left == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			leftAny, leftAny,
		)
		panic(msg)
	}
	right := RowValuePredicandFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			rightAny, rightAny,
		)
		panic(msg)
	}
	return &grammar.ComparisonPredicate{
		Operator: grammar.ComparisonOperatorGreaterThanEquals,
		A:        *left,
		B:        *right,
	}
}

// LessThan accepts two things and returns a ComparisonPredicate
// representing greater than expression that can be passed to a Join or Where
// clause.
func LessThan(
	leftAny interface{},
	rightAny interface{},
) *grammar.ComparisonPredicate {
	left := RowValuePredicandFromAny(leftAny)
	if left == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			leftAny, leftAny,
		)
		panic(msg)
	}
	right := RowValuePredicandFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			rightAny, rightAny,
		)
		panic(msg)
	}
	return &grammar.ComparisonPredicate{
		Operator: grammar.ComparisonOperatorLessThan,
		A:        *left,
		B:        *right,
	}
}

// LessThanOrEqual accepts two things and returns a ComparisonPredicate
// representing greater than or equal expression that can be passed to a Join
// or Where clause.
func LessThanOrEqual(
	leftAny interface{},
	rightAny interface{},
) *grammar.ComparisonPredicate {
	left := RowValuePredicandFromAny(leftAny)
	if left == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			leftAny, leftAny,
		)
		panic(msg)
	}
	right := RowValuePredicandFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			rightAny, rightAny,
		)
		panic(msg)
	}
	return &grammar.ComparisonPredicate{
		Operator: grammar.ComparisonOperatorLessThanEquals,
		A:        *left,
		B:        *right,
	}
}

// Between accepts three things and returns a BetweenPredicate representing a
// SQL BETWEEN expression that can be passed to a Join or Where clause.
func Between(
	targetAny interface{},
	startAny interface{},
	endAny interface{},
) *grammar.BetweenPredicate {
	target := RowValuePredicandFromAny(targetAny)
	if target == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			targetAny, targetAny,
		)
		panic(msg)
	}
	start := RowValuePredicandFromAny(startAny)
	if start == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			startAny, startAny,
		)
		panic(msg)
	}
	end := RowValuePredicandFromAny(endAny)
	if end == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			endAny, endAny,
		)
		panic(msg)
	}
	return &grammar.BetweenPredicate{
		Target: *target,
		Start:  *start,
		End:    *end,
	}
}

// In accepts two things and returns an InPredicate representing an IN
// expression that can be passed to a Join or Where clause.
func In(
	targetAny interface{},
	values ...interface{},
) *grammar.InPredicate {
	target := RowValuePredicandFromAny(targetAny)
	if target == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			targetAny, targetAny,
		)
		panic(msg)
	}
	rves := []grammar.RowValueExpression{}
	for _, v := range values {
		npvep := NonParenthesizedValueExpressionPrimaryFromAny(v)
		if npvep == nil {
			msg := fmt.Sprintf(
				"could not convert %s(%T) to expected NonParenthesizedValueExpressionPrimary",
				v, v,
			)
			panic(msg)
		}
		rves = append(rves, grammar.RowValueExpression{
			NonParenthesizedValueExpressionPrimary: npvep,
		})
	}
	return &grammar.InPredicate{
		Target: *target,
		Values: rves,
	}
}

// IsNull accepts a thing and returns a NullPredicate representing an IS NULL
// expression that can be passed to a Join or Where clause.
func IsNull(
	targetAny interface{},
) *grammar.NullPredicate {
	target := RowValuePredicandFromAny(targetAny)
	if target == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			targetAny, targetAny,
		)
		panic(msg)
	}
	return &grammar.NullPredicate{
		Target: *target,
	}
}

// IsNotNull accepts a thing and returns a NullPredicate representing an IS NOT
// NULL expression that can be passed to a Join or Where clause.
func IsNotNull(
	targetAny interface{},
) *grammar.NullPredicate {
	target := RowValuePredicandFromAny(targetAny)
	if target == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected RowValuePredicand",
			targetAny, targetAny,
		)
		panic(msg)
	}
	return &grammar.NullPredicate{
		Target: *target,
		Not:    true,
	}
}

// ReferredFromPredicate returns a slice of string names of any tables or
// derived tables (subqueries in the FROM clause) that are referenced within a
// supplied Predicate.
func ReferredFromPredicate(
	p *grammar.Predicate,
) []string {
	if p.Comparison != nil {
		found := ReferredFromRowValuePredicand(&p.Comparison.A)
		found = append(found, ReferredFromRowValuePredicand(&p.Comparison.B)...)
		return found
	} else if p.In != nil {
		found := ReferredFromRowValuePredicand(&p.In.Target)
		for _, rve := range p.In.Values {
			found = append(
				found,
				ReferredFromNonParenthesizedValueExpressionPrimary(rve.NonParenthesizedValueExpressionPrimary)...,
			)
		}
		return found
	} else if p.Between != nil {
		found := ReferredFromRowValuePredicand(&p.Between.Target)
		found = append(found, ReferredFromRowValuePredicand(&p.Between.Start)...)
		found = append(found, ReferredFromRowValuePredicand(&p.Between.End)...)
		return found
	} else if p.Null != nil {
		return ReferredFromRowValuePredicand(&p.Null.Target)
	}
	return []string{}
}
