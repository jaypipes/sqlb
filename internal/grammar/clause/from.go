//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
	"strings"

	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/scanner"
)

// From represents the SQL FROM clause
type From struct {
	selections []scanner.Selection
	joins      []*Join
}

func (f *From) Selections() []scanner.Selection {
	return f.selections
}

func (f *From) Joins() []*Join {
	return f.joins
}

func (f *From) ReplaceSelections(sels []scanner.Selection) {
	f.selections = sels
}

func (f *From) ArgCount() int {
	argc := 0
	for _, sel := range f.selections {
		argc += sel.ArgCount()
	}
	for _, join := range f.joins {
		argc += join.ArgCount()
	}
	return argc
}

func (f *From) Size(s *scanner.Scanner) int {
	size := 0
	nsels := len(f.selections)
	if nsels > 0 {
		size += len(grammar.Symbols[grammar.SYM_FROM])
		for _, sel := range f.selections {
			size += sel.Size(s)
		}
		size += (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (nsels - 1)) // the commas...
		for _, join := range f.joins {
			size += join.Size(s)
		}
	}
	return size
}

func (f *From) Scan(s *scanner.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	nsels := len(f.selections)
	if nsels > 0 {
		b.Write(grammar.Symbols[grammar.SYM_FROM])
		for x, sel := range f.selections {
			sel.Scan(s, b, args, curArg)
			if x != (nsels - 1) {
				b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
			}
		}
		for _, join := range f.joins {
			join.Scan(s, b, args, curArg)
		}
	}
}

func (f *From) AddJoin(jc *Join) *From {
	f.joins = append(f.joins, jc)
	return f
}

func (f *From) RemoveSelection(toRemove scanner.Selection) {
	idx := -1
	for x, sel := range f.selections {
		if sel == toRemove {
			idx = x
			break
		}
	}
	if idx == -1 {
		return
	}
	f.selections = append(f.selections[:idx], f.selections[idx+1:]...)
}

// NewFrom returns a new From struct that scans into a
// FROM clause.
func NewFrom(
	selections []scanner.Selection,
	joins []*Join,
) *From {
	return &From{
		selections: selections,
		joins:      joins,
	}
}
