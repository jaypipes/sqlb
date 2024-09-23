//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <order by clause>    ::=   ORDER BY <sort specification list>
//
// <sort specification list>    ::=   <sort specification> [ { <comma> <sort specification> }... ]

type OrderByClause struct {
	SortSpecifications []SortSpecification
}

func (c *OrderByClause) ArgCount(count *int) {
	for _, ss := range c.SortSpecifications {
		ss.ArgCount(count)
	}
}
