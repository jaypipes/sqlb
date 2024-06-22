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

// OrderBy represents the SQL ORDER BY clause
type OrderBy struct {
	scols []api.Orderable
}

func (ob *OrderBy) AddSortColumn(sc api.Orderable) {
	ob.scols = append(ob.scols, sc)
}

func (ob *OrderBy) ArgCount() int {
	argc := 0
	return argc
}

func (ob *OrderBy) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	b.WriteString(opts.FormatSeparateClauseWith())
	b.Write(grammar.Symbols[grammar.SYM_ORDER_BY])
	ncols := len(ob.scols)
	for x, sc := range ob.scols {
		b.WriteString(sc.String(opts, qargs, curarg))
		if x != (ncols - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
	return b.String()
}

// NewOrderBy returns a new OrderBy with zero or more sort columns
func NewOrderBy(scols ...api.Orderable) *OrderBy {
	return &OrderBy{
		scols: scols,
	}
}
