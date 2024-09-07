//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doLimitClause(
	el *grammar.LimitClause,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(b.opts.FormatSeparateClauseWith())
	b.Write(grammar.Symbols[grammar.SYM_LIMIT])
	b.WriteString(InterpolationMarker(b.opts, *curarg))
	qargs[*curarg] = el.Count
	*curarg++
	if el.Offset != nil {
		b.Write(grammar.Symbols[grammar.SYM_OFFSET])
		b.WriteString(InterpolationMarker(b.opts, *curarg))
		qargs[*curarg] = *el.Offset
		*curarg++
	}
}
