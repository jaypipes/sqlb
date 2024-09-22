//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <datetime value expression>    ::=
//          <datetime term>
//      |     <interval value expression> <plus sign> <datetime term>
//      |     <datetime value expression> <plus sign> <interval term>
//      |     <datetime value expression> <minus sign> <interval term>
//
// <datetime term>    ::=   <datetime factor>
//
// <datetime factor>    ::=   <datetime primary> [ <time zone> ]
//
// <datetime primary>    ::=   <value expression primary> | <datetime value function>
//
// <time zone>    ::=   AT <time zone specifier>
//
// <time zone specifier>    ::=   LOCAL | TIME ZONE <interval primary>

type DatetimeValueExpression struct {
	Unary       *DatetimeTerm
	AddInterval *AddIntervalExpression
	AddSubtract *AddSubtractDatetimeExpression
}

func (e *DatetimeValueExpression) ArgCount(count *int) {
	if e.Unary != nil {
		e.Unary.ArgCount(count)
	} else if e.AddInterval != nil {
		e.AddInterval.ArgCount(count)
	} else if e.AddSubtract != nil {
		e.AddSubtract.ArgCount(count)
	}
}

type AddIntervalExpression struct {
	Left  IntervalValueExpression
	Right DatetimeTerm
}

func (e *AddIntervalExpression) ArgCount(count *int) {
	e.Left.ArgCount(count)
	e.Right.ArgCount(count)
}

type AddSubtractDatetimeExpression struct {
	Left     DatetimeValueExpression
	Right    IntervalTerm
	Subtract bool
}

func (e *AddSubtractDatetimeExpression) ArgCount(count *int) {
	e.Left.ArgCount(count)
	e.Right.ArgCount(count)
}

type DatetimeTerm struct {
	Factor DatetimeFactor
}

func (t *DatetimeTerm) ArgCount(count *int) {
	t.Factor.ArgCount(count)
}

type DatetimeFactor struct {
	Primary  DatetimePrimary
	Timezone *TimezoneSpecifier
}

func (f *DatetimeFactor) ArgCount(count *int) {
	f.Primary.ArgCount(count)
	if f.Timezone != nil {
		f.Timezone.ArgCount(count)
	}
}

type DatetimePrimary struct {
	Primary  *ValueExpressionPrimary
	Function *DatetimeValueFunction
}

func (p *DatetimePrimary) ArgCount(count *int) {
	if p.Primary != nil {
		p.Primary.ArgCount(count)
	}
}

type TimezoneSpecifier struct {
	Local    bool
	Timezone *IntervalPrimary
}

func (s *TimezoneSpecifier) ArgCount(count *int) {
	if s.Timezone != nil {
		s.Timezone.ArgCount(count)
	}
}
