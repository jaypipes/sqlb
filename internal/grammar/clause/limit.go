//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
	"strings"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar"
)

type Limit struct {
	limit  int
	offset *int
}

func (lc *Limit) ArgCount() int {
	if lc.offset == nil {
		return 1
	}
	return 2
}

func (lc *Limit) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	b.WriteString(opts.FormatSeparateClauseWith())
	b.Write(grammar.Symbols[grammar.SYM_LIMIT])
	b.WriteString(builder.InterpolationMarker(opts, *curarg))
	qargs[*curarg] = lc.limit
	*curarg++
	if lc.offset != nil {
		b.Write(grammar.Symbols[grammar.SYM_OFFSET])
		b.WriteString(builder.InterpolationMarker(opts, *curarg))
		qargs[*curarg] = *lc.offset
		*curarg++
	}
	return b.String()
}

// NewLimit returns a new Limit struct
func NewLimit(limit int, offset *int) *Limit {
	return &Limit{
		limit:  limit,
		offset: offset,
	}
}
