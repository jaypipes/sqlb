//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
	"strings"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/grammar"
)

// GroupBy represents the SQL GROUP BY clause
type GroupBy struct {
	cols []api.Projection
}

func (gb *GroupBy) Columns() []api.Projection {
	return gb.cols
}

func (gb *GroupBy) AddColumn(c api.Projection) {
	gb.cols = append(gb.cols, c)
}

func (gb *GroupBy) ArgCount() int {
	return 0
}

func (gb *GroupBy) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	b.WriteString(opts.FormatSeparateClauseWith())
	b.Write(grammar.Symbols[grammar.SYM_GROUP_BY])
	ncols := len(gb.cols)
	for x, c := range gb.cols {
		reset := c.DisableAliasScan()
		defer reset()
		b.WriteString(c.String(opts, qargs, curarg))
		if x != (ncols - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
	return b.String()
}

// NewGroupBy returns a new GroupBy across one or more projections
func NewGroupBy(cols ...api.Projection) *GroupBy {
	return &GroupBy{
		cols: cols,
	}
}
