//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/grammar/expression"
	"github.com/jaypipes/sqlb/pkg/types"
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

func (w *Where) Size(scanner types.Scanner) int {
	size := 0
	nfilters := len(w.filters)
	if nfilters > 0 {
		size += len(scanner.FormatOptions().SeparateClauseWith)
		size += len(grammar.Symbols[grammar.SYM_WHERE])
		size += len(grammar.Symbols[grammar.SYM_AND]) * (nfilters - 1)
		for _, filter := range w.filters {
			size += filter.Size(scanner)
		}
	}
	return size
}

func (w *Where) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	if len(w.filters) > 0 {
		bw += copy(b[bw:], scanner.FormatOptions().SeparateClauseWith)
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_WHERE])
		for x, filter := range w.filters {
			if x > 0 {
				bw += copy(b[bw:], grammar.Symbols[grammar.SYM_AND])
			}
			bw += filter.Scan(scanner, b[bw:], args, curArg)
		}
	}
	return bw
}

// NewWhere returns a Where populated with zero or more expressions
func NewWhere(exprs ...*expression.Expression) *Where {
	return &Where{
		filters: exprs,
	}
}
