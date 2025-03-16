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

func (b *Builder) doFromClause(
	el *grammar.FromClause,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(b.opts.FormatSeparateClauseWith())
	b.WriteString(symbol.From)
	b.WriteString(symbol.Space)
	for x, tr := range el.TableReferences {
		if x > 0 {
			b.WriteString(symbol.Comma)
			b.WriteString(symbol.Space)
		}
		b.doTableReference(&tr, qargs, curarg)
	}
}
