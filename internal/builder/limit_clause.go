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

func (b *Builder) doLimitClause(
	el *grammar.LimitClause,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(b.opts.FormatSeparateClauseWith())
	b.WriteString(symbol.Limit)
	b.WriteString(symbol.Space)
	b.WriteString(InterpolationMarker(b.opts, *curarg))
	qargs[*curarg] = el.Count
	*curarg++
	if el.Offset != nil {
		b.WriteString(symbol.Space)
		b.WriteString(symbol.Offset)
		b.WriteString(symbol.Space)
		b.WriteString(InterpolationMarker(b.opts, *curarg))
		qargs[*curarg] = *el.Offset
		*curarg++
	}
}
