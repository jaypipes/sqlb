//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package ast

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/types"
)

// SortColumn describes a column listed in the ORDER BY clause
type SortColumn struct {
	p    types.Projection
	desc bool
}

func (sc *SortColumn) ArgCount() int {
	return sc.p.ArgCount()
}

func (sc *SortColumn) Size(scanner types.Scanner) int {
	reset := sc.p.DisableAliasScan()
	defer reset()
	size := sc.p.Size(scanner)
	if sc.desc {
		size += len(grammar.Symbols[grammar.SYM_DESC])
	}
	return size
}

func (sc *SortColumn) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	reset := sc.p.DisableAliasScan()
	defer reset()
	bw := 0
	bw += sc.p.Scan(scanner, b[bw:], args, curArg)
	if sc.desc {
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_DESC])
	}
	return bw
}

// NewSortColumn returns a new SortColumn on a supplied projection and sort
// order.
func NewSortColumn(p types.Projection, desc bool) *SortColumn {
	return &SortColumn{
		p:    p,
		desc: desc,
	}
}
