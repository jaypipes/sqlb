//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
	"strings"

	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/types"
)

// GroupBy represents the SQL GROUP BY clause
type GroupBy struct {
	cols []types.Projection
}

func (gb *GroupBy) Columns() []types.Projection {
	return gb.cols
}

func (gb *GroupBy) AddColumn(c types.Projection) {
	gb.cols = append(gb.cols, c)
}

func (gb *GroupBy) ArgCount() int {
	argc := 0
	return argc
}

func (gb *GroupBy) Size(scanner types.Scanner) int {
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

func (gb *GroupBy) Scan(scanner types.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	b.WriteString(scanner.FormatOptions().SeparateClauseWith)
	b.Write(grammar.Symbols[grammar.SYM_GROUP_BY])
	ncols := len(gb.cols)
	for x, c := range gb.cols {
		reset := c.DisableAliasScan()
		defer reset()
		c.Scan(scanner, b, args, curArg)
		if x != (ncols - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
}

// NewGroupBy returns a new GroupBy across one or more projections
func NewGroupBy(cols ...types.Projection) *GroupBy {
	return &GroupBy{
		cols: cols,
	}
}
