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

// OrderBy represents the SQL ORDER BY clause
type OrderBy struct {
	scols []types.Sortable
}

func (ob *OrderBy) AddSortColumn(sc types.Sortable) {
	ob.scols = append(ob.scols, sc)
}

func (ob *OrderBy) ArgCount() int {
	argc := 0
	return argc
}

func (ob *OrderBy) Size(scanner types.Scanner) int {
	size := 0
	size += len(scanner.FormatOptions().SeparateClauseWith)
	size += len(grammar.Symbols[grammar.SYM_ORDER_BY])
	ncols := len(ob.scols)
	for _, sc := range ob.scols {
		size += sc.Size(scanner)
	}
	return size + (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (ncols - 1)) // the commas...
}

func (ob *OrderBy) Scan(scanner types.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	b.WriteString(scanner.FormatOptions().SeparateClauseWith)
	b.Write(grammar.Symbols[grammar.SYM_ORDER_BY])
	ncols := len(ob.scols)
	for x, sc := range ob.scols {
		sc.Scan(scanner, b, args, curArg)
		if x != (ncols - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
}

// NewOrderBy returns a new OrderBy with zero or more sort columns
func NewOrderBy(scols ...types.Sortable) *OrderBy {
	return &OrderBy{
		scols: scols,
	}
}
