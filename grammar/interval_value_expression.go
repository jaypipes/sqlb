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
	Unary *IntervalTerm
}

type IntervalTerm struct {
}

type IntervalFactor struct {
}

type IntervalPrimary struct {
}
