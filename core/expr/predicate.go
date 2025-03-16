//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package expr

import (
	"fmt"

	"github.com/jaypipes/sqlb/core/grammar"
	"github.com/jaypipes/sqlb/internal/inspect"
)

// Equal accepts two things and returns a ComparisonPredicate representing an
// equality expression that can be passed to a Join or Where clause.
//
// Equal panics if sqlb cannot compile the supplied arguments into a valid
// ComparisonPredicate. This is intentional, as we want compile-time failures
// for invalid SQL construction and we want the result of Equal() to be passed
// directly into other `core/expr` functions.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `EqualE` function which returns a checkable `error` object.
func Equal(
	leftAny interface{},
	rightAny interface{},
) *grammar.ComparisonPredicate {
	p, err := EqualE(leftAny, rightAny)
	if err != nil {
		panic(err)
	}
	return p
}

// EqualE accepts two things and returns a ComparisonPredicate representing an
// equality expression that can be passed to a Join or Where clause. If the two
// supplied parameters cannot be evaluated into a ComparisonPredicate, an error
// is returned.
func EqualE(
	leftAny interface{},
	rightAny interface{},
) (*grammar.ComparisonPredicate, error) {
	left := inspect.RowValuePredicandFromAny(leftAny)
	if left == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			leftAny, leftAny,
		)
	}
	right := inspect.RowValuePredicandFromAny(rightAny)
	if right == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			rightAny, rightAny,
		)
	}
	return &grammar.ComparisonPredicate{
		Operator: grammar.ComparisonOperatorEquals,
		A:        *left,
		B:        *right,
	}, nil
}

// NotEqual accepts two things and returns a ComparisonPredicate representing
// an inequality expression that can be passed to a Join or Where clause.
//
// NotEqual panics if sqlb cannot compile the supplied arguments into a valid
// ComparisonPredicate. This is intentional, as we want compile-time failures
// for invalid SQL construction and we want the result of NotEqual() to be
// passed directly into other `core/expr` functions.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `NotEqualE` function which returns a checkable `error` object.
func NotEqual(
	leftAny interface{},
	rightAny interface{},
) *grammar.ComparisonPredicate {
	p, err := NotEqualE(leftAny, rightAny)
	if err != nil {
		panic(err)
	}
	return p
}

// NotEqualE accepts two things and returns a ComparisonPredicate representing
// an inequality expression that can be passed to a Join or Where clause.  If
// the two supplied parameters cannot be evaluated into a ComparisonPredicate,
// an error is returned.
func NotEqualE(
	leftAny interface{},
	rightAny interface{},
) (*grammar.ComparisonPredicate, error) {
	left := inspect.RowValuePredicandFromAny(leftAny)
	if left == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			leftAny, leftAny,
		)
	}
	right := inspect.RowValuePredicandFromAny(rightAny)
	if right == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			rightAny, rightAny,
		)
	}
	return &grammar.ComparisonPredicate{
		Operator: grammar.ComparisonOperatorNotEquals,
		A:        *left,
		B:        *right,
	}, nil
}

// GreaterThan accepts two things and returns a ComparisonPredicate
// representing greater than expression that can be passed to a Join or Where
// clause.
//
// GreaterThan panics if sqlb cannot compile the supplied arguments into a
// valid ComparisonPredicate. This is intentional, as we want compile-time
// failures for invalid SQL construction and we want the result of
// GreaterThan() to be passed directly into other `core/expr` functions.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `GreaterThanE` function which returns a checkable `error`
// object.
func GreaterThan(
	leftAny interface{},
	rightAny interface{},
) *grammar.ComparisonPredicate {
	p, err := GreaterThanE(leftAny, rightAny)
	if err != nil {
		panic(err)
	}
	return p
}

// GreaterThanE accepts two things and returns a ComparisonPredicate
// representing greater than expression that can be passed to a Join or Where
// clause. If the two supplied parameters cannot be evaluated into a
// ComparisonPredicate, an error is returned.
func GreaterThanE(
	leftAny interface{},
	rightAny interface{},
) (*grammar.ComparisonPredicate, error) {
	left := inspect.RowValuePredicandFromAny(leftAny)
	if left == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			leftAny, leftAny,
		)
	}
	right := inspect.RowValuePredicandFromAny(rightAny)
	if right == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			rightAny, rightAny,
		)
	}
	return &grammar.ComparisonPredicate{
		Operator: grammar.ComparisonOperatorGreaterThan,
		A:        *left,
		B:        *right,
	}, nil
}

