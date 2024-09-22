//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <value expression>    ::=
//          <common value expression>
//      |     <boolean value expression>
//      |     <row value expression>
//
// <common value expression>    ::=
//          <numeric value expression>
//      |     <string value expression>
//      |     <datetime value expression>
//      |     <interval value expression>
//      |     <user-defined type value expression>
//      |     <reference value expression>
//      |     <collection value expression>
//
// <user-defined type value expression>    ::=   <value expression primary>
//
// <reference value expression>    ::=   <value expression primary>
//
// <collection value expression>    ::=   <array value expression> | <multiset value expression>
//
// <collection value constructor>    ::=   <array value constructor> | <multiset value constructor>
//
// <row value expression>    ::=
//          <row value special case>
//      |     <explicit row value constructor>
//
// <table row value expression>    ::=
//          <row value special case>
//      |     <row value constructor>
//
// <contextually typed row value expression>    ::=
//          <row value special case>
//      |     <contextually typed row value constructor>
//
// <row value predicand>    ::=
//          <row value special case>
//      |     <row value constructor predicand>
//
// <row value special case>    ::=   <nonparenthesized value expression primary>

type ValueExpression struct {
	Common  *CommonValueExpression
	Boolean *BooleanValueExpression
	Row     *RowValueExpression
}

func (e *ValueExpression) ArgCount(count *int) {
	if e.Common != nil {
		e.Common.ArgCount(count)
	} else if e.Boolean != nil {
		e.Boolean.ArgCount(count)
	} else if e.Row != nil {
		e.Row.Primary.ArgCount(count)
	}
}

type CommonValueExpression struct {
	Numeric  *NumericValueExpression
	String   *StringValueExpression
	Datetime *DatetimeValueExpression
	Interval *IntervalValueExpression
	//UserDefinedTypeValueExpression *UserDefinedTypeValueExpression
	//ReferenceValueExpression       *ReferenceValueExpression
	//CollectionValueExpression      *CollectionValueExpression
}

func (e *CommonValueExpression) ArgCount(count *int) {
	if e.Numeric != nil {
		e.Numeric.ArgCount(count)
	} else if e.String != nil {
		e.String.ArgCount(count)
	} else if e.Datetime != nil {
		e.Datetime.ArgCount(count)
	} else if e.Interval != nil {
		e.Interval.ArgCount(count)
	}
}

type RowValueExpression struct {
	Primary *NonParenthesizedValueExpressionPrimary
}

func (e *RowValueExpression) ArgCount(count *int) {
	e.Primary.ArgCount(count)
}
