package sqlb

type joinType int

const (
    JOIN_INNER joinType = iota
)

type JoinClause struct {
    joinType joinType
    left *Table
    right *Table
    onExprs []*Expression
}

func (j *JoinClause) ArgCount() int {
    argc := 0
    for _, onExpr := range j.onExprs {
        argc += onExpr.ArgCount()
    }
    return argc
}

func (j *JoinClause) Size() int {
    size := len(Symbols[SYM_JOIN])
    size += j.right.Size()
    nonExprs := len(j.onExprs)
    if nonExprs > 0 {
        size += len(Symbols[SYM_ON])
        size += len(Symbols[SYM_AND]) * (nonExprs - 1)
        for _, onExpr := range j.onExprs {
            size += onExpr.Size()
        }
    }
    return size
}

func (j *JoinClause) Scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    bw += copy(b[bw:], Symbols[SYM_JOIN])
    pbw, pac := j.right.Scan(b[bw:], args)
    bw += pbw
    ac += pac
    if len(j.onExprs) > 0 {
        bw += copy(b[bw:], Symbols[SYM_ON])
        for x, onExpr := range j.onExprs {
            if x > 0 {
                bw += copy(b[bw:], Symbols[SYM_AND])
            }
            fbw, fac := onExpr.Scan(b[bw:], args[ac:])
            bw += fbw
            ac += fac
        }
    }
    return bw, ac
}
