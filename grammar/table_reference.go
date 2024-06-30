//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <table reference>    ::=   <table primary or joined table> [ <sample clause> ]
//
// <table primary or joined table>    ::=   <table primary> | <joined table>
//
// <sample clause>    ::=
//          TABLESAMPLE <sample method> <left paren> <sample percentage> <right paren> [ <repeatable clause> ]
//
// <sample method>    ::=   BERNOULLI | SYSTEM
//
// <repeatable clause>    ::=   REPEATABLE <left paren> <repeat argument> <right paren>
//
// <sample percentage>    ::=   <numeric value expression>
//
// <repeat argument>    ::=   <numeric value expression>
//
// <table primary>    ::=
//          <table or query name> [ [ AS ] <correlation name> [ <left paren> <derived column list> <right paren> ] ]
//      |     <derived table> [ AS ] <correlation name> [ <left paren> <derived column list> <right paren> ]
//      |     <lateral derived table> [ AS ] <correlation name> [ <left paren> <derived column list> <right paren> ]
//      |     <collection derived table> [ AS ] <correlation name> [ <left paren> <derived column list> <right paren> ]
//      |     <table function derived table> [ AS ] <correlation name> [ <left paren> <derived column list> <right paren> ]
//      |     <only spec> [ [ AS ] <correlation name> [ <left paren> <derived column list> <right paren> ] ]
//      |     <left paren> <joined table> <right paren>
//
// <only spec>    ::=   ONLY <left paren> <table or query name> <right paren>
//
// <lateral derived table>    ::=   LATERAL <table subquery>
//
// <collection derived table>    ::=   UNNEST <left paren> <collection value expression> <right paren> [ WITH ORDINALITY ]
//
// <table function derived table>    ::=   TABLE <left paren> <collection value expression> <right paren>
//
// <derived table>    ::=   <table subquery>
//
// <table or query name>    ::=   <table name> | <query name>
//
// <derived column list>    ::=   <column name list>
//
// <column name list>    ::=   <column name> [ { <comma> <column name> }... ]

// TableReference represents the <table reference> SQL grammar element
type TableReference struct{}
