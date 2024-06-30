//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <group by clause>    ::=   GROUP BY [ <set quantifier> ] <grouping element list>
//
// <grouping element list>    ::=   <grouping element> [ { <comma> <grouping element> }... ]
//
// <grouping element>    ::=
//          <ordinary grouping set>
//      |     <rollup list>
//      |     <cube list>
//      |     <grouping sets specification>
//      |     <empty grouping set>
//
// <ordinary grouping set>    ::=
//          <grouping column reference>
//      |     <left paren> <grouping column reference list> <right paren>
//
// <grouping column reference>    ::=   <column reference> [ <collate clause> ]
//
// <grouping column reference list>    ::=   <grouping column reference> [ { <comma> <grouping column reference> }... ]
//
// <rollup list>    ::=   ROLLUP <left paren> <ordinary grouping set list> <right paren>
//
// <ordinary grouping set list>    ::=   <ordinary grouping set> [ { <comma> <ordinary grouping set> }... ]
//
// <cube list>    ::=   CUBE <left paren> <ordinary grouping set list> <right paren>
//
// <grouping sets specification>    ::=   GROUPING SETS <left paren> <grouping set list> <right paren>
//
// <grouping set list>    ::=   <grouping set> [ { <comma> <grouping set> }... ]
//
// <grouping set>    ::=
//          <ordinary grouping set>
//      |     <rollup list>
//      |     <cube list>
//      |     <grouping sets specification>
//      |     <empty grouping set>
//
// <empty grouping set>    ::=   <left paren> <right paren>

// GroupBy represents the SQL GROUP BY clause
type GroupBy struct {
	Columns []string
}
