package sqlb

type whereClause struct {
	filters []*Expression
}

// Sets the statement's dialect and pushes the dialect down into any of the
// statement's sub-clauses
func (w *whereClause) setDialect(dialect Dialect) {
	for _, filter := range w.filters {
		filter.setDialect(dialect)
	}
}

func (w *whereClause) argCount() int {
	argc := 0
	for _, filter := range w.filters {
		argc += filter.argCount()
	}
	return argc
}

func (w *whereClause) size() int {
	size := 0
	nfilters := len(w.filters)
	if nfilters > 0 {
		size += len(Symbols[SYM_WHERE])
		size += len(Symbols[SYM_AND]) * (nfilters - 1)
		for _, filter := range w.filters {
			size += filter.size()
		}
	}
	return size
}

func (w *whereClause) scan(b []byte, args []interface{}, curArg *int) int {
	bw := 0
	if len(w.filters) > 0 {
		bw += copy(b[bw:], Symbols[SYM_WHERE])
		for x, filter := range w.filters {
			if x > 0 {
				bw += copy(b[bw:], Symbols[SYM_AND])
			}
			bw += filter.scan(b[bw:], args, curArg)
		}
	}
	return bw
}
