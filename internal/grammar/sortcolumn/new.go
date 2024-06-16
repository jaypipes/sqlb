//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sortcolumn

import (
	"strings"

	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/scanner"
)

// SortColumn describes a column listed in the ORDER BY clause
type SortColumn struct {
	p   scanner.Projection
	asc bool
}

func (sc *SortColumn) ArgCount() int {
	return sc.p.ArgCount()
}

func (sc *SortColumn) Size(s *scanner.Scanner) int {
	reset := sc.p.DisableAliasScan()
	defer reset()
	size := sc.p.Size(s)
	if !sc.asc {
		size += len(grammar.Symbols[grammar.SYM_DESC])
	}
	return size
}

func (sc *SortColumn) Scan(s *scanner.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	reset := sc.p.DisableAliasScan()
	defer reset()
	sc.p.Scan(s, b, args, curArg)
	if !sc.asc {
		b.Write(grammar.Symbols[grammar.SYM_DESC])
	}
}

func (sc *SortColumn) IsAsc() bool {
	return sc.asc
}

// NewAsc returns a new SortColumn on a supplied projection and ascending sort
// order.
func NewAsc(p scanner.Projection) scanner.Sortable {
	return &SortColumn{
		p:   p,
		asc: true,
	}
}

// NewDesc returns a new SortColumn on a supplied projection and descending
// sort order.
func NewDesc(p scanner.Projection) scanner.Sortable {
	return &SortColumn{
		p:   p,
		asc: false,
	}
}
