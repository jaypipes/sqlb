//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

type sortColumn struct {
	p    projection
	desc bool
}

func (sc *sortColumn) argCount() int {
	return sc.p.argCount()
}

func (sc *sortColumn) size(scanner *sqlScanner) int {
	reset := sc.p.disableAliasScan()
	defer reset()
	size := sc.p.size(scanner)
	if sc.desc {
		size += len(Symbols[SYM_DESC])
	}
	return size
}

func (sc *sortColumn) scan(scanner *sqlScanner, b []byte, args []interface{}, curArg *int) int {
	reset := sc.p.disableAliasScan()
	defer reset()
	bw := 0
	bw += sc.p.scan(scanner, b[bw:], args, curArg)
	if sc.desc {
		bw += copy(b[bw:], Symbols[SYM_DESC])
	}
	return bw
}

type orderByClause struct {
	scols []*sortColumn
}

func (ob *orderByClause) argCount() int {
	argc := 0
	return argc
}

func (ob *orderByClause) size(scanner *sqlScanner) int {
	size := 0
	size += len(scanner.format.SeparateClauseWith)
	size += len(Symbols[SYM_ORDER_BY])
	ncols := len(ob.scols)
	for _, sc := range ob.scols {
		size += sc.size(scanner)
	}
	return size + (len(Symbols[SYM_COMMA_WS]) * (ncols - 1)) // the commas...
}

func (ob *orderByClause) scan(scanner *sqlScanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], scanner.format.SeparateClauseWith)
	bw += copy(b[bw:], Symbols[SYM_ORDER_BY])
	ncols := len(ob.scols)
	for x, sc := range ob.scols {
		bw += sc.scan(scanner, b[bw:], args, curArg)
		if x != (ncols - 1) {
			bw += copy(b[bw:], Symbols[SYM_COMMA_WS])
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
