package sqlb

type exprType int

const (
    EXP_EQUAL = iota
    EXP_NEQUAL
    EXP_AND
    EXP_OR
    EXP_IN
    EXP_BETWEEN
)

type exprScanInfo []Symbol

var (
    // A static table containing information used in constructing the
    // expression's SQL string during Scan() calls
    exprScanTable = map[exprType]exprScanInfo{
        EXP_EQUAL: exprScanInfo{
            SYM_ELEMENT, SYM_EQUAL, SYM_ELEMENT,
        },
        EXP_NEQUAL: exprScanInfo{
            SYM_ELEMENT, SYM_NEQUAL, SYM_ELEMENT,
        },
        EXP_AND: exprScanInfo{
            SYM_ELEMENT, SYM_AND, SYM_ELEMENT,
        },
        EXP_OR: exprScanInfo{
            SYM_ELEMENT, SYM_OR, SYM_ELEMENT,
        },
        EXP_IN: exprScanInfo{
            SYM_ELEMENT, SYM_IN, SYM_ELEMENT, SYM_RPAREN,
        },
        EXP_BETWEEN: exprScanInfo{
            SYM_ELEMENT, SYM_BETWEEN, SYM_ELEMENT, SYM_AND, SYM_ELEMENT,
        },
    }
)

type Expression struct {
    scanInfo exprScanInfo
    elements []Element
}

func (e *Expression) ArgCount() int {
    ac := 0
    for _, el := range e.elements {
        ac += el.ArgCount()
    }
    return ac
}

func (e *Expression) Size() int {
    size := 0
    elidx := 0
    for _, sym := range e.scanInfo {
        if sym == SYM_ELEMENT {
            el := e.elements[elidx]
            elidx++
            size += el.Size()
        } else {
            size += len(Symbols[sym])
        }
    }
    return size
}

func (e *Expression) Scan(b []byte, args []interface{}) (int, int) {
    bw, ac := 0, 0
    elidx := 0
    for _, sym := range e.scanInfo {
        if sym == SYM_ELEMENT {
            el := e.elements[elidx]
            elidx++
            ebw, eac := el.Scan(b[bw:], args[ac:])
            bw += ebw
            ac += eac
        } else {
            bw += copy(b[bw:], Symbols[sym])
        }
    }
    return bw, ac
}

func Equal(left interface{}, right interface{}) *Expression {
    els := toElements(left, right)
    return &Expression{
        scanInfo: exprScanTable[EXP_EQUAL],
        elements: els,
    }
}

func NotEqual(left interface{}, right interface{}) *Expression {
    els := toElements(left, right)
    return &Expression{
        scanInfo: exprScanTable[EXP_NEQUAL],
        elements: els,
    }
}

func And(a *Expression, b *Expression) *Expression {
    return &Expression{
        scanInfo: exprScanTable[EXP_AND],
        elements: []Element{a, b},
    }
}

func Or(a *Expression, b *Expression) *Expression {
    return &Expression{
        scanInfo: exprScanTable[EXP_OR],
        elements: []Element{a, b},
    }
}

func In(subject Element, values ...interface{}) *Expression {
    return &Expression{
        scanInfo: exprScanTable[EXP_IN],
        elements: []Element{subject, toValueList(values...)},
    }
}

func Between(a *Expression, b *Expression) *Expression {
    return &Expression{
        scanInfo: exprScanTable[EXP_BETWEEN],
        elements: []Element{a, b},
    }
}
