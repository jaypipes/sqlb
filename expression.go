package sqlb

type Op int

const (
    OP_EQUAL = iota
    OP_NEQUAL
    OP_AND
)

type Expression struct {
    op Op
    left Element
    right Element
}

func (e *Expression) ArgCount() int {
    return e.left.ArgCount() + e.right.ArgCount()
}

func (e *Expression) Size() int {
    return len(SYM_OP[e.op]) + e.left.Size() + e.right.Size()
}

func (e *Expression) Scan(b []byte, args []interface{}) (int, int) {
    bw, ac := e.left.Scan(b, args)
    bw += copy(b[bw:], SYM_OP[e.op])
    rbw, rac := e.right.Scan(b[bw:], args[ac:])
    bw += rbw
    ac += rac
    return bw, ac
}

func Equal(left interface{}, right interface{}) *Expression {
    els := toElements(left, right)
    return &Expression{
        op: OP_EQUAL,
        left: els[0],
        right: els[1],
    }
}

func NotEqual(left interface{}, right interface{}) *Expression {
    els := toElements(left, right)
    return &Expression{
        op: OP_NEQUAL,
        left: els[0],
        right: els[1],
    }
}

func And(a *Expression, b *Expression) *Expression {
    return &Expression{op: OP_AND, left: a, right: b}
}

type InExpression struct {
    subject Element
    values *List
}

func (e *InExpression) ArgCount() int {
    return e.values.ArgCount()
}

func (e *InExpression) Size() int {
    return e.subject.Size() + SYM_IN_LEN + e.values.Size() + SYM_RPAREN_LEN
}

func (e *InExpression) Scan(b []byte, args []interface{}) (int, int) {
    bw, ac := e.subject.Scan(b, args)
    bw += copy(b[bw:], SYM_IN)
    lbw, lac := e.values.Scan(b[bw:], args[ac:])
    bw += lbw
    ac += lac
    bw += copy(b[bw:], SYM_RPAREN)
    return bw, ac
}

func In(subject Element, values ...interface{}) *InExpression {
    return &InExpression{
        subject: subject,
        values: toValueList(values...),
    }
}
