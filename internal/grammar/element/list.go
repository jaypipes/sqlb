//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package element

import (
	"strings"

	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/scanner"
)

// List is a concrete struct wrapper around an array of elements that
// implements the element interface.
type List struct {
	elements []scanner.Element
}

func (l *List) ArgCount() int {
	ac := 0
	for _, el := range l.elements {
		ac += el.ArgCount()
	}
	return ac
}

func (l *List) Size(s *scanner.Scanner) int {
	nels := len(l.elements)
	size := 0
	for _, el := range l.elements {
		size += el.Size(s)
	}
	return size + (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (nels - 1)) // the commas...
}

func (l *List) Scan(s *scanner.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	nels := len(l.elements)
	for x, el := range l.elements {
		el.Scan(s, b, args, curArg)
		if x != (nels - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
}

// NewList returns a new List struct containing zero or more elements.
func NewList(els ...scanner.Element) *List {
	return &List{elements: els}
}
