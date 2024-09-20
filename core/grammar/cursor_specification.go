//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <cursor specification>    ::=   <query expression> [ <order by clause> ] [ <limit clause> ] [ <updatability clause> ]

type CursorSpecification struct {
	QueryExpression QueryExpression
	OrderByClause   *OrderByClause
	LimitClause     *LimitClause
}
