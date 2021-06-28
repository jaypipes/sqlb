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

type GroupByClause struct {
	cols []types.Projection
}

func (gb *GroupByClause) Columns() []types.Projection {
	return gb.cols
}

func (gb *GroupByClause) AddColumn(c types.Projection) {
	gb.cols = append(gb.cols, c)
}

func (gb *GroupByClause) ArgCount() int {
	argc := 0
	return argc
}

func (gb *GroupByClause) Size(scanner types.Scanner) int {
	size := 0
	size += len(scanner.FormatOptions().SeparateClauseWith)
	size += len(grammar.Symbols[grammar.SYM_GROUP_BY])
	ncols := len(gb.cols)
	for _, c := range gb.cols {
		reset := c.DisableAliasScan()
		defer reset()
		size += c.Size(scanner)
	}
	return size + (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (ncols - 1)) // the commas...
}

func (gb *GroupByClause) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], scanner.FormatOptions().SeparateClauseWith)
	bw += copy(b[bw:], grammar.Symbols[grammar.SYM_GROUP_BY])
	ncols := len(gb.cols)
	for x, c := range gb.cols {
		reset := c.DisableAliasScan()
		defer reset()
		bw += c.Scan(scanner, b[bw:], args, curArg)
		if x != (ncols - 1) {
			bw += copy(b[bw:], grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
	return bw
}

// NewGroupByClause returns a new GroupByClause across one or more projections
func NewGroupByClause(cols ...types.Projection) *GroupByClause {
	return &GroupByClause{
		cols: cols,
	}
}
