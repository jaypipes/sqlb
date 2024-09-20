//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <query expression>    ::=   [ <with clause> ] <query expression body>
//
// <with clause>    ::=   WITH [ RECURSIVE ] <with list>
//
// <with list>    ::=   <with list element> [ { <comma> <with list element> }... ]
//
// <with list element>    ::=
//          <query name> [ <left paren> <with column list> <right paren> ]
//          AS <left paren> <query expression> <right paren> [ <search or cycle clause> ]
//
// <with column list>    ::=   <column name list>
//
// <query expression body>    ::=   <non-join query expression> | <joined table>

type QueryExpression struct {
	// WithClause *WithClause
	Body QueryExpressionBody
}

type QueryExpressionBody struct {
	NonJoinQueryExpression *NonJoinQueryExpression
	JoinedTable            *JoinedTable
}

// <non-join query expression>    ::=
//
//	    <non-join query term>
//	|     <query expression body> UNION [ ALL | DISTINCT ] [ <corresponding spec> ] <query term>
//	|     <query expression body> EXCEPT [ ALL | DISTINCT ] [ <corresponding spec> ] <query term>
//
// <query term>    ::=   <non-join query term> | <joined table>
//
// <non-join query term>    ::=
//
//	    <non-join query primary>
//	|     <query term> INTERSECT [ ALL | DISTINCT ] [ <corresponding spec> ] <query primary>
//
// <query primary>    ::=   <non-join query primary> | <joined table>
//
// <non-join query primary>    ::=   <simple table> | <left paren> <non-join query expression> <right paren>
//
// <simple table>    ::=
//
//	    <query specification>
//	|     <table value constructor>
//	|     <explicit table>
//
// <explicit table>    ::=   TABLE <table or query name>
//
// <corresponding spec>    ::=   CORRESPONDING [ BY <left paren> <corresponding column list> <right paren> ]
//
// <corresponding column list>    ::=   <column name list>

type NonJoinQueryExpression struct {
	NonJoinQueryTerm *NonJoinQueryTerm
	// Union *UnionQueryExpression
	// Except *ExceptQueryExpression
}

type NonJoinQueryTerm struct {
	Primary   *NonJoinQueryPrimary
	Intersect *IntersectQuery
}

type NonJoinQueryPrimary struct {
	SimpleTable                         *SimpleTable
	ParenthesizedNonJoinQueryExpression *NonJoinQueryExpression
}

type SimpleTable struct {
	QuerySpecification *QuerySpecification
	//TableValueConstructor *TableValueConstructor
	//ExplicitTable string
}

type QueryTerm struct {
	NonJoinQueryTerm *NonJoinQueryTerm
	JoinedTable      *JoinedTable
}

type QueryPrimary struct {
	NonJoinQueryPrimary *NonJoinQueryPrimary
	JoinedTable         *JoinedTable
}

type IntersectModifier int

const (
	IntersectModifierAll IntersectModifier = iota
	IntersectModifierDistinct
)

type IntersectQuery struct {
	Term     QueryTerm
	Modifier IntersectModifier
	Primary  QueryPrimary
}