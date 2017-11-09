package sqlb

type sortColumn struct {
	p    projection
	desc bool
}

func (sc *sortColumn) argCount() int {
	return sc.p.argCount()
}

func (sc *sortColumn) size() int {
	reset := sc.p.disableAliasScan()
	defer reset()
	size := sc.p.size()
	if sc.desc {
		size += len(Symbols[SYM_DESC])
	}
	return size
}

func (sc *sortColumn) scan(b []byte, args []interface{}, curArg *int) int {
	reset := sc.p.disableAliasScan()
	defer reset()
	bw := 0
	bw += sc.p.scan(b[bw:], args, curArg)
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

func (ob *orderByClause) size() int {
	size := len(Symbols[SYM_ORDER_BY])
	ncols := len(ob.scols)
	for _, sc := range ob.scols {
		size += sc.size()
	}
	return size + (len(Symbols[SYM_COMMA_WS]) * (ncols - 1)) // the commas...
}

func (ob *orderByClause) scan(b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], Symbols[SYM_ORDER_BY])
	ncols := len(ob.scols)
	for x, sc := range ob.scols {
		bw += sc.scan(b[bw:], args, curArg)
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
