package sqlb

// A List is a concrete struct wrapper around an array of elements that
// implements the Element interface.
type List struct {
    elements []Element
}

func (l *List) ArgCount() int {
    ac := 0
    for _, el := range l.elements {
        ac += el.ArgCount()
    }
    return ac
}

func (l  *List) Size() int {
    nels := len(l.elements)
    size := 0
    for _, el := range l.elements {
        size += el.Size()
    }
    return size + (len(Symbols[SYM_COMMA_WS]) * (nels - 1))  // the commas...
}

func (l *List) Scan(b []byte, args []interface{}) (int, int) {
    bw, ac := 0, 0
    nels := len(l.elements)
    for x, el := range l.elements {
        ebw, eac := el.Scan(b[bw:], args[ac:])
        bw += ebw
        if x != (nels - 1) {
            bw += copy(b[bw:], Symbols[SYM_COMMA_WS])
        }
        ac += eac
    }
    return bw, ac
}
