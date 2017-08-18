package sqlb

type SortColumn struct {
    el Element
    desc bool
}

func (sc *SortColumn) ArgCount() int {
    return sc.el.ArgCount()
}

func (sc *SortColumn) Size() int {
    size := sc.el.Size()
    if sc.desc {
        size += len(Symbols[SYM_DESC])
    }
    return size
}

func (sc *SortColumn) Scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    ebw, eac := sc.el.Scan(b[bw:], args[ac:])
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

func (ob *OrderByClause) ArgCount() int {
    argc := 0
    return argc
}

func (ob *OrderByClause) Size() int {
    size := len(Symbols[SYM_ORDER_BY])
    size += ob.cols.Size()
    return size
}

func (ob *OrderByClause) Scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    bw += copy(b[bw:], Symbols[SYM_ORDER_BY])
    ebw, eac := ob.cols.Scan(b[bw:], args[ac:])
    bw += ebw
    ac += eac
    return bw, ac
}
