//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package expr

import (
	"fmt"

	"github.com/jaypipes/sqlb/internal/inspect"
	"github.com/jaypipes/sqlb/grammar"
)

// Equal accepts two things and returns a ComparisonPredicate representing an
// equality expression that can be passed to a Join or Where clause.
func Equal(
	leftAny interface{},
	rightAny interface{},
) *grammar.ComparisonPredicate {
	left := inspect.RowValuePredicandFromAny(leftAny)
	if left == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			leftAny, leftAny,
		)
		panic(msg)
	}
	right := inspect.RowValuePredicandFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
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
	left := inspect.RowValuePredicandFromAny(leftAny)
	if left == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			leftAny, leftAny,
		)
		panic(msg)
	}
	right := inspect.RowValuePredicandFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
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
	left := inspect.RowValuePredicandFromAny(leftAny)
	if left == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			leftAny, leftAny,
		)
		panic(msg)
	}
	right := inspect.RowValuePredicandFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
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
	left := inspect.RowValuePredicandFromAny(leftAny)
	if left == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			leftAny, leftAny,
		)
		panic(msg)
	}
	right := inspect.RowValuePredicandFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
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
	left := inspect.RowValuePredicandFromAny(leftAny)
	if left == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			leftAny, leftAny,
		)
		panic(msg)
	}
	right := inspect.RowValuePredicandFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
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
	left := inspect.RowValuePredicandFromAny(leftAny)
	if left == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			leftAny, leftAny,
		)
		panic(msg)
	}
	right := inspect.RowValuePredicandFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
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
	target := inspect.RowValuePredicandFromAny(targetAny)
	if target == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			targetAny, targetAny,
		)
		panic(msg)
	}
	start := inspect.RowValuePredicandFromAny(startAny)
	if start == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			startAny, startAny,
		)
		panic(msg)
	}
	end := inspect.RowValuePredicandFromAny(endAny)
	if end == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
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
	target := inspect.RowValuePredicandFromAny(targetAny)
	if target == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			targetAny, targetAny,
		)
		panic(msg)
	}
	rves := []grammar.RowValueExpression{}
	for _, v := range values {
		npvep := inspect.NonParenthesizedValueExpressionPrimaryFromAny(v)
		if npvep == nil {
			msg := fmt.Sprintf(
				"could not convert %s(%T) to expected inspect.NonParenthesizedValueExpressionPrimary",
				v, v,
			)
			panic(msg)
		}
		rves = append(rves, grammar.RowValueExpression{
			Primary: npvep,
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
	target := inspect.RowValuePredicandFromAny(targetAny)
	if target == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
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
	target := inspect.RowValuePredicandFromAny(targetAny)
	if target == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			targetAny, targetAny,
		)
		panic(msg)
	}
	return &grammar.NullPredicate{
		Target: *target,
		Not:    true,
	}
}
