//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <table expression>    ::=
//          <from clause>
//          [ <where clause> ]
//          [ <group by clause> ]
//          [ <having clause> ]
//          [ <window clause> ]

// TableExpression represents a table expression in the SQL
// statement, e.g. "FROM t WHERE a = b"
type TableExpression struct {
	FromClause    FromClause
	WhereClause   *WhereClause
	GroupByClause *GroupByClause
	HavingClause  *HavingClause
}
