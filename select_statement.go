//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/types"
)

type SelectStatement struct {
	projs      []types.Projection
	selections []types.Selection
	joins      []*JoinClause
	where      *ast.WhereClause
	groupBy    *GroupByClause
	having     *HavingClause
	orderBy    *OrderByClause
	limit      *LimitClause
}

func (s *SelectStatement) ArgCount() int {
	argc := 0
	for _, p := range s.projs {
		argc += p.ArgCount()
	}
	for _, sel := range s.selections {
		argc += sel.ArgCount()
	}
	for _, join := range s.joins {
		argc += join.ArgCount()
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
	nsels := len(s.selections)
	if nsels > 0 {
		size += len(scanner.FormatOptions().SeparateClauseWith)
		size += len(grammar.Symbols[grammar.SYM_FROM])
		for _, sel := range s.selections {
			size += sel.Size(scanner)
		}
		size += (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (nsels - 1)) // the commas...
		for _, join := range s.joins {
			size += join.Size(scanner)
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
	nsels := len(s.selections)
	if nsels > 0 {
		bw += copy(b[bw:], scanner.FormatOptions().SeparateClauseWith)
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_FROM])
		for x, sel := range s.selections {
			bw += sel.Scan(scanner, b[bw:], args, curArg)
			if x != (nsels - 1) {
				bw += copy(b[bw:], grammar.Symbols[grammar.SYM_COMMA_WS])
			}
		}
		for _, join := range s.joins {
			bw += join.Scan(scanner, b[bw:], args, curArg)
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
	s.joins = append(s.joins, jc)
	return s
}

func (s *SelectStatement) AddWhere(e *ast.Expression) *SelectStatement {
	if s.where == nil {
		s.where = ast.NewWhereClause(e)
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
	gb := s.groupBy
	if gb == nil {
		gb = &GroupByClause{
			cols: make([]types.Projection, len(cols)),
		}
		for x, c := range cols {
			gb.cols[x] = c
		}
	} else {
		for _, c := range cols {
			gb.cols = append(gb.cols, c)
		}
	}
	s.groupBy = gb
	return s
}

func (s *SelectStatement) AddHaving(e *ast.Expression) *SelectStatement {
	if s.having == nil {
		s.having = &HavingClause{conditions: make([]*ast.Expression, 0)}
	}
	s.having.conditions = append(s.having.conditions, e)
	return s
}

// Given one or more sort columns, either set or add to the ORDER BY clause for
// the SelectStatement
func (s *SelectStatement) AddOrderBy(sortCols ...*ast.SortColumn) *SelectStatement {
	if len(sortCols) == 0 {
		return s
	}
	ob := s.orderBy
	if ob == nil {
		ob = &OrderByClause{
			scols: make([]*ast.SortColumn, len(sortCols)),
		}
		for x, sc := range sortCols {
			ob.scols[x] = sc
		}
	} else {
		for _, sc := range sortCols {
			ob.scols = append(ob.scols, sc)
		}
	}
	s.orderBy = ob
	return s
}

func (s *SelectStatement) SetLimitWithOffset(limit int, offset int) *SelectStatement {
	lc := &LimitClause{limit: limit}
	lc.offset = &offset
	s.limit = lc
	return s
}

func (s *SelectStatement) SetLimit(limit int) *SelectStatement {
	lc := &LimitClause{limit: limit}
	s.limit = lc
	return s
}

func containsJoin(s *SelectStatement, j *JoinClause) bool {
	for _, sj := range s.joins {
		if j == sj {
			return true
		}
	}
	return false
}

func addToProjections(s *SelectStatement, p types.Projection) {
	s.projs = append(s.projs, p)
}

func (s *SelectStatement) RemoveSelection(toRemove types.Selection) {
	idx := -1
	for x, sel := range s.selections {
		if sel == toRemove {
			idx = x
			break
		}
	}
	if idx == -1 {
		return
	}
	s.selections = append(s.selections[:idx], s.selections[idx+1:]...)
}
