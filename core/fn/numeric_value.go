//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package fn

import (
	"fmt"

	"github.com/jaypipes/sqlb/core/inspect"
	"github.com/jaypipes/sqlb/core/types"
	"github.com/jaypipes/sqlb/grammar"
)

// CharacterLength returns a LengthExpression that produces a CHAR_LENGTH() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The first argument is the subject of the CHAR_LENGTH function and must be
// coercible to a string value expression.
func CharacterLength(
	subjectAny interface{},
) *LengthExpression {
	var ref types.Relation
	switch subjectAny := subjectAny.(type) {
	case types.Projection:
		ref = subjectAny.References()
	}
	subject := inspect.StringValueExpressionFromAny(subjectAny)
	if subject == nil {
		msg := fmt.Sprintf(
			"expected coerceable StringValueExpression but got %+v(%T)",
			subjectAny, subjectAny,
		)
		panic(msg)
	}
	return &LengthExpression{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		LengthExpression: &grammar.LengthExpression{
			Character: &grammar.CharacterLengthExpression{
				Subject: *subject,
			},
		},
	}
}

var CharLength = CharacterLength

// OctetLength returns a LengthExpression that produces a OCTET_LENGTH() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The first argument is the subject of the OCTET_LENGTH function and must be
// coercible to a string value expression.
func OctetLength(
	subjectAny interface{},
) *LengthExpression {
	var ref types.Relation
	switch subjectAny := subjectAny.(type) {
	case types.Projection:
		ref = subjectAny.References()
	}
	subject := inspect.StringValueExpressionFromAny(subjectAny)
	if subject == nil {
		msg := fmt.Sprintf(
			"expected coerceable StringValueExpression but got %+v(%T)",
			subjectAny, subjectAny,
		)
		panic(msg)
	}
	return &LengthExpression{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		LengthExpression: &grammar.LengthExpression{
			Octet: &grammar.OctetLengthExpression{
				Subject: *subject,
			},
		},
	}
}

// LengthExpression wraps the CHAR_LENGTH() SQL function grammar element
type LengthExpression struct {
	BaseFunction
	*grammar.LengthExpression
}

// CommonValueExpression returns the object as a
// `*grammar.CommonValueExpression`
func (f *LengthExpression) CommonValueExpression() *grammar.CommonValueExpression {
	return &grammar.CommonValueExpression{
		Numeric: &grammar.NumericValueExpression{
			Unary: &grammar.Term{
				Unary: &grammar.Factor{
					Primary: grammar.NumericPrimary{
						Function: &grammar.NumericValueFunction{
							Length: f.LengthExpression,
						},
					},
				},
			},
		},
	}
}

// DerivedColumn returns the `*grammar.DerivedColumn` element representing
// the Projection
func (f *LengthExpression) DerivedColumn() *grammar.DerivedColumn {
	dc := &grammar.DerivedColumn{
		ValueExpression: grammar.ValueExpression{
			Common: f.CommonValueExpression(),
		},
	}
	if f.alias != "" {
		dc.As = &f.alias
	}
	return dc
}

// As aliases the SQL function as the supplied column name
func (f *LengthExpression) As(alias string) types.Projection {
	f.alias = alias
	return f
}

// Using modifies the CHAR_LENGTH function with a character length units.
func (f *LengthExpression) Using(
	using grammar.CharacterLengthUnits,
) *LengthExpression {
	if f.LengthExpression.Character == nil {
		return f
	}
	f.LengthExpression.Character.Using = using
	return f
}
