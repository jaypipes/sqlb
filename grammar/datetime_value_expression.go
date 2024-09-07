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

type DatetimeValueExpression struct{}
