//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
	"strings"

	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/scanner"
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

func (lc *Limit) Size(s *scanner.Scanner) int {
	// Due to dialect handling, we do not include the length of interpolation
	// markers for query parameters. This is calculated separately by the
	// top-level scanning struct before malloc'ing the buffer to inject the SQL
	// string into.
	size := 0
	size += len(s.Format.SeparateClauseWith)
	size += len(grammar.Symbols[grammar.SYM_LIMIT])
	if lc.offset != nil {
		size += len(grammar.Symbols[grammar.SYM_OFFSET])
	}
	return size
}

func (lc *Limit) Scan(s *scanner.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	b.WriteString(s.Format.SeparateClauseWith)
	b.Write(grammar.Symbols[grammar.SYM_LIMIT])
	scanner.ScanInterpolationMarker(s.Dialect, b, *curArg)
	args[*curArg] = lc.limit
	*curArg++
	if lc.offset != nil {
		b.Write(grammar.Symbols[grammar.SYM_OFFSET])
		scanner.ScanInterpolationMarker(s.Dialect, b, *curArg)
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
