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

type CommonValueExpression struct {
	Numeric  *NumericValueExpression
	String   *StringValueExpression
	Datetime *DatetimeValueExpression
	Interval *IntervalValueExpression
	//UserDefinedTypeValueExpression *UserDefinedTypeValueExpression
	//ReferenceValueExpression       *ReferenceValueExpression
	//CollectionValueExpression      *CollectionValueExpression
}

type RowValueExpression struct {
	Primary *NonParenthesizedValueExpressionPrimary
}
