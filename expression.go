package sqlb

type Op int

const (
    OP_EQUAL = iota
    OP_NEQUAL
    OP_AND
    OP_OR
    OP_IN
)

type exprScanInfo struct {
    opSym []byte
    suffix []byte
}

var (
    // A static table containing information used in constructing the
    // expression's SQL string during Scan() calls
    exprScanTable = map[Op]*exprScanInfo{
        OP_EQUAL: &exprScanInfo{opSym: SYM_EQUAL, suffix: SYM_EMPTY},
        OP_NEQUAL: &exprScanInfo{opSym: SYM_NEQUAL, suffix: SYM_EMPTY},
        OP_AND: &exprScanInfo{opSym: SYM_AND, suffix: SYM_EMPTY},
        OP_OR: &exprScanInfo{opSym: SYM_OR, suffix: SYM_EMPTY},
        OP_IN: &exprScanInfo{opSym: SYM_IN, suffix: SYM_RPAREN},
    }
)

type Expression struct {
    scanInfo *exprScanInfo
    left Element
    right Element
}

func (e *Expression) ArgCount() int {
    return e.left.ArgCount() + e.right.ArgCount()
}

func (e *Expression) Size() int {
    return (e.left.Size() +
            len(e.scanInfo.opSym) +
            e.right.Size() +
            len(e.scanInfo.suffix))
}

func (e *Expression) Scan(b []byte, args []interface{}) (int, int) {
    bw, ac := e.left.Scan(b, args)
    bw += copy(b[bw:], e.scanInfo.opSym)
    rbw, rac := e.right.Scan(b[bw:], args[ac:])
    bw += rbw
    ac += rac
    bw += copy(b[bw:], e.scanInfo.suffix)
    return bw, ac
}

func Equal(left interface{}, right interface{}) *Expression {
    els := toElements(left, right)
    return &Expression{
        scanInfo: exprScanTable[OP_EQUAL],
        left: els[0],
        right: els[1],
    }
}

func NotEqual(left interface{}, right interface{}) *Expression {
    els := toElements(left, right)
    return &Expression{
        scanInfo: exprScanTable[OP_NEQUAL],
        left: els[0],
        right: els[1],
    }
}

func And(a *Expression, b *Expression) *Expression {
    return &Expression{
        scanInfo: exprScanTable[OP_AND],
        left: a,
        right: b,
    }
}

func Or(a *Expression, b *Expression) *Expression {
    return &Expression{
        scanInfo: exprScanTable[OP_OR],
        left: a,
        right: b,
    }
}

func In(subject Element, values ...interface{}) *Expression {
    return &Expression{
        scanInfo: exprScanTable[OP_IN],
        left: subject,
        right: toValueList(values...),
    }
}
