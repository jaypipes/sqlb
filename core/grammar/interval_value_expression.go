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

func (e *IntervalValueExpression) ArgCount(count *int) {
	if e.Unary != nil {
		e.Unary.ArgCount(count)
	} else if e.AddSubtract != nil {
		e.AddSubtract.ArgCount(count)
	} else if e.SubtractDatetime != nil {
		e.SubtractDatetime.ArgCount(count)
	}
}

type AddSubtractIntervalExpression struct {
	Left     IntervalValueExpression
	Right    Term
	Subtract bool
}

func (e *AddSubtractIntervalExpression) ArgCount(count *int) {
	e.Left.ArgCount(count)
	e.Right.ArgCount(count)
}

type SubtractDatetimeExpression struct {
	Left  DatetimeValueExpression
	Right DatetimeTerm
}

func (e *SubtractDatetimeExpression) ArgCount(count *int) {
	e.Left.ArgCount(count)
	e.Right.ArgCount(count)
}

type IntervalTerm struct {
	Unary           *IntervalFactor
	MultiplyDivide  *MultiplyDivideIntervalTerm
	MultiplyNumeric *MultiplyNumericIntervalFactor
}

func (t *IntervalTerm) ArgCount(count *int) {
	if t.Unary != nil {
		t.Unary.ArgCount(count)
	} else if t.MultiplyDivide != nil {
		t.MultiplyDivide.ArgCount(count)
	} else if t.MultiplyNumeric != nil {
		t.MultiplyNumeric.ArgCount(count)
	}
}

type IntervalFactor struct {
	Sign    Sign
	Primary IntervalPrimary
}

func (f *IntervalFactor) ArgCount(count *int) {
	f.Primary.ArgCount(count)
}

type MultiplyDivideIntervalTerm struct {
	Left   IntervalTerm
	Right  Factor
	Divide bool
}

func (t *MultiplyDivideIntervalTerm) ArgCount(count *int) {
	t.Left.ArgCount(count)
	t.Right.ArgCount(count)
}

type MultiplyNumericIntervalFactor struct {
	Left  Term
	Right IntervalFactor
}

func (t *MultiplyNumericIntervalFactor) ArgCount(count *int) {
	t.Left.ArgCount(count)
	t.Right.ArgCount(count)
}

type IntervalPrimary struct {
	Primary   *ValueExpressionPrimary
	Qualifier *IntervalQualifier
	Function  *IntervalValueFunction
}

func (p *IntervalPrimary) ArgCount(count *int) {
	if p.Primary != nil {
		p.Primary.ArgCount(count)
	} else if p.Function != nil {
		p.Function.ArgCount(count)
	}
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

type NonsecondPrimaryDatetimeField int

const (
	NonsecondPrimaryDatetimeFieldYear NonsecondPrimaryDatetimeField = iota
	NonsecondPrimaryDatetimeFieldMonth
	NonsecondPrimaryDatetimeFieldDay
	NonsecondPrimaryDatetimeFieldHour
	NonsecondPrimaryDatetimeFieldMinute
)

var NonsecondPrimaryDatetimeFieldSymbols = map[NonsecondPrimaryDatetimeField]string{
	NonsecondPrimaryDatetimeFieldYear:   "YEAR",
	NonsecondPrimaryDatetimeFieldMonth:  "MONTH",
	NonsecondPrimaryDatetimeFieldDay:    "DAY",
	NonsecondPrimaryDatetimeFieldHour:   "HOUR",
	NonsecondPrimaryDatetimeFieldMinute: "MINUTE",
}

type StartField struct {
	Nonsecond NonsecondPrimaryDatetimeField
}

type EndField struct {
	Nonsecond *NonsecondPrimaryDatetimeField
	Second    *SecondPrimaryDatetimeField
}

type SecondPrimaryDatetimeField struct {
	Precision           *uint
	FractionalPrecision *uint
}

type PrimaryDatetimeField struct {
	Nonsecond NonsecondPrimaryDatetimeField
	Second    bool
}
