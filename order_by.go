package sqlb

type sortColumn struct {
    p projection
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

func (sc *sortColumn) scan(b []byte, args []interface{}) (int, int) {
    reset := sc.p.disableAliasScan()
    defer reset()
    var bw, ac int
    ebw, eac := sc.p.scan(b[bw:], args[ac:])
    bw += ebw
    ac += eac
    if sc.desc {
        bw += copy(b[bw:], Symbols[SYM_DESC])
    }
    return bw, ac
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
    return size + (len(Symbols[SYM_COMMA_WS]) * (ncols - 1))  // the commas...
}

func (ob *orderByClause) scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    bw += copy(b[bw:], Symbols[SYM_ORDER_BY])
    ncols := len(ob.scols)
    for x, sc := range ob.scols {
        ebw, eac := sc.scan(b[bw:], args[ac:])
        bw += ebw
        if x != (ncols - 1) {
            bw += copy(b[bw:], Symbols[SYM_COMMA_WS])
        }
        ac += eac
    }
    return bw, ac
}

func (c *Column) Desc() *sortColumn {
    return &sortColumn{p: c, desc: true}
}

func (c *Column) Asc() *sortColumn {
    return &sortColumn{p: c}
}

func (c *ColumnDef) Desc() *sortColumn {
    return &sortColumn{p: c, desc: true}
}

func (c *ColumnDef) Asc() *sortColumn {
    return &sortColumn{p: c}
}
