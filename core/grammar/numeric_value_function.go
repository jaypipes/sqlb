//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <numeric value function>    ::=
//          <position expression>
//      |     <extract expression>
//      |     <length expression>
//      |     <cardinality expression>
//      |     <absolute value expression>
//      |     <modulus expression>
//      |     <natural logarithm>
//      |     <exponential function>
//      |     <power function>
//      |     <square root>
//      |     <floor function>
//      |     <ceiling function>
//      |     <width bucket function>

type NumericValueFunction struct {
	Position      *PositionExpression
	Extract       *ExtractExpression
	Length        *LengthExpression
	Cardinality   *CardinalityExpression
	AbsoluteValue *AbsoluteValueExpression
	Modulus       *ModulusExpression
	Natural       *NaturalLogarithm
	Exponential   *ExponentialFunction
	Power         *PowerFunction
	SquareRoot    *SquareRoot
	Floor         *FloorFunction
	Ceiling       *CeilingFunction
	WidthBucket   *WidthBucketFunction
}

func (f *NumericValueFunction) ArgCount(count *int) {
	if f.Position != nil {
		f.Position.ArgCount(count)
	} else if f.Length != nil {
		f.Length.ArgCount(count)
	} else if f.Extract != nil {
		f.Extract.ArgCount(count)
	} else if f.Natural != nil {
		f.Natural.Subject.ArgCount(count)
	} else if f.AbsoluteValue != nil {
		f.AbsoluteValue.Subject.ArgCount(count)
	} else if f.Exponential != nil {
		f.Exponential.Subject.ArgCount(count)
	} else if f.SquareRoot != nil {
		f.SquareRoot.Subject.ArgCount(count)
	} else if f.Ceiling != nil {
		f.Ceiling.Subject.ArgCount(count)
	} else if f.Floor != nil {
		f.Floor.Subject.ArgCount(count)
	}
}

// <position expression>    ::=
//          <string position expression>
//      |     <blob position expression>
//
// <string position expression>    ::=
//          POSITION <left paren> <string value expression> IN <string value expression> [ USING <char length units> ] <right paren>
//
// <blob position expression>    ::=
//          POSITION <left paren> <blob value expression> IN <blob value expression> <right paren>

type PositionExpression struct {
	String *StringPositionExpression
	Blob   *BlobPositionExpression
}

func (e *PositionExpression) ArgCount(count *int) {
	if e.String != nil {
		e.String.Subject.ArgCount(count)
		e.String.In.ArgCount(count)
	} else if e.Blob != nil {
		e.Blob.Subject.ArgCount(count)
		e.Blob.In.ArgCount(count)
	}
}

type StringPositionExpression struct {
	Subject StringValueExpression
	In      StringValueExpression
	Using   CharacterLengthUnits
}

func (e *StringPositionExpression) ArgCount(count *int) {
	e.Subject.ArgCount(count)
	e.In.ArgCount(count)
}

type BlobPositionExpression struct {
	Subject BlobValueExpression
	In      BlobValueExpression
}

func (e *BlobPositionExpression) ArgCount(count *int) {
	e.Subject.ArgCount(count)
	e.In.ArgCount(count)
}

// <extract expression>    ::=   EXTRACT <left paren> <extract field> FROM <extract source> <right paren>
//
// <extract field>    ::=   <primary datetime field> | <time zone field>
//
// <time zone field>    ::=   TIMEZONE_HOUR | TIMEZONE_MINUTE
//
// <extract source>    ::=   <datetime value expression> | <interval value expression>

type TimezoneField int

const (
	TimezoneFieldHour TimezoneField = iota
	TimezoneFieldMinute
)

var TimezoneFieldSymbols = map[TimezoneField]string{
	TimezoneFieldHour:   "TIMEZONE_HOUR",
	TimezoneFieldMinute: "TIMEZONE_MINUTE",
}

type ExtractExpression struct {
	What ExtractField
	From ExtractSource
}

func (e *ExtractExpression) ArgCount(count *int) {
	e.From.ArgCount(count)
}

type ExtractField struct {
	Datetime *PrimaryDatetimeField
	Timezone *TimezoneField
}

type ExtractSource struct {
	Datetime *DatetimeValueExpression
	Interval *IntervalValueExpression
}

