package sqlb

type groupByClause struct {
    cols *List
}

func (ob *groupByClause) argCount() int {
    argc := 0
    return argc
}

func (ob *groupByClause) size() int {
    size := len(Symbols[SYM_GROUP_BY])
    size += ob.cols.size()
    return size
}

func (ob *groupByClause) scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    bw += copy(b[bw:], Symbols[SYM_GROUP_BY])
    ebw, eac := ob.cols.scan(b[bw:], args[ac:])
    bw += ebw
    ac += eac
    return bw, ac
}
