//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package statement

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/grammar/clause"
	"github.com/jaypipes/sqlb/pkg/grammar/expression"
	"github.com/jaypipes/sqlb/pkg/types"
)

type Select struct {
	projs   []types.Projection
	from    *clause.From
	where   *clause.Where
	groupBy *clause.GroupBy
	having  *clause.Having
	orderBy *clause.OrderBy
	limit   *clause.Limit
}

func (s *Select) Projections() []types.Projection {
	return s.projs
}

func (s *Select) Selections() []types.Selection {
	return s.from.Selections()
}

func (s *Select) Joins() []*clause.Join {
	return s.from.Joins()
}

func (s *Select) AddProjection(p types.Projection) {
	s.projs = append(s.projs, p)
}

func (s *Select) ReplaceSelections(sels []types.Selection) {
	s.from.ReplaceSelections(sels)
}

func (s *Select) ArgCount() int {
	argc := 0
	for _, p := range s.projs {
		argc += p.ArgCount()
	}
	if s.from != nil {
		argc += s.from.ArgCount()
	}
	if s.where != nil {
		argc += s.where.ArgCount()
	}
	if s.groupBy != nil {
		argc += s.groupBy.ArgCount()
	}
	if s.having != nil {
		argc += s.having.ArgCount()
	}
	if s.orderBy != nil {
		argc += s.orderBy.ArgCount()
	}
	if s.limit != nil {
		argc += s.limit.ArgCount()
	}
	return argc
}

func (s *Select) Size(scanner types.Scanner) int {
	size := len(grammar.Symbols[grammar.SYM_SELECT])
	nprojs := len(s.projs)
	for _, p := range s.projs {
		size += p.Size(scanner)
	}
	size += (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (nprojs - 1)) // the commas...
	if s.from != nil {
		if s.from.Size(scanner) > 0 {
			size += len(scanner.FormatOptions().SeparateClauseWith)
			size += s.from.Size(scanner)
		}
	}
	if s.where != nil {
		size += s.where.Size(scanner)
	}
	if s.groupBy != nil {
		size += s.groupBy.Size(scanner)
	}
	if s.having != nil {
		size += s.having.Size(scanner)
	}
	if s.orderBy != nil {
		size += s.orderBy.Size(scanner)
	}
	if s.limit != nil {
		size += s.limit.Size(scanner)
	}
	return size
}

func (s *Select) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], grammar.Symbols[grammar.SYM_SELECT])
	nprojs := len(s.projs)
	for x, p := range s.projs {
		bw += p.Scan(scanner, b[bw:], args, curArg)
		if x != (nprojs - 1) {
			bw += copy(b[bw:], grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
	if s.from != nil {
		if s.from.Size(scanner) > 0 {
			bw += copy(b[bw:], scanner.FormatOptions().SeparateClauseWith)
			bw += s.from.Scan(scanner, b[bw:], args, curArg)
		}
	}
	if s.where != nil {
		bw += s.where.Scan(scanner, b[bw:], args, curArg)
	}
	if s.groupBy != nil {
		bw += s.groupBy.Scan(scanner, b[bw:], args, curArg)
	}
	if s.having != nil {
		bw += s.having.Scan(scanner, b[bw:], args, curArg)
	}
	if s.orderBy != nil {
		bw += s.orderBy.Scan(scanner, b[bw:], args, curArg)
	}
	if s.limit != nil {
		bw += s.limit.Scan(scanner, b[bw:], args, curArg)
	}
	return bw
}

func (s *Select) AddJoin(jc *clause.Join) *Select {
	s.from.AddJoin(jc)
	return s
}

func (s *Select) AddWhere(e *expression.Expression) *Select {
	if s.where == nil {
		s.where = clause.NewWhere(e)
		return s
	}
	s.where.AddExpression(e)
	return s
}

// Given one or more columns, either set or add to the GROUP BY clause for
// the Select
func (s *Select) AddGroupBy(cols ...types.Projection) *Select {
	if len(cols) == 0 {
		return s
	}
	if s.groupBy == nil {
		s.groupBy = clause.NewGroupBy(cols...)
		return s
	}
	for _, c := range cols {
		s.groupBy.AddColumn(c)
	}
	return s
}

func (s *Select) AddHaving(e *expression.Expression) *Select {
	if s.having == nil {
		s.having = clause.NewHaving(e)
		return s
	}
	s.having.AddCondition(e)
	return s
}

// Given one or more sort columns, either set or add to the ORDER BY clause for
// the Select
func (s *Select) AddOrderBy(sortCols ...types.Sortable) *Select {
	if len(sortCols) == 0 {
		return s
	}
	if s.orderBy == nil {
		s.orderBy = clause.NewOrderBy(sortCols...)
		return s
	}

	for _, sc := range sortCols {
		s.orderBy.AddSortColumn(sc)
	}
	return s
}

func (s *Select) SetLimitWithOffset(limit int, offset int) *Select {
	tmpOffset := offset
	lc := clause.NewLimit(limit, &tmpOffset)
	s.limit = lc
	return s
}

func (s *Select) SetLimit(limit int) *Select {
	lc := clause.NewLimit(limit, nil)
	s.limit = lc
	return s
}

func (s *Select) RemoveSelection(toRemove types.Selection) {
	s.from.RemoveSelection(toRemove)
}

// NewSelect returns a new Select struct that scans into a
// SELECT SQL statement.
func NewSelect(
	projs []types.Projection,
	selections []types.Selection,
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
