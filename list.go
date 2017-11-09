package sqlb

// A List is a concrete struct wrapper around an array of elements that
// implements the element interface.
type List struct {
	elements []element
}

func (l *List) argCount() int {
	ac := 0
	for _, el := range l.elements {
		ac += el.argCount()
	}
	return ac
}

func (l *List) size() int {
	nels := len(l.elements)
	size := 0
	for _, el := range l.elements {
		size += el.size()
	}
	return size + (len(Symbols[SYM_COMMA_WS]) * (nels - 1)) // the commas...
}

func (l *List) scan(b []byte, args []interface{}, curArg *int) int {
	bw := 0
	nels := len(l.elements)
	for x, el := range l.elements {
		bw += el.scan(b[bw:], args, curArg)
		if x != (nels - 1) {
			bw += copy(b[bw:], Symbols[SYM_COMMA_WS])
		}
	}
	return bw
}
