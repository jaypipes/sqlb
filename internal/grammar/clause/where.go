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
	"github.com/jaypipes/sqlb/internal/grammar/expression"
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

func (w *Where) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	if len(w.filters) > 0 {
		b.WriteString(opts.FormatSeparateClauseWith())
		b.Write(grammar.Symbols[grammar.SYM_WHERE])
		for x, filter := range w.filters {
			if x > 0 {
				b.Write(grammar.Symbols[grammar.SYM_AND])
			}
			b.WriteString(filter.String(opts, qargs, curarg))
		}
	}
	return b.String()
}

// NewWhere returns a Where populated with zero or more expressions
func NewWhere(exprs ...*expression.Expression) *Where {
	return &Where{
		filters: exprs,
	}
}
