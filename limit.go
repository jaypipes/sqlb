//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	pkgscanner "github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/types"
)

type LimitClause struct {
	limit  int
	offset *int
}

func (lc *LimitClause) ArgCount() int {
	if lc.offset == nil {
		return 1
	}
	return 2
}

func (lc *LimitClause) Size(scanner types.Scanner) int {
	// Due to dialect handling, we do not include the length of interpolation
	// markers for query parameters. This is calculated separately by the
	// top-level scanning struct before malloc'ing the buffer to inject the SQL
	// string into.
	size := 0
	size += len(scanner.FormatOptions().SeparateClauseWith)
	size += len(grammar.Symbols[grammar.SYM_LIMIT])
	if lc.offset != nil {
		size += len(grammar.Symbols[grammar.SYM_OFFSET])
	}
	return size
}

func (lc *LimitClause) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], scanner.FormatOptions().SeparateClauseWith)
	bw += copy(b[bw:], grammar.Symbols[grammar.SYM_LIMIT])
	bw += pkgscanner.ScanInterpolationMarker(scanner.Dialect(), b[bw:], *curArg)
	args[*curArg] = lc.limit
	*curArg++
	if lc.offset != nil {
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_OFFSET])
		bw += pkgscanner.ScanInterpolationMarker(scanner.Dialect(), b[bw:], *curArg)
		args[*curArg] = *lc.offset
		*curArg++
	}
	return bw
}
