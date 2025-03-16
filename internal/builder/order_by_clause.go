//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
	"github.com/jaypipes/sqlb/core/grammar/symbol"
)

func (b *Builder) doOrderByClause(
	el *grammar.OrderByClause,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(b.opts.FormatSeparateClauseWith())
	b.WriteString(symbol.Order)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.By)
	b.WriteString(symbol.Space)
	for x, ss := range el.SortSpecifications {
		if x > 0 {
			b.WriteString(symbol.Comma)
			b.WriteString(symbol.Space)
		}
		b.doSortSpecification(&ss, qargs, curarg)
	}
}
