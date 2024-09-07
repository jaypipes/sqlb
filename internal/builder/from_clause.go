//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doFromClause(
	el *grammar.FromClause,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(b.opts.FormatSeparateClauseWith())
	b.Write(grammar.Symbols[grammar.SYM_FROM])
	for x, tr := range el.TableReferences {
		if x > 0 {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
		b.doTableReference(&tr, qargs, curarg)
	}
}
