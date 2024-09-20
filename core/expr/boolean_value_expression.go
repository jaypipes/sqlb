//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package expr

import (
	"fmt"

	"github.com/jaypipes/sqlb/grammar"
	"github.com/jaypipes/sqlb/internal/inspect"
)

// And accepts two things and returns a BooleanValueExpression ANDing the two
// things together. This boolean value expression can be passed to a Join or
// Where clause.
func And(
	leftAny interface{},
	rightAny interface{},
) *grammar.BooleanValueExpression {
	left := inspect.BooleanTermFromAny(leftAny)
	if left == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected BooleanTerm",
			leftAny, leftAny,
		)
		panic(msg)
	}
	right := inspect.BooleanFactorFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected BooleanFactor",
			rightAny, rightAny,
		)
		panic(msg)
	}
	return &grammar.BooleanValueExpression{
		Unary: &grammar.BooleanTerm{
			AndLeft:  left,
			AndRight: right,
		},
	}
}

// Or accepts two things and returns an Element representing an OR expression
// that can be passed to a Join or Where clause.
func Or(
	leftAny interface{},
	rightAny interface{},
) *grammar.BooleanValueExpression {
	left := inspect.BooleanValueExpressionFromAny(leftAny)
	if left == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected BooleanValueExpression",
			leftAny, leftAny,
		)
		panic(msg)
	}
	right := inspect.BooleanTermFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected BooleanTerm",
			rightAny, rightAny,
		)
		panic(msg)
	}
	return &grammar.BooleanValueExpression{
		OrLeft:  left,
		OrRight: right,
	}
}
