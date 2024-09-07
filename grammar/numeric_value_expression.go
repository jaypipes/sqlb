//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <numeric value expression>    ::=
//          <term>
//      |     <numeric value expression> <plus sign> <term>
//      |     <numeric value expression> <minus sign> <term>
//
// <term>    ::=
//          <factor>
//      |     <term> <asterisk> <factor>
//      |     <term> <solidus> <factor>
//
// <factor>    ::=   [ <sign> ] <numeric primary>
//
// <numeric primary>    ::=
//          <value expression primary>
//      |     <numeric value function>

type NumericValueExpression struct{}
