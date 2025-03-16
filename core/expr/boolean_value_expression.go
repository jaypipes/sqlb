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

// And accepts two things and returns a BooleanValueExpression ANDing the two
// things together. This boolean value expression can be passed to a Join or
// Where clause.
//
// And panics if sqlb cannot compile the supplied arguments into a valid
// BooleanValueExpression. This is intentional, as we want compile-time
// failures for invalid SQL construction and we want the result of And() to be
// passed directly into other `core/expr` functions.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `AndE` function which returns a checkable `error` object.
func And(
	leftAny interface{},
	rightAny interface{},
) *grammar.BooleanValueExpression {
	p, err := AndE(leftAny, rightAny)
	if err != nil {
		panic(err)
	}
	return p
}

// AndE accepts two things and returns a BooleanValueExpression ANDing the two
// things together. This boolean value expression can be passed to a Join or
// Where clause. If the two parameters cannot be compiled into a
// BooleanValueExpression, an error is returned.
func AndE(
	leftAny interface{},
	rightAny interface{},
) (*grammar.BooleanValueExpression, error) {
	left := inspect.BooleanTermFromAny(leftAny)
	if left == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected BooleanTerm",
			leftAny, leftAny,
		)
	}
	right := inspect.BooleanFactorFromAny(rightAny)
	if right == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected BooleanFactor",
			rightAny, rightAny,
		)
	}
	return &grammar.BooleanValueExpression{
		Unary: &grammar.BooleanTerm{
			AndLeft:  left,
			AndRight: right,
		},
	}, nil
}

// Or accepts two things and returns an Element representing an OR expression
// that can be passed to a Join or Where clause.
//
// Or panics if sqlb cannot compile the supplied arguments into a valid
// BooleanValueExpression. This is intentional, as we want compile-time
// failures for invalid SQL construction and we want the result of Or() to be
// passed directly into other `core/expr` functions.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `OrE` function which returns a checkable `error` object.
func Or(
	leftAny interface{},
	rightAny interface{},
) *grammar.BooleanValueExpression {
	p, err := OrE(leftAny, rightAny)
	if err != nil {
		panic(err)
	}
	return p
}

// OrE accepts two things and returns a BooleanValueExpression ORing the two
// things together. This boolean value expression can be passed to a Join or
// Where clause. If the two parameters cannot be compiled into a
// BooleanValueExpression, an error is returned.
func OrE(
	leftAny interface{},
	rightAny interface{},
) (*grammar.BooleanValueExpression, error) {
	left := inspect.BooleanValueExpressionFromAny(leftAny)
	if left == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected BooleanValueExpression",
			leftAny, leftAny,
		)
	}
	right := inspect.BooleanTermFromAny(rightAny)
	if right == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected BooleanTerm",
			rightAny, rightAny,
		)
	}
	return &grammar.BooleanValueExpression{
		OrLeft:  left,
		OrRight: right,
	}, nil
}
