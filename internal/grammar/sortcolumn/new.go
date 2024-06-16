//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sortcolumn

import (
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar"
)

// SortColumn describes a column listed in the ORDER BY clause
type SortColumn struct {
	p   builder.Projection
	asc bool
}

func (sc *SortColumn) ArgCount() int {
	return sc.p.ArgCount()
}

func (sc *SortColumn) Size(b *builder.Builder) int {
	reset := sc.p.DisableAliasScan()
	defer reset()
	size := sc.p.Size(b)
	if !sc.asc {
		size += len(grammar.Symbols[grammar.SYM_DESC])
	}
	return size
}

func (sc *SortColumn) Scan(b *builder.Builder, args []interface{}, curArg *int) {
	reset := sc.p.DisableAliasScan()
	defer reset()
	sc.p.Scan(b, args, curArg)
	if !sc.asc {
		b.Write(grammar.Symbols[grammar.SYM_DESC])
	}
}

func (sc *SortColumn) IsAsc() bool {
	return sc.asc
}

// NewAsc returns a new SortColumn on a supplied projection and ascending sort
// order.
func NewAsc(p builder.Projection) builder.Sortable {
	return &SortColumn{
		p:   p,
		asc: true,
	}
}

// NewDesc returns a new SortColumn on a supplied projection and descending
// sort order.
func NewDesc(p builder.Projection) builder.Sortable {
	return &SortColumn{
		p:   p,
		asc: false,
	}
}
