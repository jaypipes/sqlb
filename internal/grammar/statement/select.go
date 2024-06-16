//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package statement

import (
	"strings"

	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/clause"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
	"github.com/jaypipes/sqlb/internal/scanner"
)

type Select struct {
	projs   []scanner.Projection
	from    *clause.From
	where   *clause.Where
	groupBy *clause.GroupBy
	having  *clause.Having
	orderBy *clause.OrderBy
	limit   *clause.Limit
}

func (st *Select) Projections() []scanner.Projection {
	return st.projs
}

func (st *Select) Selections() []scanner.Selection {
	return st.from.Selections()
}

func (st *Select) Joins() []*clause.Join {
	return st.from.Joins()
}

func (st *Select) AddProjection(p scanner.Projection) {
	st.projs = append(st.projs, p)
}

func (st *Select) ReplaceSelections(sels []scanner.Selection) {
	st.from.ReplaceSelections(sels)
}

func (st *Select) ArgCount() int {
	argc := 0
	for _, p := range st.projs {
		argc += p.ArgCount()
	}
	if st.from != nil {
		argc += st.from.ArgCount()
	}
	if st.where != nil {
		argc += st.where.ArgCount()
	}
	if st.groupBy != nil {
		argc += st.groupBy.ArgCount()
	}
	if st.having != nil {
		argc += st.having.ArgCount()
	}
	if st.orderBy != nil {
		argc += st.orderBy.ArgCount()
	}
	if st.limit != nil {
		argc += st.limit.ArgCount()
	}
	return argc
}

func (st *Select) Size(s *scanner.Scanner) int {
	size := len(grammar.Symbols[grammar.SYM_SELECT])
	nprojs := len(st.projs)
	for _, p := range st.projs {
		size += p.Size(s)
	}
	size += (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (nprojs - 1)) // the commas...
	if st.from != nil {
		if st.from.Size(s) > 0 {
			size += len(s.Format.SeparateClauseWith)
			size += st.from.Size(s)
		}
	}
	if st.where != nil {
		size += st.where.Size(s)
	}
	if st.groupBy != nil {
		size += st.groupBy.Size(s)
	}
	if st.having != nil {
		size += st.having.Size(s)
	}
	if st.orderBy != nil {
		size += st.orderBy.Size(s)
	}
	if st.limit != nil {
		size += st.limit.Size(s)
	}
	return size
}

func (st *Select) Scan(s *scanner.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	b.Write(grammar.Symbols[grammar.SYM_SELECT])
	nprojs := len(st.projs)
	for x, p := range st.projs {
		p.Scan(s, b, args, curArg)
		if x != (nprojs - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
	if st.from != nil {
		if st.from.Size(s) > 0 {
			b.WriteString(s.Format.SeparateClauseWith)
			st.from.Scan(s, b, args, curArg)
		}
	}
	if st.where != nil {
		st.where.Scan(s, b, args, curArg)
	}
	if st.groupBy != nil {
		st.groupBy.Scan(s, b, args, curArg)
	}
	if st.having != nil {
		st.having.Scan(s, b, args, curArg)
	}
	if st.orderBy != nil {
		st.orderBy.Scan(s, b, args, curArg)
	}
	if st.limit != nil {
		st.limit.Scan(s, b, args, curArg)
	}
}

func (st *Select) AddJoin(jc *clause.Join) *Select {
	st.from.AddJoin(jc)
	return st
}

func (st *Select) AddWhere(e *expression.Expression) *Select {
	if st.where == nil {
		st.where = clause.NewWhere(e)
		return st
	}
	st.where.AddExpression(e)
	return st
}

// Given one or more columns, either set or add to the GROUP BY clause for
// the Select
func (st *Select) AddGroupBy(cols ...scanner.Projection) *Select {
	if len(cols) == 0 {
		return st
	}
	if st.groupBy == nil {
		st.groupBy = clause.NewGroupBy(cols...)
		return st
	}
	for _, c := range cols {
		st.groupBy.AddColumn(c)
	}
	return st
}

func (st *Select) AddHaving(e *expression.Expression) *Select {
	if st.having == nil {
		st.having = clause.NewHaving(e)
		return st
	}
	st.having.AddCondition(e)
	return st
}

// Given one or more sort columns, either set or add to the ORDER BY clause for
// the Select
func (st *Select) AddOrderBy(sortCols ...scanner.Sortable) *Select {
	if len(sortCols) == 0 {
		return st
	}
	if st.orderBy == nil {
		st.orderBy = clause.NewOrderBy(sortCols...)
		return st
	}

	for _, sc := range sortCols {
		st.orderBy.AddSortColumn(sc)
	}
	return st
}

func (st *Select) SetLimitWithOffset(limit int, offset int) *Select {
	tmpOffset := offset
	lc := clause.NewLimit(limit, &tmpOffset)
	st.limit = lc
	return st
}

func (st *Select) SetLimit(limit int) *Select {
	lc := clause.NewLimit(limit, nil)
	st.limit = lc
	return st
}

func (st *Select) RemoveSelection(toRemove scanner.Selection) {
	st.from.RemoveSelection(toRemove)
}

// NewSelect returns a new Select struct that scans into a
// SELECT SQL statement.
func NewSelect(
	projs []scanner.Projection,
	selections []scanner.Selection,
	joins []*clause.Join,
	where *clause.Where,
	groupBy *clause.GroupBy,
	having *clause.Having,
	orderBy *clause.OrderBy,
	limit *clause.Limit,
) *Select {
	return &Select{
		projs:   projs,
		from:    clause.NewFrom(selections, joins),
		where:   where,
		groupBy: groupBy,
		having:  having,
		orderBy: orderBy,
		limit:   limit,
	}
}
