//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ast

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/types"
)

type SelectStatement struct {
	projs   []types.Projection
	from    *FromClause
	where   *WhereClause
	groupBy *GroupByClause
	having  *HavingClause
	orderBy *OrderByClause
	limit   *LimitClause
}

func (s *SelectStatement) Projections() []types.Projection {
	return s.projs
}

func (s *SelectStatement) Selections() []types.Selection {
	return s.from.Selections()
}

func (s *SelectStatement) Joins() []*JoinClause {
	return s.from.Joins()
}

func (s *SelectStatement) AddProjection(p types.Projection) {
	s.projs = append(s.projs, p)
}

func (s *SelectStatement) ReplaceSelections(sels []types.Selection) {
	s.from.ReplaceSelections(sels)
}

func (s *SelectStatement) ArgCount() int {
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

func (s *SelectStatement) Size(scanner types.Scanner) int {
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

func (s *SelectStatement) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
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

func (s *SelectStatement) AddJoin(jc *JoinClause) *SelectStatement {
	s.from.AddJoin(jc)
	return s
}

func (s *SelectStatement) AddWhere(e *Expression) *SelectStatement {
	if s.where == nil {
		s.where = NewWhereClause(e)
		return s
	}
	s.where.AddExpression(e)
	return s
}

// Given one or more columns, either set or add to the GROUP BY clause for
// the SelectStatement
func (s *SelectStatement) AddGroupBy(cols ...types.Projection) *SelectStatement {
	if len(cols) == 0 {
		return s
	}
	if s.groupBy == nil {
		s.groupBy = NewGroupByClause(cols...)
		return s
	}
	for _, c := range cols {
		s.groupBy.AddColumn(c)
	}
	return s
}

func (s *SelectStatement) AddHaving(e *Expression) *SelectStatement {
	if s.having == nil {
		s.having = NewHavingClause(e)
		return s
	}
	s.having.AddCondition(e)
	return s
}

// Given one or more sort columns, either set or add to the ORDER BY clause for
// the SelectStatement
func (s *SelectStatement) AddOrderBy(sortCols ...*SortColumn) *SelectStatement {
	if len(sortCols) == 0 {
		return s
	}
	if s.orderBy == nil {
		s.orderBy = NewOrderByClause(sortCols...)
		return s
	}

	for _, sc := range sortCols {
		s.orderBy.AddSortColumn(sc)
	}
	return s
}

func (s *SelectStatement) SetLimitWithOffset(limit int, offset int) *SelectStatement {
	tmpOffset := offset
	lc := NewLimitClause(limit, &tmpOffset)
	s.limit = lc
	return s
}

func (s *SelectStatement) SetLimit(limit int) *SelectStatement {
	lc := NewLimitClause(limit, nil)
	s.limit = lc
	return s
}

func (s *SelectStatement) RemoveSelection(toRemove types.Selection) {
	s.from.RemoveSelection(toRemove)
}

// NewSelectStatement returns a new SelectStatement struct that scans into a
// SELECT SQL statement.
func NewSelectStatement(
	projs []types.Projection,
	selections []types.Selection,
	joins []*JoinClause,
	where *WhereClause,
	groupBy *GroupByClause,
	having *HavingClause,
	orderBy *OrderByClause,
	limit *LimitClause,
) *SelectStatement {
	return &SelectStatement{
		projs:   projs,
		from:    NewFromClause(selections, joins),
		where:   where,
		groupBy: groupBy,
		having:  having,
		orderBy: orderBy,
		limit:   limit,
	}
}
