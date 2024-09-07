//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doOrderByClause(
	el *grammar.OrderByClause,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(b.opts.FormatSeparateClauseWith())
	b.Write(grammar.Symbols[grammar.SYM_ORDER_BY])
	for x, ss := range el.SortSpecifications {
		if x > 0 {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
		b.doSortSpecification(&ss, qargs, curarg)
	}
}
