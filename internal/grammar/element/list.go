//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package element

import (
	"strings"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/grammar"
)

// List is a concrete struct wrapper around an array of elements that
// implements the element interface.
type List struct {
	elements []api.Element
}

func (l *List) ArgCount() int {
	ac := 0
	for _, el := range l.elements {
		ac += el.ArgCount()
	}
	return ac
}

func (l *List) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	nels := len(l.elements)
	for x, el := range l.elements {
		b.WriteString(el.String(opts, qargs, curarg))
		if x != (nels - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
	return b.String()
}

// NewList returns a new List struct containing zero or more elements.
func NewList(els ...api.Element) *List {
	return &List{elements: els}
}
