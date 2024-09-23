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
	Primary *NonParenthesizedValueExpressionPrimary
	Common  *CommonValueExpression
	Boolean *BooleanPredicand
	//ExplictRowValueConstructor             *ExplicitRowValueConstructor
}

func (p *RowValuePredicand) ArgCount(count *int) {
	if p.Primary != nil {
		p.Primary.ArgCount(count)
	} else if p.Common != nil {
		p.Common.ArgCount(count)
	} else if p.Boolean != nil {
		p.Boolean.ArgCount(count)
	}
}
