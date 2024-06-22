//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package statement

import (
	"strings"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/clause"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
)

type Select struct {
	projs   []api.Projection
	from    *clause.From
	where   *clause.Where
	groupBy *clause.GroupBy
	having  *clause.Having
	orderBy *clause.OrderBy
	limit   *clause.Limit
}

func (st *Select) Projections() []api.Projection {
	return st.projs
}

func (st *Select) Selections() []api.Selection {
	return st.from.Selections()
}

func (st *Select) Joins() []*clause.Join {
	return st.from.Joins()
}

func (st *Select) AddProjection(p api.Projection) {
	st.projs = append(st.projs, p)
}

func (st *Select) ReplaceSelections(sels []api.Selection) {
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

func (st *Select) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	b.Write(grammar.Symbols[grammar.SYM_SELECT])
	nprojs := len(st.projs)
	for x, p := range st.projs {
		b.WriteString(p.String(opts, qargs, curarg))
		if x != (nprojs - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
	if st.from != nil {
		fstr := st.from.String(opts, qargs, curarg)
		if len(fstr) > 0 {
			b.WriteString(opts.FormatSeparateClauseWith())
			b.WriteString(fstr)
		}
	}
	if st.where != nil {
		b.WriteString(st.where.String(opts, qargs, curarg))
	}
	if st.groupBy != nil {
		b.WriteString(st.groupBy.String(opts, qargs, curarg))
	}
	if st.having != nil {
		b.WriteString(st.having.String(opts, qargs, curarg))
	}
	if st.orderBy != nil {
		b.WriteString(st.orderBy.String(opts, qargs, curarg))
	}
	if st.limit != nil {
		b.WriteString(st.limit.String(opts, qargs, curarg))
	}
	return b.String()
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
func (st *Select) AddGroupBy(cols ...api.Projection) *Select {
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
func (st *Select) AddOrderBy(sortCols ...api.Orderable) *Select {
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

func (st *Select) RemoveSelection(toRemove api.Selection) {
	st.from.RemoveSelection(toRemove)
}

// NewSelect returns a new Select struct that scans into a
// SELECT SQL statement.
func NewSelect(
	projs []api.Projection,
	selections []api.Selection,
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