// GreaterThanOrEqual accepts two things and returns a ComparisonPredicate
// representing greater than or equal expression that can be passed to a Join
// or Where clause.
//
// GreaterThanOrEqual panics if sqlb cannot compile the supplied arguments into
// a valid ComparisonPredicate. This is intentional, as we want compile-time
// failures for invalid SQL construction and we want the result of
// GreaterThanOrEqual() to be passed directly into other `core/expr` functions.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `GreaterThanOrEqualE` function which returns a checkable
// `error` object.
func GreaterThanOrEqual(
	leftAny interface{},
	rightAny interface{},
) *grammar.ComparisonPredicate {
	p, err := GreaterThanOrEqualE(leftAny, rightAny)
	if err != nil {
		panic(err)
	}
	return p
}

// GreaterThanOrEqualE accepts two things and returns a ComparisonPredicate
// representing greater than or equal expression that can be passed to a Join
// or Where clause. If the two supplied parameters cannot be evaluated into a
// ComparisonPredicate, an error is returned.
func GreaterThanOrEqualE(
	leftAny interface{},
	rightAny interface{},
) (*grammar.ComparisonPredicate, error) {
	left := inspect.RowValuePredicandFromAny(leftAny)
	if left == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			leftAny, leftAny,
		)
	}
	right := inspect.RowValuePredicandFromAny(rightAny)
	if right == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			rightAny, rightAny,
		)
	}
	return &grammar.ComparisonPredicate{
		Operator: grammar.ComparisonOperatorGreaterThanEquals,
		A:        *left,
		B:        *right,
	}, nil
}

// LessThan accepts two things and returns a ComparisonPredicate representing
// less than expression that can be passed to a Join or Where clause.
//
// LessThan panics if sqlb cannot compile the supplied arguments into a valid
// ComparisonPredicate. This is intentional, as we want compile-time failures
// for invalid SQL construction and we want the result of LessThan() to be
// passed directly into other `core/expr` functions.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `LessThanE` function which returns a checkable `error`
// object.
func LessThan(
	leftAny interface{},
	rightAny interface{},
) *grammar.ComparisonPredicate {
	p, err := LessThanE(leftAny, rightAny)
	if err != nil {
		panic(err)
	}
	return p
}

// LessThanE accepts two things and returns a ComparisonPredicate representing
// less than expression that can be passed to a Join or Where clause. If the
// two supplied parameters cannot be evaluated into a ComparisonPredicate, an
// error is returned.
func LessThanE(
	leftAny interface{},
	rightAny interface{},
) (*grammar.ComparisonPredicate, error) {
	left := inspect.RowValuePredicandFromAny(leftAny)
	if left == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			leftAny, leftAny,
		)
	}
	right := inspect.RowValuePredicandFromAny(rightAny)
	if right == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			rightAny, rightAny,
		)
	}
	return &grammar.ComparisonPredicate{
		Operator: grammar.ComparisonOperatorLessThan,
		A:        *left,
		B:        *right,
	}, nil
}

// LessThanOrEqual accepts two things and returns a ComparisonPredicate
// representing less than or equal expression that can be passed to a Join or
// Where clause.
//
// LessThanOrEqual panics if sqlb cannot compile the supplied arguments into
// a valid ComparisonPredicate. This is intentional, as we want compile-time
// failures for invalid SQL construction and we want the result of
// LessThanOrEqual() to be passed directly into other `core/expr` functions.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `LessThanOrEqualE` function which returns a checkable
// `error` object.
func LessThanOrEqual(
	leftAny interface{},
	rightAny interface{},
) *grammar.ComparisonPredicate {
	p, err := LessThanOrEqualE(leftAny, rightAny)
	if err != nil {
		panic(err)
	}
	return p
}

// LessThanOrEqualE accepts two things and returns a ComparisonPredicate
// representing less than or equal expression that can be passed to a Join or
// Where clause. If the two supplied parameters cannot be evaluated into a
// ComparisonPredicate, an error is returned.
func LessThanOrEqualE(
	leftAny interface{},
	rightAny interface{},
) (*grammar.ComparisonPredicate, error) {
	left := inspect.RowValuePredicandFromAny(leftAny)
	if left == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			leftAny, leftAny,
		)
	}
	right := inspect.RowValuePredicandFromAny(rightAny)
	if right == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			rightAny, rightAny,
		)
	}
	return &grammar.ComparisonPredicate{
		Operator: grammar.ComparisonOperatorLessThanEquals,
		A:        *left,
		B:        *right,
	}, nil
}

// Between accepts three things and returns a BetweenPredicate representing a
// SQL BETWEEN expression that can be passed to a Join or Where clause.
//
// Between panics if sqlb cannot compile the supplied arguments into a valid
// BetweenPredicate. This is intentional, as we want compile-time failures for
// invalid SQL construction and we want the result of Between() to be passed
// directly into other `core/expr` functions.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `BetweenE` function which returns a checkable `error` object.
func Between(
	targetAny interface{},
	startAny interface{},
	endAny interface{},
) *grammar.BetweenPredicate {
	p, err := BetweenE(targetAny, startAny, endAny)
	if err != nil {
		panic(err)
	}
	return p
}

