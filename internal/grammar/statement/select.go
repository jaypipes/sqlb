//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package statement

import (
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/clause"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
)

type Select struct {
	projs   []builder.Projection
	from    *clause.From
	where   *clause.Where
	groupBy *clause.GroupBy
	having  *clause.Having
	orderBy *clause.OrderBy
	limit   *clause.Limit
}

func (st *Select) Projections() []builder.Projection {
	return st.projs
}

func (st *Select) Selections() []builder.Selection {
	return st.from.Selections()
}

func (st *Select) Joins() []*clause.Join {
	return st.from.Joins()
}

func (st *Select) AddProjection(p builder.Projection) {
	st.projs = append(st.projs, p)
}

func (st *Select) ReplaceSelections(sels []builder.Selection) {
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

func (st *Select) Size(b *builder.Builder) int {
	size := len(grammar.Symbols[grammar.SYM_SELECT])
	nprojs := len(st.projs)
	for _, p := range st.projs {
		size += p.Size(b)
	}
	size += (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (nprojs - 1)) // the commas...
	if st.from != nil {
		if st.from.Size(b) > 0 {
			size += len(b.Format.SeparateClauseWith)
			size += st.from.Size(b)
		}
	}
	if st.where != nil {
		size += st.where.Size(b)
	}
	if st.groupBy != nil {
		size += st.groupBy.Size(b)
	}
	if st.having != nil {
		size += st.having.Size(b)
	}
	if st.orderBy != nil {
		size += st.orderBy.Size(b)
	}
	if st.limit != nil {
		size += st.limit.Size(b)
	}
	return size
}

func (st *Select) Scan(b *builder.Builder, args []interface{}, curArg *int) {
	b.Write(grammar.Symbols[grammar.SYM_SELECT])
	nprojs := len(st.projs)
	for x, p := range st.projs {
		p.Scan(b, args, curArg)
		if x != (nprojs - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
	if st.from != nil {
		if st.from.Size(b) > 0 {
			b.WriteString(b.Format.SeparateClauseWith)
			st.from.Scan(b, args, curArg)
		}
	}
	if st.where != nil {
		st.where.Scan(b, args, curArg)
	}
	if st.groupBy != nil {
		st.groupBy.Scan(b, args, curArg)
	}
	if st.having != nil {
		st.having.Scan(b, args, curArg)
	}
	if st.orderBy != nil {
		st.orderBy.Scan(b, args, curArg)
	}
	if st.limit != nil {
		st.limit.Scan(b, args, curArg)
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
func (st *Select) AddGroupBy(cols ...builder.Projection) *Select {
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
func (st *Select) AddOrderBy(sortCols ...builder.Sortable) *Select {
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

func (st *Select) RemoveSelection(toRemove builder.Selection) {
	st.from.RemoveSelection(toRemove)
}

// NewSelect returns a new Select struct that scans into a
// SELECT SQL statement.
func NewSelect(
	projs []builder.Projection,
	selections []builder.Selection,
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
