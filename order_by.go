//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/types"
)

type sortColumn struct {
	p    types.Projection
	desc bool
}

func (sc *sortColumn) ArgCount() int {
	return sc.p.ArgCount()
}

func (sc *sortColumn) Size(scanner types.Scanner) int {
	reset := sc.p.DisableAliasScan()
	defer reset()
	size := sc.p.Size(scanner)
	if sc.desc {
		size += len(grammar.Symbols[grammar.SYM_DESC])
	}
	return size
}

func (sc *sortColumn) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	reset := sc.p.DisableAliasScan()
	defer reset()
	bw := 0
	bw += sc.p.Scan(scanner, b[bw:], args, curArg)
	if sc.desc {
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_DESC])
	}
	return bw
}

type orderByClause struct {
	scols []*sortColumn
}

func (ob *orderByClause) ArgCount() int {
	argc := 0
	return argc
}

func (ob *orderByClause) Size(scanner types.Scanner) int {
	size := 0
	size += len(scanner.FormatOptions().SeparateClauseWith)
	size += len(grammar.Symbols[grammar.SYM_ORDER_BY])
	ncols := len(ob.scols)
	for _, sc := range ob.scols {
		size += sc.Size(scanner)
	}
	return size + (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (ncols - 1)) // the commas...
}

func (ob *orderByClause) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], scanner.FormatOptions().SeparateClauseWith)
	bw += copy(b[bw:], grammar.Symbols[grammar.SYM_ORDER_BY])
	ncols := len(ob.scols)
	for x, sc := range ob.scols {
		bw += sc.Scan(scanner, b[bw:], args, curArg)
		if x != (ncols - 1) {
			bw += copy(b[bw:], grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
	return bw
}

func (c *Column) Desc() *sortColumn {
	return &sortColumn{p: c, desc: true}
}

func (c *Column) Asc() *sortColumn {
	return &sortColumn{p: c}
}

func (f *sqlFunc) Desc() *sortColumn {
	return &sortColumn{p: f, desc: true}
}

func (f *sqlFunc) Asc() *sortColumn {
	return &sortColumn{p: f}
}
