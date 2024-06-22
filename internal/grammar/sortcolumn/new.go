//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sortcolumn

import (
	"strings"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/grammar"
)

// SortColumn describes a column listed in the ORDER BY clause
type SortColumn struct {
	p   api.Projection
	asc bool
}

func (sc *SortColumn) ArgCount() int {
	return sc.p.ArgCount()
}

func (sc *SortColumn) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	reset := sc.p.DisableAliasScan()
	defer reset()
	b.WriteString(sc.p.String(opts, qargs, curarg))
	if !sc.asc {
		b.Write(grammar.Symbols[grammar.SYM_DESC])
	}
	return b.String()
}

func (sc *SortColumn) IsAsc() bool {
	return sc.asc
}

// NewAsc returns a new SortColumn on a supplied projection and ascending sort
// order.
func NewAsc(p api.Projection) api.Orderable {
	return &SortColumn{
		p:   p,
		asc: true,
	}
}

// NewDesc returns a new SortColumn on a supplied projection and descending
// sort order.
func NewDesc(p api.Projection) api.Orderable {
	return &SortColumn{
		p:   p,
		asc: false,
	}
}
