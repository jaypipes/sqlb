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
//
// <position expression>    ::=
//          <string position expression>
//      |     <blob position expression>
//
// <string position expression>    ::=
//          POSITION <left paren> <string value expression> IN <string value expression> [ USING <char length units> ] <right paren>
//
// <blob position expression>    ::=
//          POSITION <left paren> <blob value expression> IN <blob value expression> <right paren>
//
// <length expression>    ::=
//          <char length expression>
//      |     <octet length expression>
//
// <char length expression>    ::=
//          { CHAR_LENGTH | CHARACTER_LENGTH } <left paren> <string value expression> [ USING <char length units> ] <right paren>
//
// <octet length expression>    ::=   OCTET_LENGTH <left paren> <string value expression> <right paren>
//
// <extract expression>    ::=   EXTRACT <left paren> <extract field> FROM <extract source> <right paren>
//
// <extract field>    ::=   <primary datetime field> | <time zone field>
//
// <time zone field>    ::=   TIMEZONE_HOUR | TIMEZONE_MINUTE
//
// <extract source>    ::=   <datetime value expression> | <interval value expression>
//
// <cardinality expression>    ::=   CARDINALITY <left paren> <collection value expression> <right paren>
//
// <absolute value expression>    ::=   ABS <left paren> <numeric value expression> <right paren>
//
// <modulus expression>    ::=   MOD <left paren> <numeric value expression dividend> <comma> <numeric value expression divisor> <right paren>
//
// <natural logarithm>    ::=   LN <left paren> <numeric value expression> <right paren>
//
// <exponential function>    ::=   EXP <left paren> <numeric value expression> <right paren>
//
// <power function>    ::=   POWER <left paren> <numeric value expression base> <comma> <numeric value expression exponent> <right paren>
//
// <numeric value expression base>    ::=   <numeric value expression>
//
// <numeric value expression exponent>    ::=   <numeric value expression>
//
// <square root>    ::=   SQRT <left paren> <numeric value expression> <right paren>
//
// <floor function>    ::=   FLOOR <left paren> <numeric value expression> <right paren>
//
// <ceiling function>    ::=   { CEIL | CEILING } <left paren> <numeric value expression> <right paren>
//
// <width bucket function>    ::=   WIDTH_BUCKET <left paren> <width bucket operand> <comma> <width bucket bound 1> <comma> <width bucket bound 2> <comma> <width bucket count> <right paren>
//
// <width bucket operand>    ::=   <numeric value expression>
//
// <width bucket bound 1>    ::=   <numeric value expression>
//
// <width bucket bound 2>    ::=   <numeric value expression>
//
// <width bucket count>    ::=   <numeric value expression>

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

type PositionExpression struct{}

type ExtractExpression struct{}

type LengthExpression struct{}

type CardinalityExpression struct{}

type AbsoluteValueExpression struct{}

type ModulusExpression struct{}

type NaturalLogarithm struct{}

type ExponentialFunction struct{}

type PowerFunction struct{}

type SquareRoot struct{}

type FloorFunction struct{}

type CeilingFunction struct{}

type WidthBucketFunction struct{}
