package sqlb

// A Literal is a concrete struct wrapper around a constant that implements the
// Scannable interface. Typically, users won't directly construct Literal
// structs but instead helper functions like sqlb.Equal() will construct a
// Literal and bind it to the containing Element.
type Literal struct {
    value interface{}
}

func (lit *Literal) ArgCount() int {
    return 1
}

func (lit  *Literal) Size() int {
    return 1  // The literal is always injected as a question mark
}

func (lit *Literal) Scan(b []byte, args []interface{}) (int, int) {
    args[0] = lit.value
    copy(b, SYM_QM)
    return 1, 1
}
