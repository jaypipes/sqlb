//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
	"strings"

	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/scanner"
)

// OrderBy represents the SQL ORDER BY clause
type OrderBy struct {
	scols []scanner.Sortable
}

func (ob *OrderBy) AddSortColumn(sc scanner.Sortable) {
	ob.scols = append(ob.scols, sc)
}

func (ob *OrderBy) ArgCount() int {
	argc := 0
	return argc
}

func (ob *OrderBy) Size(s *scanner.Scanner) int {
	size := 0
	size += len(s.Format.SeparateClauseWith)
	size += len(grammar.Symbols[grammar.SYM_ORDER_BY])
	ncols := len(ob.scols)
	for _, sc := range ob.scols {
		size += sc.Size(s)
	}
	return size + (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (ncols - 1)) // the commas...
}

func (ob *OrderBy) Scan(s *scanner.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	b.WriteString(s.Format.SeparateClauseWith)
	b.Write(grammar.Symbols[grammar.SYM_ORDER_BY])
	ncols := len(ob.scols)
	for x, sc := range ob.scols {
		sc.Scan(s, b, args, curArg)
		if x != (ncols - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
}

// NewOrderBy returns a new OrderBy with zero or more sort columns
func NewOrderBy(scols ...scanner.Sortable) *OrderBy {
	return &OrderBy{
		scols: scols,
	}
}
