package sqlb


type sortColumn struct {
    el Element
    desc bool
}

type OrderByClause struct {
    cols []*sortColumn
}

func (ob *OrderByClause) ArgCount() int {
    argc := 0
    for _, col := range ob.cols {
        argc += col.el.ArgCount()
    }
    return argc
}

func (ob *OrderByClause) Size() int {
    size := len(Symbols[SYM_ORDER_BY])
    for _, col := range ob.cols {
        size += col.el.Size()
        if col.desc {
            size += len(Symbols[SYM_DESC])
        }
    }
    size += ((len(ob.cols) - 1) * len(Symbols[SYM_COMMA_WS]))
    return size
}

func (ob *OrderByClause) Scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    bw += copy(b[bw:], Symbols[SYM_ORDER_BY])
    for x, col := range ob.cols {
        if x > 0 {
            bw += copy(b[bw:], Symbols[SYM_COMMA_WS])
        }
        ebw, eac := col.el.Scan(b[bw:], args[ac:])
        bw += ebw
        ac += eac
        if col.desc {
            bw += copy(b[bw:], Symbols[SYM_DESC])
        }
    }
    return bw, ac
}
