//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

type whereClause struct {
	filters []*Expression
}

func (w *whereClause) argCount() int {
	argc := 0
	for _, filter := range w.filters {
		argc += filter.argCount()
	}
	return argc
}

func (w *whereClause) size(scanner *sqlScanner) int {
	size := 0
	nfilters := len(w.filters)
	if nfilters > 0 {
		size += len(scanner.format.SeparateClauseWith)
		size += len(Symbols[SYM_WHERE])
		size += len(Symbols[SYM_AND]) * (nfilters - 1)
		for _, filter := range w.filters {
			size += filter.size(scanner)
		}
	}
	return size
}

func (w *whereClause) scan(scanner *sqlScanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	if len(w.filters) > 0 {
		bw += copy(b[bw:], scanner.format.SeparateClauseWith)
		bw += copy(b[bw:], Symbols[SYM_WHERE])
		for x, filter := range w.filters {
			if x > 0 {
				bw += copy(b[bw:], Symbols[SYM_AND])
			}
			bw += filter.scan(scanner, b[bw:], args, curArg)
		}
	}
	return bw
}