func (s *ExtractSource) ArgCount(count *int) {
	if s.Datetime != nil {
		s.Datetime.ArgCount(count)
	} else if s.Interval != nil {
		s.Interval.ArgCount(count)
	}
}

// <length expression>    ::=
//          <char length expression>
//      |     <octet length expression>
//
// <char length expression>    ::=
//          { CHAR_LENGTH | CHARACTER_LENGTH } <left paren> <string value expression> [ USING <char length units> ] <right paren>
//
// <octet length expression>    ::=   OCTET_LENGTH <left paren> <string value expression> <right paren>

type LengthExpression struct {
	Character *CharacterLengthExpression
	Octet     *OctetLengthExpression
}

func (e *LengthExpression) ArgCount(count *int) {
	if e.Character != nil {
		e.Character.Subject.ArgCount(count)
	} else if e.Octet != nil {
		e.Octet.Subject.ArgCount(count)
	}
}

type CharacterLengthExpression struct {
	Subject StringValueExpression
	Using   CharacterLengthUnits
}

func (e *CharacterLengthExpression) ArgCount(count *int) {
	e.Subject.ArgCount(count)
}

type OctetLengthExpression struct {
	Subject StringValueExpression
}

func (e *OctetLengthExpression) ArgCount(count *int) {
	e.Subject.ArgCount(count)
}

// <cardinality expression>    ::=   CARDINALITY <left paren> <collection value expression> <right paren>

type CardinalityExpression struct {
}

func (e *CardinalityExpression) ArgCount(count *int) {

}

// <absolute value expression>    ::=   ABS <left paren> <numeric value expression> <right paren>

type AbsoluteValueExpression struct {
	Subject NumericValueExpression
}

func (e *AbsoluteValueExpression) ArgCount(count *int) {

}

// <modulus expression>    ::=   MOD <left paren> <numeric value expression dividend> <comma> <numeric value expression divisor> <right paren>

type ModulusExpression struct {
	Dividend NumericValueExpression
	Divisor  NumericValueExpression
}

func (e *ModulusExpression) ArgCount(count *int) {

}

// <natural logarithm>    ::=   LN <left paren> <numeric value expression> <right paren>

type NaturalLogarithm struct {
	Subject NumericValueExpression
}

func (l *NaturalLogarithm) ArgCount(count *int) {

}

// <exponential function>    ::=   EXP <left paren> <numeric value expression> <right paren>

type ExponentialFunction struct {
	Subject NumericValueExpression
}

func (f *ExponentialFunction) ArgCount(count *int) {

}

// <power function>    ::=   POWER <left paren> <numeric value expression base> <comma> <numeric value expression exponent> <right paren>
//
// <numeric value expression base>    ::=   <numeric value expression>
//
// <numeric value expression exponent>    ::=   <numeric value expression>

type PowerFunction struct {
	Base     NumericValueExpression
	Exponent NumericValueExpression
}

func (f *PowerFunction) ArgCount(count *int) {

}

// <square root>    ::=   SQRT <left paren> <numeric value expression> <right paren>

type SquareRoot struct {
	Subject NumericValueExpression
}

func (r *SquareRoot) ArgCount(count *int) {

}

// <floor function>    ::=   FLOOR <left paren> <numeric value expression> <right paren>

type FloorFunction struct {
	Subject NumericValueExpression
}

func (f *FloorFunction) ArgCount(count *int) {

}

// <ceiling function>    ::=   { CEIL | CEILING } <left paren> <numeric value expression> <right paren>

type CeilingFunction struct {
	Subject NumericValueExpression
}

func (f *CeilingFunction) ArgCount(count *int) {

}

// <width bucket function>    ::=   WIDTH_BUCKET <left paren> <width bucket operand> <comma> <width bucket bound 1> <comma> <width bucket bound 2> <comma> <width bucket count> <right paren>
//
// <width bucket operand>    ::=   <numeric value expression>
//
// <width bucket bound 1>    ::=   <numeric value expression>
//
// <width bucket bound 2>    ::=   <numeric value expression>
//
// <width bucket count>    ::=   <numeric value expression>

type WidthBucketFunction struct {
	Operand NumericValueExpression
	Bound1  NumericValueExpression
	Bound2  NumericValueExpression
	Count   NumericValueExpression
}

func (f *WidthBucketFunction) ArgCount(count *int) {

}
