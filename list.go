//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import "github.com/jaypipes/sqlb/pkg/types"

// A List is a concrete struct wrapper around an array of elements that
// implements the element interface.
type List struct {
	elements []types.Element
}

func (l *List) ArgCount() int {
	ac := 0
	for _, el := range l.elements {
		ac += el.ArgCount()
	}
	return ac
}

func (l *List) Size(scanner types.Scanner) int {
	nels := len(l.elements)
	size := 0
	for _, el := range l.elements {
		size += el.Size(scanner)
	}
	return size + (len(Symbols[SYM_COMMA_WS]) * (nels - 1)) // the commas...
}

func (l *List) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	nels := len(l.elements)
	for x, el := range l.elements {
		bw += el.Scan(scanner, b[bw:], args, curArg)
		if x != (nels - 1) {
			bw += copy(b[bw:], Symbols[SYM_COMMA_WS])
		}
	}
	return bw
}
