//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

type limitClause struct {
	limit  int
	offset *int
}

func (lc *limitClause) argCount() int {
	if lc.offset == nil {
		return 1
	}
	return 2
}

func (lc *limitClause) size(scanner *sqlScanner) int {
	// Due to dialect handling, we do not include the length of interpolation
	// markers for query parameters. This is calculated separately by the
	// top-level scanning struct before malloc'ing the buffer to inject the SQL
	// string into.
	size := 0
	size += len(scanner.format.SeparateClauseWith)
	size += len(Symbols[SYM_LIMIT])
	if lc.offset != nil {
		size += len(Symbols[SYM_OFFSET])
	}
	return size
}

func (lc *limitClause) scan(scanner *sqlScanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], scanner.format.SeparateClauseWith)
	bw += copy(b[bw:], Symbols[SYM_LIMIT])
	bw += scanInterpolationMarker(scanner.dialect, b[bw:], *curArg)
	args[*curArg] = lc.limit
	*curArg++
	if lc.offset != nil {
		bw += copy(b[bw:], Symbols[SYM_OFFSET])
		bw += scanInterpolationMarker(scanner.dialect, b[bw:], *curArg)
		args[*curArg] = *lc.offset
		*curArg++
	}
	return bw
}
