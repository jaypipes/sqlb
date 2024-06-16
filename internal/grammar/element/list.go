//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package element

import (
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar"
)

// List is a concrete struct wrapper around an array of elements that
// implements the element interface.
type List struct {
	elements []builder.Element
}

func (l *List) ArgCount() int {
	ac := 0
	for _, el := range l.elements {
		ac += el.ArgCount()
	}
	return ac
}

func (l *List) Size(b *builder.Builder) int {
	nels := len(l.elements)
	size := 0
	for _, el := range l.elements {
		size += el.Size(b)
	}
	return size + (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (nels - 1)) // the commas...
}

func (l *List) Scan(b *builder.Builder, args []interface{}, curArg *int) {
	nels := len(l.elements)
	for x, el := range l.elements {
		el.Scan(b, args, curArg)
		if x != (nels - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
}

// NewList returns a new List struct containing zero or more elements.
func NewList(els ...builder.Element) *List {
	return &List{elements: els}
}
