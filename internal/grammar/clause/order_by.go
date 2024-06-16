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

// OrderBy represents the SQL ORDER BY clause
type OrderBy struct {
	scols []builder.Sortable
}

func (ob *OrderBy) AddSortColumn(sc builder.Sortable) {
	ob.scols = append(ob.scols, sc)
}

func (ob *OrderBy) ArgCount() int {
	argc := 0
	return argc
}

func (ob *OrderBy) Size(b *builder.Builder) int {
	size := 0
	size += len(b.Format.SeparateClauseWith)
	size += len(grammar.Symbols[grammar.SYM_ORDER_BY])
	ncols := len(ob.scols)
	for _, sc := range ob.scols {
		size += sc.Size(b)
	}
	return size + (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (ncols - 1)) // the commas...
}

func (ob *OrderBy) Scan(b *builder.Builder, args []interface{}, curArg *int) {
	b.WriteString(b.Format.SeparateClauseWith)
	b.Write(grammar.Symbols[grammar.SYM_ORDER_BY])
	ncols := len(ob.scols)
	for x, sc := range ob.scols {
		sc.Scan(b, args, curArg)
		if x != (ncols - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
}

// NewOrderBy returns a new OrderBy with zero or more sort columns
func NewOrderBy(scols ...builder.Sortable) *OrderBy {
	return &OrderBy{
		scols: scols,
	}
}
