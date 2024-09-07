//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <row value predicand>    ::=
//          <row value special case>
//      |     <row value constructor predicand>
//
// <row value special case>    ::=   <nonparenthesized value expression primary>
//
// <row value constructor predicand>    ::=
//          <common value expression>
//      |     <boolean predicand>
//      |     <explicit row value constructor>

type RowValuePredicand struct {
	NonParenthesizedValueExpressionPrimary *NonParenthesizedValueExpressionPrimary
	CommonValueExpression                  *CommonValueExpression
	BooleanPredicand                       *BooleanPredicand
	//ExplictRowValueConstructor             *ExplicitRowValueConstructor
}
