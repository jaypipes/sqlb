package sqlb

type joinType int

const (
    JOIN_INNER joinType = iota
    JOIN_OUTER
    JOIN_CROSS
)

type joinClause struct {
    joinType joinType
    left selection
    right selection
    on *Expression
}

func (j *joinClause) argCount() int {
    ac := 0
    if j.on != nil {
        ac = j.on.argCount()
    }
    return ac + j.left.argCount() + j.right.argCount()
}

func (j *joinClause) size() int {
    size := 0
    switch j.joinType {
        case JOIN_INNER:
            size += len(Symbols[SYM_JOIN])
        case JOIN_OUTER:
            size += len(Symbols[SYM_LEFT_JOIN])
        case JOIN_CROSS:
            size += len(Symbols[SYM_CROSS_JOIN])
            // CROSS JOIN has no ON condition so just short-circuit here
            return size + j.right.size()
    }
    size += j.right.size()
    size += len(Symbols[SYM_ON])
    size += j.on.size()
    return size
}

func (j *joinClause) scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    switch j.joinType {
        case JOIN_INNER:
            bw += copy(b[bw:], Symbols[SYM_JOIN])
        case JOIN_OUTER:
            bw += copy(b[bw:], Symbols[SYM_LEFT_JOIN])
        case JOIN_CROSS:
            bw += copy(b[bw:], Symbols[SYM_CROSS_JOIN])
    }
    pbw, pac := j.right.scan(b[bw:], args)
    bw += pbw
    ac += pac
    if j.on != nil {
        bw += copy(b[bw:], Symbols[SYM_ON])
        fbw, fac := j.on.scan(b[bw:], args[ac:])
        bw += fbw
        ac += fac
    }
    return bw, ac
}

func Join(left selection, right selection, on *Expression) *joinClause {
    return &joinClause{left: left, right: right, on: on}
}

func OuterJoin(left selection, right selection, on *Expression) *joinClause {
    return &joinClause{
        joinType: JOIN_OUTER,
        left: left,
        right: right,
        on: on,
    }
}

func CrossJoin(left selection, right selection) *joinClause {
    return &joinClause{joinType: JOIN_CROSS, left: left, right: right}
}
