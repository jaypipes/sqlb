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
	From    FromClause
	Where   *WhereClause
	GroupBy *GroupByClause
	Having  *HavingClause
}

func (e *TableExpression) ArgCount(count *int) {
	e.From.ArgCount(count)
	if e.Where != nil {
		e.Where.ArgCount(count)
	}
	if e.GroupBy != nil {
		e.GroupBy.ArgCount(count)
	}
	if e.Having != nil {
		e.Having.ArgCount(count)
	}
}
