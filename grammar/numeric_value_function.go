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
	Square        *SquareRoot
	Floor         *FloorFunction
	Ceiling       *CeilingFunction
	WidthBucket   *WidthBucketFunction
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

type StringPositionExpression struct {
	Subject StringValueExpression
	In      StringValueExpression
	Using   CharacterLengthUnits
}

type BlobPositionExpression struct {
	Subject BlobValueExpression
	In      BlobValueExpression
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

type ExtractField struct {
	Datetime *PrimaryDatetimeField
	Timezone *TimezoneField
}

type ExtractSource struct {
	Datetime *DatetimeValueExpression
	Interval *IntervalValueExpression
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

type CharacterLengthExpression struct {
	Subject StringValueExpression
	Using   CharacterLengthUnits
}

type OctetLengthExpression struct {
	Subject StringValueExpression
}

// <cardinality expression>    ::=   CARDINALITY <left paren> <collection value expression> <right paren>

type CardinalityExpression struct {
}

// <absolute value expression>    ::=   ABS <left paren> <numeric value expression> <right paren>

type AbsoluteValueExpression struct {
	Subject NumericValueExpression
}

// <modulus expression>    ::=   MOD <left paren> <numeric value expression dividend> <comma> <numeric value expression divisor> <right paren>

type ModulusExpression struct {
	Dividend NumericValueExpression
	Divisor  NumericValueExpression
}

// <natural logarithm>    ::=   LN <left paren> <numeric value expression> <right paren>

type NaturalLogarithm struct {
	Subject NumericValueExpression
}

// <exponential function>    ::=   EXP <left paren> <numeric value expression> <right paren>

type ExponentialFunction struct {
	Subject NumericValueExpression
}

// <power function>    ::=   POWER <left paren> <numeric value expression base> <comma> <numeric value expression exponent> <right paren>
//
// <numeric value expression base>    ::=   <numeric value expression>
//
// <numeric value expression exponent>    ::=   <numeric value expression>

type PowerFunction struct {
	Base      NumericValueExpression
	Expontent NumericValueExpression
}

// <square root>    ::=   SQRT <left paren> <numeric value expression> <right paren>

type SquareRoot struct {
	Subject NumericValueExpression
}

// <floor function>    ::=   FLOOR <left paren> <numeric value expression> <right paren>

type FloorFunction struct {
	Subject NumericValueExpression
}

// <ceiling function>    ::=   { CEIL | CEILING } <left paren> <numeric value expression> <right paren>

type CeilingFunction struct {
	Subject NumericValueExpression
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
	Coun    NumericValueExpression
}
