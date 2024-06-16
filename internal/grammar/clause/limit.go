//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
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

func (lc *Limit) Size(b *builder.Builder) int {
	// Due to dialect handling, we do not include the length of interpolation
	// markers for query parameters. This is calculated separately by the
	// top-level scanning struct before malloc'ing the buffer to inject the SQL
	// string into.
	size := 0
	size += len(b.Format.SeparateClauseWith)
	size += len(grammar.Symbols[grammar.SYM_LIMIT])
	if lc.offset != nil {
		size += len(grammar.Symbols[grammar.SYM_OFFSET])
	}
	return size
}

func (lc *Limit) Scan(b *builder.Builder, args []interface{}, curArg *int) {
	b.WriteString(b.Format.SeparateClauseWith)
	b.Write(grammar.Symbols[grammar.SYM_LIMIT])
	b.AddInterpolationMarker(*curArg)
	args[*curArg] = lc.limit
	*curArg++
	if lc.offset != nil {
		b.Write(grammar.Symbols[grammar.SYM_OFFSET])
		b.AddInterpolationMarker(*curArg)
		args[*curArg] = *lc.offset
		*curArg++
	}
}

// NewLimit returns a new Limit struct
func NewLimit(limit int, offset *int) *Limit {
	return &Limit{
		limit:  limit,
		offset: offset,
	}
}
