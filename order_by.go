package sqlb

type sortColumn struct {
    el element
    desc bool
}

func (sc *sortColumn) argCount() int {
    return sc.el.argCount()
}

func (sc *sortColumn) size() int {
    size := sc.el.size()
    if sc.desc {
        size += len(Symbols[SYM_DESC])
    }
    return size
}

func (sc *sortColumn) scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    ebw, eac := sc.el.scan(b[bw:], args[ac:])
    bw += ebw
    ac += eac
    if sc.desc {
        bw += copy(b[bw:], Symbols[SYM_DESC])
    }
    return bw, ac
}

type orderByClause struct {
    cols *List
}

func (ob *orderByClause) argCount() int {
    argc := 0
    return argc
}

func (ob *orderByClause) size() int {
    size := len(Symbols[SYM_ORDER_BY])
    size += ob.cols.size()
    return size
}

func (ob *orderByClause) scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    bw += copy(b[bw:], Symbols[SYM_ORDER_BY])
    ebw, eac := ob.cols.scan(b[bw:], args[ac:])
    bw += ebw
    ac += eac
    return bw, ac
}

func (c *Column) Desc() *sortColumn {
    return &sortColumn{el: c, desc: true}
}

func (c *Column) Asc() *sortColumn {
    return &sortColumn{el: c}
}

func (c *ColumnDef) Desc() *sortColumn {
    return &sortColumn{el: c, desc: true}
}

func (c *ColumnDef) Asc() *sortColumn {
    return &sortColumn{el: c}
}
