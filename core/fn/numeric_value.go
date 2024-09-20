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

type ExtractField int

const (
	ExtractFieldYear ExtractField = iota
	ExtractFieldMonth
	ExtractFieldDay
	ExtractFieldHour
	ExtractFieldMinute
	ExtractFieldSecond
	ExtractFieldTimezoneHour
	ExtractFieldTimezoneMinute
)

func grammarExtractField(field ExtractField) *grammar.ExtractField {
	switch field {
	case ExtractFieldSecond:
		return &grammar.ExtractField{
			Datetime: &grammar.PrimaryDatetimeField{
				Second: true,
			},
		}
	case ExtractFieldTimezoneHour:
		tzf := grammar.TimezoneFieldHour
		return &grammar.ExtractField{
			Timezone: &tzf,
		}
	case ExtractFieldTimezoneMinute:
		tzf := grammar.TimezoneFieldMinute
		return &grammar.ExtractField{
			Timezone: &tzf,
		}
	default:
		return &grammar.ExtractField{
			Datetime: &grammar.PrimaryDatetimeField{
				Nonsecond: grammar.NonsecondPrimaryDatetimeField(field),
			},
		}
	}
}

// Extract returns a ExtractExpression that produces a EXTRACT() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The first argument is the subject of the EXTRACT function and must be
// coercible to a datetime value expression or interval value expression. The
// second argument specifies which datetime or timezone field to extract from
// the value expression identified in the first parameter.
func Extract(
	fromAny interface{},
	what ExtractField,
) *ExtractExpression {
	var ref types.Relation
	switch fromAny := fromAny.(type) {
	case types.Projection:
		ref = fromAny.References()
	}
	var source *grammar.ExtractSource
	fromDatetime := inspect.DatetimeValueExpressionFromAny(fromAny)
	if fromDatetime != nil {
		source = &grammar.ExtractSource{Datetime: fromDatetime}
	} else {
		fromInterval := inspect.IntervalValueExpressionFromAny(fromAny)
		if fromInterval == nil {
			msg := fmt.Sprintf(
				"expected coerceable DatetimeValueExpression or IntervalValueExpression but got %+v(%T)",
				fromAny, fromAny,
			)
			panic(msg)
		}
		source = &grammar.ExtractSource{Interval: fromInterval}
	}
	return &ExtractExpression{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		ExtractExpression: &grammar.ExtractExpression{
			From: *source,
			What: *grammarExtractField(what),
		},
	}
}

// ExtractExpression wraps the CHAR_LENGTH() SQL function grammar element
type ExtractExpression struct {
	BaseFunction
	*grammar.ExtractExpression
}

// CommonValueExpression returns the object as a
// `*grammar.CommonValueExpression`
func (f *ExtractExpression) CommonValueExpression() *grammar.CommonValueExpression {
	return &grammar.CommonValueExpression{
		Numeric: &grammar.NumericValueExpression{
			Unary: &grammar.Term{
				Unary: &grammar.Factor{
					Primary: grammar.NumericPrimary{
						Function: &grammar.NumericValueFunction{
							Extract: f.ExtractExpression,
						},
					},
				},
			},
		},
	}
}

