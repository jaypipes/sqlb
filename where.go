package sqlb

type whereClause struct {
    filters []*Expression
}

func (w *whereClause) argCount() int {
    argc := 0
    for _, filter := range w.filters {
        argc += filter.argCount()
    }
    return argc
}

func (w *whereClause) size() int {
    size := 0
    nfilters := len(w.filters)
    if nfilters > 0 {
        size += len(Symbols[SYM_WHERE])
        size += len(Symbols[SYM_AND]) * (nfilters - 1)
        for _, filter := range w.filters {
            size += filter.size()
        }
    }
    return size
}

func (w *whereClause) scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    if len(w.filters) > 0 {
        bw += copy(b[bw:], Symbols[SYM_WHERE])
        for x, filter := range w.filters {
            if x > 0 {
                bw += copy(b[bw:], Symbols[SYM_AND])
            }
            fbw, fac := filter.scan(b[bw:], args[ac:])
            bw += fbw
            ac += fac
        }
    }
    return bw, ac
}