// BetweenE accepts three things and returns a BetweenPredicate representing a
// SQL BETWEEN expression that can be passed to a Join or Where clause. If the
// supplied arguments cannot be compiled into a valid BetweenPredicate, an
// error is returned.
func BetweenE(
	targetAny interface{},
	startAny interface{},
	endAny interface{},
) (*grammar.BetweenPredicate, error) {
	target := inspect.RowValuePredicandFromAny(targetAny)
	if target == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			targetAny, targetAny,
		)
	}
	start := inspect.RowValuePredicandFromAny(startAny)
	if start == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			startAny, startAny,
		)
	}
	end := inspect.RowValuePredicandFromAny(endAny)
	if end == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			endAny, endAny,
		)
	}
	return &grammar.BetweenPredicate{
		Target: *target,
		Start:  *start,
		End:    *end,
	}, nil
}

// In accepts two things and returns an InPredicate representing an IN
// expression that can be passed to a Join or Where clause.
//
// In panics if sqlb cannot compile the supplied arguments into a valid
// InPredicate. This is intentional, as we want compile-time failures for
// invalid SQL construction and we want the result of In() to be passed
// directly into other `core/expr` functions.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `InE` function which returns a checkable `error` object.
func In(
	targetAny interface{},
	values ...interface{},
) *grammar.InPredicate {
	p, err := InE(targetAny, values...)
	if err != nil {
		panic(err)
	}
	return p
}

// In accepts two things and returns an InPredicate representing an IN
// expression that can be passed to a Join or Where clause. If the supplied
// arguments cannot be compiled into a valid InPredicate, an error is returned.
func InE(
	targetAny interface{},
	values ...interface{},
) (*grammar.InPredicate, error) {
	target := inspect.RowValuePredicandFromAny(targetAny)
	if target == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			targetAny, targetAny,
		)
	}
	rves := []grammar.RowValueExpression{}
	for _, v := range values {
		npvep := inspect.NonParenthesizedValueExpressionPrimaryFromAny(v)
		if npvep == nil {
			return nil, fmt.Errorf(
				"could not convert %s(%T) to expected inspect.NonParenthesizedValueExpressionPrimary",
				v, v,
			)
		}
		rves = append(rves, grammar.RowValueExpression{
			Primary: npvep,
		})
	}
	return &grammar.InPredicate{
		Target: *target,
		Values: rves,
	}, nil
}

// IsNull accepts a thing and returns a NullPredicate representing an IS NULL
// expression that can be passed to a Join or Where clause.
//
// IsNull panics if sqlb cannot compile the supplied arguments into a valid
// NullPredicate. This is intentional, as we want compile-time failures for
// invalid SQL construction and we want the result of IsNull() to be passed
// directly into other `core/expr` functions.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `IsNullE` function which returns a checkable `error` object.
func IsNull(
	targetAny interface{},
) *grammar.NullPredicate {
	p, err := IsNullE(targetAny)
	if err != nil {
		panic(err)
	}
	return p
}

// IsNullE accepts a thing and returns a NullPredicate representing an IS NULL
// expression that can be passed to a Join or Where clause. If the supplied
// parameter cannot be converted into a RowValuePredicand, an error is
// returned.
func IsNullE(
	targetAny interface{},
) (*grammar.NullPredicate, error) {
	target := inspect.RowValuePredicandFromAny(targetAny)
	if target == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			targetAny, targetAny,
		)
	}
	return &grammar.NullPredicate{
		Target: *target,
	}, nil
}

// IsNotNull accepts a thing and returns a NullPredicate representing an IS NOT
// NULL expression that can be passed to a Join or Where clause.
//
// IsNotNull panics if sqlb cannot compile the supplied arguments into a valid
// NullPredicate. This is intentional, as we want compile-time failures for
// invalid SQL construction and we want the result of IsNotNull() to be passed
// directly into other `core/expr` functions.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `IsNotNullE` function which returns a checkable `error` object.
func IsNotNull(
	targetAny interface{},
) *grammar.NullPredicate {
	p, err := IsNotNullE(targetAny)
	if err != nil {
		panic(err)
	}
	return p
}

// IsNotNullE accepts a thing and returns a NullPredicate representing an IS
// NOT NULL expression that can be passed to a Join or Where clause. If the
// supplied parameter cannot be converted into a RowValuePredicand, an error is
// returned.
func IsNotNullE(
	targetAny interface{},
) (*grammar.NullPredicate, error) {
	target := inspect.RowValuePredicandFromAny(targetAny)
	if target == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected inspect.RowValuePredicand",
			targetAny, targetAny,
		)
	}
	return &grammar.NullPredicate{
		Target: *target,
		Not:    true,
	}, nil
}
