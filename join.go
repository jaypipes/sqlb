package sqlb

type joinType int

const (
    JOIN_INNER joinType = iota
    JOIN_OUTER
)

type joinClause struct {
    joinType joinType
    left selection
    right selection
    onExprs []*Expression
}

func (j *joinClause) argCount() int {
    argc := 0
    for _, onExpr := range j.onExprs {
        argc += onExpr.argCount()
    }
    return argc
}

func (j *joinClause) size() int {
    size := 0
    switch j.joinType {
        case JOIN_INNER:
            size += len(Symbols[SYM_JOIN])
        case JOIN_OUTER:
            size += len(Symbols[SYM_LEFT_JOIN])
    }
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

func (j *joinClause) scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    switch j.joinType {
        case JOIN_INNER:
            bw += copy(b[bw:], Symbols[SYM_JOIN])
        case JOIN_OUTER:
            bw += copy(b[bw:], Symbols[SYM_LEFT_JOIN])
    }
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

func (j *joinClause) On(onExprs ...*Expression) *joinClause {
    for _, onExpr := range onExprs {
        j.onExprs = append(j.onExprs, onExpr)
    }
    return j
}

func Join(left selection, right selection, onExpr ...*Expression) *joinClause {
    return &joinClause{left: left, right: right, onExprs: onExpr}
}

func OuterJoin(left selection, right selection, onExpr ...*Expression) *joinClause {
    return &joinClause{
        joinType: JOIN_OUTER,
        left: left,
        right: right,
        onExprs: onExpr,
    }
}