// DerivedColumn returns the `*grammar.DerivedColumn` element representing
// the Projection
func (f *ExtractExpression) DerivedColumn() *grammar.DerivedColumn {
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
func (f *ExtractExpression) As(alias string) types.Projection {
	f.alias = alias
	return f
}

// NaturalLogarithm returns a NumericUnaryfunction that produces a LN() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The argument is the subject of the LN function and must be coercible to a
// numeric value expression.
func NaturalLogarithm(
	subjectAny interface{},
) *NumericValueFunction {
	ref, subject := relationsAndSubjectAsNumericValueExpression(
		subjectAny,
	)
	return &NumericValueFunction{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		NumericValueFunction: &grammar.NumericValueFunction{
			Natural: &grammar.NaturalLogarithm{
				Subject: *subject,
			},
		},
	}
}

var Ln = NaturalLogarithm

// Absolute returns a NumericUnaryfunction that produces a ABS() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The argument is the subject of the ABS function and must be coercible to a
// numeric value expression.
func Absolute(
	subjectAny interface{},
) *NumericValueFunction {
	ref, subject := relationsAndSubjectAsNumericValueExpression(
		subjectAny,
	)
	return &NumericValueFunction{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		NumericValueFunction: &grammar.NumericValueFunction{
			AbsoluteValue: &grammar.AbsoluteValueExpression{
				Subject: *subject,
			},
		},
	}
}

var Abs = Absolute

// Exponential returns a NumericUnaryfunction that produces a EXP() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The argument is the subject of the EXP function and must be coercible to a
// numeric value expression.
func Exponential(
	subjectAny interface{},
) *NumericValueFunction {
	ref, subject := relationsAndSubjectAsNumericValueExpression(
		subjectAny,
	)
	return &NumericValueFunction{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		NumericValueFunction: &grammar.NumericValueFunction{
			Exponential: &grammar.ExponentialFunction{
				Subject: *subject,
			},
		},
	}
}

var Exp = Exponential

// SquareRoot returns a NumericUnaryfunction that produces a SQRT() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The argument is the subject of the SQRT function and must be coercible to a
// numeric value expression.
func SquareRoot(
	subjectAny interface{},
) *NumericValueFunction {
	ref, subject := relationsAndSubjectAsNumericValueExpression(
		subjectAny,
	)
	return &NumericValueFunction{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		NumericValueFunction: &grammar.NumericValueFunction{
			SquareRoot: &grammar.SquareRoot{
				Subject: *subject,
			},
		},
	}
}

var SqRt = SquareRoot

// Ceiling returns a NumericUnaryfunction that produces a CEIL() SQL function
// that can be passed to sqlb constructs and functions like Select()
//
// The argument is the subject of the CEIL function and must be coercible to a
// numeric value expression.
func Ceiling(
	subjectAny interface{},
) *NumericValueFunction {
	ref, subject := relationsAndSubjectAsNumericValueExpression(
		subjectAny,
	)
	return &NumericValueFunction{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		NumericValueFunction: &grammar.NumericValueFunction{
			Ceiling: &grammar.CeilingFunction{
				Subject: *subject,
			},
		},
	}
}

var Ceil = Ceiling

// Floor returns a NumericUnaryfunction that produces a FLOOR() SQL function
// that can be passed to sqlb constructs and functions like Select()
//
// The argument is the subject of the FLOOR function and must be coercible to a
// numeric value expression.
func Floor(
	subjectAny interface{},
) *NumericValueFunction {
	ref, subject := relationsAndSubjectAsNumericValueExpression(
		subjectAny,
	)
	return &NumericValueFunction{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		NumericValueFunction: &grammar.NumericValueFunction{
			Floor: &grammar.FloorFunction{
				Subject: *subject,
			},
		},
	}
}

// NumericValueFunction wraps a number of unary numeric value SQL function
// grammar elements
type NumericValueFunction struct {
	BaseFunction
	*grammar.NumericValueFunction
}

func relationsAndSubjectAsNumericValueExpression(
	subjectAny interface{},
) (types.Relation, *grammar.NumericValueExpression) {
	var ref types.Relation
	switch subjectAny := subjectAny.(type) {
	case types.Projection:
		ref = subjectAny.References()
	}
	subject := inspect.NumericValueExpressionFromAny(subjectAny)
	if subject == nil {
		msg := fmt.Sprintf(
			"expected coerceable NumericValueExpression but got %+v(%T)",
			subjectAny, subjectAny,
		)
		panic(msg)
	}
	return ref, subject
}

// CommonValueExpression returns the object as a
// `*grammar.CommonValueExpression`
func (f *NumericValueFunction) CommonValueExpression() *grammar.CommonValueExpression {
	return &grammar.CommonValueExpression{
		Numeric: &grammar.NumericValueExpression{
			Unary: &grammar.Term{
				Unary: &grammar.Factor{
					Primary: grammar.NumericPrimary{
						Function: f.NumericValueFunction,
					},
				},
			},
		},
	}
}

// DerivedColumn returns the `*grammar.DerivedColumn` element representing
// the Projection
func (f *NumericValueFunction) DerivedColumn() *grammar.DerivedColumn {
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
func (f *NumericValueFunction) As(alias string) types.Projection {
	f.alias = alias
	return f
}
