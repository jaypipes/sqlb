package sqlb

type Op int

const (
    OP_EQUAL = iota
    OP_NEQUAL
)

type Expression struct {
    op Op
    left Element
    right Element
}

func (o *Expression) ArgCount() int {
    return o.left.ArgCount() + o.right.ArgCount()
}

func (o *Expression) Size() int {
    return len(SYM_OP[o.op]) + o.left.Size() + o.right.Size()
}

func (o *Expression) Scan(b []byte, args []interface{}) (int, int) {
    idx, argc := o.left.Scan(b, args)
    idx += copy(b[idx:], SYM_OP[o.op])
    bc, ac := o.right.Scan(b[idx:], args)
    idx += bc
    argc += ac
    return idx, argc
}
