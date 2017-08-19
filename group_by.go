package sqlb

type GroupByClause struct {
    cols *List
}

func (ob *GroupByClause) ArgCount() int {
    argc := 0
    return argc
}

func (ob *GroupByClause) Size() int {
    size := len(Symbols[SYM_GROUP_BY])
    size += ob.cols.Size()
    return size
}

func (ob *GroupByClause) Scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    bw += copy(b[bw:], Symbols[SYM_GROUP_BY])
    ebw, eac := ob.cols.Scan(b[bw:], args[ac:])
    bw += ebw
    ac += eac
    return bw, ac
}
