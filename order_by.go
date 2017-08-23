package sqlb

type SortColumn struct {
    el element
    desc bool
}

func (sc *SortColumn) argCount() int {
    return sc.el.argCount()
}

func (sc *SortColumn) size() int {
    size := sc.el.size()
    if sc.desc {
        size += len(Symbols[SYM_DESC])
    }
    return size
}

func (sc *SortColumn) scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    ebw, eac := sc.el.scan(b[bw:], args[ac:])
    bw += ebw
    ac += eac
    if sc.desc {
        bw += copy(b[bw:], Symbols[SYM_DESC])
    }
    return bw, ac
}

type OrderByClause struct {
    cols *List
}

func (ob *OrderByClause) argCount() int {
    argc := 0
    return argc
}

func (ob *OrderByClause) size() int {
    size := len(Symbols[SYM_ORDER_BY])
    size += ob.cols.size()
    return size
}

func (ob *OrderByClause) scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    bw += copy(b[bw:], Symbols[SYM_ORDER_BY])
    ebw, eac := ob.cols.scan(b[bw:], args[ac:])
    bw += ebw
    ac += eac
    return bw, ac
}

func (c *Column) Desc() *SortColumn {
    return &SortColumn{el: c, desc: true}
}

func (c *Column) Asc() *SortColumn {
    return &SortColumn{el: c}
}

func (c *ColumnDef) Desc() *SortColumn {
    return &SortColumn{el: c, desc: true}
}

func (c *ColumnDef) Asc() *SortColumn {
    return &SortColumn{el: c}
}
