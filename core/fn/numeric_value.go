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
// The argument is the subject of the CHAR_LENGTH function and must be
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
	if f.LengthExpression.Character == nil && f.LengthExpression.Octet == nil {
		panic("Cannot set Using on nil LengthExpression")
	}
	if f.LengthExpression.Character == nil {
		return f
	}
	f.LengthExpression.Character.Using = using
	return f
}

// Position returns a PositionExpression that produces a POSITION() SQL function
// that can be passed to sqlb constructs and functions like Select()
//
// The first argument is the subject of the POSITION function and must be
// coercible to either a blob value expression or a string value expression.
// The second argument is the thing to search for the presence of the subject.
// The second argument must also be coercible to either a blob value expression
// or a string value expression.
func Position(
	subjectAny interface{},
	inAny interface{},
) *PositionExpression {
	var ref types.Relation
	switch subjectAny := subjectAny.(type) {
	case types.Projection:
		ref = subjectAny.References()
	}
	switch inAny := inAny.(type) {
	case types.Projection:
		ref = inAny.References()
	}
	blobSubject := inspect.BlobValueExpressionFromAny(subjectAny)
	if blobSubject != nil {
		blobIn := inspect.BlobValueExpressionFromAny(inAny)
		if blobIn == nil {
			msg := fmt.Sprintf(
				"expected coerceable BlobValueExpression but got %+v(%T)",
				inAny, inAny,
			)
			panic(msg)
		}
		return &PositionExpression{
			BaseFunction: BaseFunction{
				ref: ref,
			},
			PositionExpression: &grammar.PositionExpression{
				Blob: &grammar.BlobPositionExpression{
					Subject: *blobSubject,
					In:      *blobIn,
				},
			},
		}
	}
	strSubject := inspect.StringValueExpressionFromAny(subjectAny)
	if strSubject == nil {
		msg := fmt.Sprintf(
			"expected coerceable BlobValueExpression or StringValueExpression but got %+v(%T)",
			subjectAny, subjectAny,
		)
		panic(msg)
	}
	strIn := inspect.StringValueExpressionFromAny(inAny)
	if strIn == nil {
		msg := fmt.Sprintf(
			"expected coerceable BlobValueExpression but got %+v(%T)",
			inAny, inAny,
		)
		panic(msg)
	}
	return &PositionExpression{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		PositionExpression: &grammar.PositionExpression{
			String: &grammar.StringPositionExpression{
				Subject: *strSubject,
				In:      *strIn,
			},
		},
	}
}

// PositionExpression wraps the POSITION() SQL function grammar element
type PositionExpression struct {
	BaseFunction
	*grammar.PositionExpression
}

// CommonValueExpression returns the object as a
// `*grammar.CommonValueExpression`
func (f *PositionExpression) CommonValueExpression() *grammar.CommonValueExpression {
	return &grammar.CommonValueExpression{
		Numeric: &grammar.NumericValueExpression{
			Unary: &grammar.Term{
				Unary: &grammar.Factor{
					Primary: grammar.NumericPrimary{
						Function: &grammar.NumericValueFunction{
							Position: f.PositionExpression,
						},
					},
				},
			},
		},
	}
}

// DerivedColumn returns the `*grammar.DerivedColumn` element representing
// the Projection
func (f *PositionExpression) DerivedColumn() *grammar.DerivedColumn {
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
func (f *PositionExpression) As(alias string) types.Projection {
	f.alias = alias
	return f
}

// Using modifies the POSITION function with a character length units.
func (f *PositionExpression) Using(
	using grammar.CharacterLengthUnits,
) *PositionExpression {
	if f.PositionExpression.String == nil {
		if f.PositionExpression.Blob == nil {
			panic("Cannot set Using on nil PositionExpression")
		}
		// Convert the blob value expression into a string value expression
		// since we're specifying character length units.
		subject := inspect.StringValueExpressionFromAny(
			f.PositionExpression.Blob.Subject,
		)
		in := inspect.StringValueExpressionFromAny(
			f.PositionExpression.Blob.In,
		)
		f.PositionExpression = &grammar.PositionExpression{
			String: &grammar.StringPositionExpression{
				Subject: *subject,
				In:      *in,
			},
		}
	}
	f.PositionExpression.String.Using = using
	return f
}
