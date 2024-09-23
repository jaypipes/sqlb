//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// WHERE <search condition>

type WhereClause struct {
	Search BooleanValueExpression
}

func (c *WhereClause) ArgCount(count *int) {
	c.Search.ArgCount(count)
}
