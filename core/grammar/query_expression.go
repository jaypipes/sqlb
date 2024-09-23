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

func (e *QueryExpression) ArgCount(count *int) {
	e.Body.ArgCount(count)
}

type QueryExpressionBody struct {
	NonJoin *NonJoinQueryExpression
	Joined  *JoinedTable
}

func (b *QueryExpressionBody) ArgCount(count *int) {
	if b.NonJoin != nil {
		b.NonJoin.ArgCount(count)
	} else if b.Joined != nil {
		b.Joined.ArgCount(count)
	}
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
	NonJoin *NonJoinQueryTerm
	// Union *UnionQueryExpression
	// Except *ExceptQueryExpression
}

func (e *NonJoinQueryExpression) ArgCount(count *int) {
	if e.NonJoin != nil {
		e.NonJoin.ArgCount(count)
	}
}

type NonJoinQueryTerm struct {
	Primary   *NonJoinQueryPrimary
	Intersect *IntersectQuery
}

func (t *NonJoinQueryTerm) ArgCount(count *int) {
	if t.Primary != nil {
		t.Primary.ArgCount(count)
	}
}

type NonJoinQueryPrimary struct {
	Simple        *SimpleTable
	Parenthesized *NonJoinQueryExpression
}

func (p *NonJoinQueryPrimary) ArgCount(count *int) {
	if p.Simple != nil {
		p.Simple.ArgCount(count)
	} else if p.Parenthesized != nil {
		p.Parenthesized.ArgCount(count)
	}
}

type SimpleTable struct {
	QuerySpecification *QuerySpecification
	//TableValueConstructor *TableValueConstructor
	//ExplicitTable string
}

func (t *SimpleTable) ArgCount(count *int) {
	if t.QuerySpecification != nil {
		t.QuerySpecification.ArgCount(count)
	}
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
