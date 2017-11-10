package sqlb

type limitClause struct {
	limit   int
	offset  *int
	dialect Dialect
}

func (lc *limitClause) setDialect(dialect Dialect) {
	lc.dialect = dialect
}

func (lc *limitClause) argCount() int {
	if lc.offset == nil {
		return 1
	}
	return 2
}

func (lc *limitClause) size() int {
	// Due to dialect handling, we do not include the length of interpolation
	// markers for query parameters. This is calculated separately by the
	// top-level scanning struct before malloc'ing the buffer to inject the SQL
	// string into.
	size := len(Symbols[SYM_LIMIT])
	if lc.offset != nil {
		size += len(Symbols[SYM_OFFSET])
	}
	return size
}

func (lc *limitClause) scan(b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], Symbols[SYM_LIMIT])
	bw += scanInterpolationMarker(lc.dialect, b[bw:], *curArg)
	args[*curArg] = lc.limit
	*curArg++
	if lc.offset != nil {
		bw += copy(b[bw:], Symbols[SYM_OFFSET])
		bw += scanInterpolationMarker(lc.dialect, b[bw:], *curArg)
		args[*curArg] = *lc.offset
		*curArg++
	}
	return bw
}
