//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar"
)

// GroupBy represents the SQL GROUP BY clause
type GroupBy struct {
	cols []builder.Projection
}

func (gb *GroupBy) Columns() []builder.Projection {
	return gb.cols
}

func (gb *GroupBy) AddColumn(c builder.Projection) {
	gb.cols = append(gb.cols, c)
}

func (gb *GroupBy) ArgCount() int {
	argc := 0
	return argc
}

func (gb *GroupBy) Size(b *builder.Builder) int {
	size := 0
	size += len(b.Format.SeparateClauseWith)
	size += len(grammar.Symbols[grammar.SYM_GROUP_BY])
	ncols := len(gb.cols)
	for _, c := range gb.cols {
		reset := c.DisableAliasScan()
		defer reset()
		size += c.Size(b)
	}
	return size + (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (ncols - 1)) // the commas...
}

func (gb *GroupBy) Scan(b *builder.Builder, args []interface{}, curArg *int) {
	b.WriteString(b.Format.SeparateClauseWith)
	b.Write(grammar.Symbols[grammar.SYM_GROUP_BY])
	ncols := len(gb.cols)
	for x, c := range gb.cols {
		reset := c.DisableAliasScan()
		defer reset()
		c.Scan(b, args, curArg)
		if x != (ncols - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
}

// NewGroupBy returns a new GroupBy across one or more projections
func NewGroupBy(cols ...builder.Projection) *GroupBy {
	return &GroupBy{
		cols: cols,
	}
}
