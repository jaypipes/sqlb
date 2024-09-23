//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <derived column>    ::=   <value expression> [ <as clause> ]
//
// <as clause>    ::=   [ AS ] <column name>

type DerivedColumn struct {
	Value ValueExpression
	As    *string
}

func (c *DerivedColumn) ArgCount(count *int) {
	c.Value.ArgCount(count)
}
