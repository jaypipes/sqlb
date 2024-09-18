//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <interval value expression>    ::=
//          <interval term>
//      |     <interval value expression 1> <plus sign> <interval term 1>
//      |     <interval value expression 1> <minus sign> <interval term 1>
//      |     <left paren> <datetime value expression> <minus sign> <datetime term> <right paren> <interval qualifier>
//
// <interval term>    ::=
//          <interval factor>
//      |     <interval term 2> <asterisk> <factor>
//      |     <interval term 2> <solidus> <factor>
//      |     <term> <asterisk> <interval factor>
//
// <interval factor>    ::=   [ <sign> ] <interval primary>
//
// <interval primary>    ::=
//          <value expression primary> [ <interval qualifier> ]
//      |     <interval value function>
//
// <interval value expression 1>    ::=   <interval value expression>
//
// <interval term 1>    ::=   <interval term>
//
// <interval term 2>    ::=   <interval term>

type IntervalValueExpression struct {
	Unary            *IntervalTerm
	AddSubtract      *AddSubtractIntervalExpression
	SubtractDatetime *SubtractDatetimeExpression
}

type AddSubtractIntervalExpression struct {
	Left     IntervalValueExpression
	Right    Term
	Subtract bool
}

type SubtractDatetimeExpression struct {
	Left  DatetimeValueExpression
	Right DatetimeTerm
}

type IntervalTerm struct {
	Unary           *IntervalFactor
	MultiplyDivide  *MultiplyDivideIntervalTerm
	MultiplyNumeric *MultiplyNumericIntervalFactor
}

type IntervalFactor struct {
	Sign    Sign
	Primary IntervalPrimary
}

type MultiplyDivideIntervalTerm struct {
	Left   IntervalTerm
	Right  Factor
	Divide bool
}

type MultiplyNumericIntervalFactor struct {
	Left  Term
	Right IntervalFactor
}

type IntervalPrimary struct {
	Primary   *ValueExpressionPrimary
	Qualifier *IntervalQualifier
	Function  *IntervalValueFunction
}

// <interval qualifier>    ::=
//          <start field> TO <end field>
//      |     <single datetime field>
//
// <start field>    ::=   <non-second primary datetime field> [ <left paren> <interval leading field precision> <right paren> ]
//
// <end field>    ::=
//          <non-second primary datetime field>
//      |     SECOND [ <left paren> <interval fractional seconds precision> <right paren> ]
//
// <single datetime field>    ::=
//          <non-second primary datetime field> [ <left paren> <interval leading field precision> <right paren> ]
//      |     SECOND [ <left paren> <interval leading field precision> [ <comma> <interval fractional seconds precision> ] <right paren> ]
//
// <primary datetime field>    ::=
//          <non-second primary datetime field>
//      |     SECOND
//
// <non-second primary datetime field>    ::=   YEAR | MONTH | DAY | HOUR | MINUTE
//
// <interval fractional seconds precision>    ::=   <unsigned integer>
//
// <interval leading field precision>    ::=   <unsigned integer>

type IntervalQualifier struct {
	Unary    *SingleDatetimeField
	StartEnd *StartEndDatetimeField
}

type SingleDatetimeField struct {
	Nonsecond *NonsecondPrimaryDatetimeField
}

type StartEndDatetimeField struct {
	Start StartField
	End   EndField
}

type NonsecondPrimaryDatetimeType int

const (
	NonsecondPrimaryDatetimeTypeYear NonsecondPrimaryDatetimeType = iota
	NonsecondPrimaryDatetimeTypeMonth
	NonsecondPrimaryDatetimeTypeDay
	NonsecondPrimaryDatetimeTypeHour
	NonsecondPrimaryDatetimeTypeMinute
)

type StartField struct {
	Nonsecond NonsecondPrimaryDatetimeField
}

type EndField struct {
	Nonsecond *NonsecondPrimaryDatetimeField
	Second    *SecondPrimaryDatetimeField
}

type NonsecondPrimaryDatetimeField struct {
	Type      NonsecondPrimaryDatetimeType
	Precision *uint
}

type SecondPrimaryDatetimeField struct {
	Precision           *uint
	FractionalPrecision *uint
}
