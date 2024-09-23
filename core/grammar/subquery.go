//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <subquery>    ::=   <left paren> <query expression> <right paren>

type Subquery struct {
	QueryExpression
}

func (q *Subquery) ArgCount(count *int) {
	q.QueryExpression.ArgCount(count)
}
