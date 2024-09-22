//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <having clause>    ::=   HAVING <search condition>

type HavingClause struct {
	Search BooleanValueExpression
}

func (c *HavingClause) ArgCount(count *int) {
	c.Search.ArgCount(count)
}
