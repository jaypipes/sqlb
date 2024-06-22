//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
	"strings"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/grammar"
)

// From represents the SQL FROM clause
type From struct {
	selections []api.Selection
	joins      []*Join
}

func (f *From) Selections() []api.Selection {
	return f.selections
}

func (f *From) Joins() []*Join {
	return f.joins
}

func (f *From) ReplaceSelections(sels []api.Selection) {
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

func (f *From) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	nsels := len(f.selections)
	if nsels > 0 {
		b.Write(grammar.Symbols[grammar.SYM_FROM])
		for x, sel := range f.selections {
			b.WriteString(sel.String(opts, qargs, curarg))
			if x != (nsels - 1) {
				b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
			}
		}
		for _, join := range f.joins {
			b.WriteString(join.String(opts, qargs, curarg))
		}
	}
	return b.String()
}

func (f *From) AddJoin(jc *Join) *From {
	f.joins = append(f.joins, jc)
	return f
}

func (f *From) RemoveSelection(toRemove api.Selection) {
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
	selections []api.Selection,
	joins []*Join,
) *From {
	return &From{
		selections: selections,
		joins:      joins,
	}
}
