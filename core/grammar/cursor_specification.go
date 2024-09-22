//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <cursor specification>    ::=   <query expression> [ <order by clause> ] [ <limit clause> ] [ <updatability clause> ]

type CursorSpecification struct {
	Query   QueryExpression
	OrderBy *OrderByClause
	Limit   *LimitClause
}

func (s *CursorSpecification) ArgCount(count *int) {
	s.Query.ArgCount(count)
	if s.OrderBy != nil {
		s.OrderBy.ArgCount(count)
	}
	if s.Limit != nil {
		s.Limit.ArgCount(count)
	}
}
