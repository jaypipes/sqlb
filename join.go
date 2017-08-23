package sqlb

type joinType int

const (
    JOIN_INNER joinType = iota
)

type JoinClause struct {
    joinType joinType
    left selection
    right selection
    onExprs []*Expression
}

func (j *JoinClause) argCount() int {
    argc := 0
    for _, onExpr := range j.onExprs {
        argc += onExpr.argCount()
    }
    return argc
}

func (j *JoinClause) size() int {
    size := len(Symbols[SYM_JOIN])
    size += j.right.size()
    nonExprs := len(j.onExprs)
    if nonExprs > 0 {
        size += len(Symbols[SYM_ON])
        size += len(Symbols[SYM_AND]) * (nonExprs - 1)
        for _, onExpr := range j.onExprs {
            size += onExpr.size()
        }
    }
    return size
}

func (j *JoinClause) scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    bw += copy(b[bw:], Symbols[SYM_JOIN])
    pbw, pac := j.right.scan(b[bw:], args)
    bw += pbw
    ac += pac
    if len(j.onExprs) > 0 {
        bw += copy(b[bw:], Symbols[SYM_ON])
        for x, onExpr := range j.onExprs {
            if x > 0 {
                bw += copy(b[bw:], Symbols[SYM_AND])
            }
            fbw, fac := onExpr.scan(b[bw:], args[ac:])
            bw += fbw
            ac += fac
        }
    }
    return bw, ac
}

func (j *JoinClause) On(onExprs ...*Expression) *JoinClause {
    for _, onExpr := range onExprs {
        j.onExprs = append(j.onExprs, onExpr)
    }
    return j
}

func Join(left selection, right selection, onExpr ...*Expression) *JoinClause {
    return &JoinClause{left: left, right: right, onExprs: onExpr}
}
