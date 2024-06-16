//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar"
)

// From represents the SQL FROM clause
type From struct {
	selections []builder.Selection
	joins      []*Join
}

func (f *From) Selections() []builder.Selection {
	return f.selections
}

func (f *From) Joins() []*Join {
	return f.joins
}

func (f *From) ReplaceSelections(sels []builder.Selection) {
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

func (f *From) Size(b *builder.Builder) int {
	size := 0
	nsels := len(f.selections)
	if nsels > 0 {
		size += len(grammar.Symbols[grammar.SYM_FROM])
		for _, sel := range f.selections {
			size += sel.Size(b)
		}
		size += (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (nsels - 1)) // the commas...
		for _, join := range f.joins {
			size += join.Size(b)
		}
	}
	return size
}

func (f *From) Scan(b *builder.Builder, args []interface{}, curArg *int) {
	nsels := len(f.selections)
	if nsels > 0 {
		b.Write(grammar.Symbols[grammar.SYM_FROM])
		for x, sel := range f.selections {
			sel.Scan(b, args, curArg)
			if x != (nsels - 1) {
				b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
			}
		}
		for _, join := range f.joins {
			join.Scan(b, args, curArg)
		}
	}
}

func (f *From) AddJoin(jc *Join) *From {
	f.joins = append(f.joins, jc)
	return f
}

func (f *From) RemoveSelection(toRemove builder.Selection) {
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
	selections []builder.Selection,
	joins []*Join,
) *From {
	return &From{
		selections: selections,
		joins:      joins,
	}
}
