//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import "github.com/jaypipes/sqlb/pkg/types"

type whereClause struct {
	filters []*Expression
}

func (w *whereClause) ArgCount() int {
	argc := 0
	for _, filter := range w.filters {
		argc += filter.ArgCount()
	}
	return argc
}

func (w *whereClause) Size(scanner types.Scanner) int {
	size := 0
	nfilters := len(w.filters)
	if nfilters > 0 {
		size += len(scanner.FormatOptions().SeparateClauseWith)
		size += len(Symbols[SYM_WHERE])
		size += len(Symbols[SYM_AND]) * (nfilters - 1)
		for _, filter := range w.filters {
			size += filter.Size(scanner)
		}
	}
	return size
}

func (w *whereClause) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	if len(w.filters) > 0 {
		bw += copy(b[bw:], scanner.FormatOptions().SeparateClauseWith)
		bw += copy(b[bw:], Symbols[SYM_WHERE])
		for x, filter := range w.filters {
			if x > 0 {
				bw += copy(b[bw:], Symbols[SYM_AND])
			}
			bw += filter.Scan(scanner, b[bw:], args, curArg)
		}
	}
	return bw
}
