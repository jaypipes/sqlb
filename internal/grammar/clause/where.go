//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
	"strings"

	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
	"github.com/jaypipes/sqlb/internal/scanner"
)

// Where represents the SQL WHERE clause
type Where struct {
	filters []*expression.Expression
}

func (w *Where) AddExpression(e *expression.Expression) {
	w.filters = append(w.filters, e)
}

func (w *Where) ArgCount() int {
	argc := 0
	for _, filter := range w.filters {
		argc += filter.ArgCount()
	}
	return argc
}

func (w *Where) Size(s *scanner.Scanner) int {
	size := 0
	nfilters := len(w.filters)
	if nfilters > 0 {
		size += len(s.Format.SeparateClauseWith)
		size += len(grammar.Symbols[grammar.SYM_WHERE])
		size += len(grammar.Symbols[grammar.SYM_AND]) * (nfilters - 1)
		for _, filter := range w.filters {
			size += filter.Size(s)
		}
	}
	return size
}

func (w *Where) Scan(s *scanner.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	if len(w.filters) > 0 {
		b.WriteString(s.Format.SeparateClauseWith)
		b.Write(grammar.Symbols[grammar.SYM_WHERE])
		for x, filter := range w.filters {
			if x > 0 {
				b.Write(grammar.Symbols[grammar.SYM_AND])
			}
			filter.Scan(s, b, args, curArg)
		}
	}
}

// NewWhere returns a Where populated with zero or more expressions
func NewWhere(exprs ...*expression.Expression) *Where {
	return &Where{
		filters: exprs,
	}
}
