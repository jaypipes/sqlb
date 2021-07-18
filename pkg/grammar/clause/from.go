//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/types"
)

// From represents the SQL FROM clause
type From struct {
	selections []types.Selection
	joins      []*Join
}

func (s *From) Selections() []types.Selection {
	return s.selections
}

func (s *From) Joins() []*Join {
	return s.joins
}

func (s *From) ReplaceSelections(sels []types.Selection) {
	s.selections = sels
}

func (s *From) ArgCount() int {
	argc := 0
	for _, sel := range s.selections {
		argc += sel.ArgCount()
	}
	for _, join := range s.joins {
		argc += join.ArgCount()
	}
	return argc
}

func (s *From) Size(scanner types.Scanner) int {
	size := 0
	nsels := len(s.selections)
	if nsels > 0 {
		size += len(grammar.Symbols[grammar.SYM_FROM])
		for _, sel := range s.selections {
			size += sel.Size(scanner)
		}
		size += (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (nsels - 1)) // the commas...
		for _, join := range s.joins {
			size += join.Size(scanner)
		}
	}
	return size
}

func (s *From) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	nsels := len(s.selections)
	if nsels > 0 {
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_FROM])
		for x, sel := range s.selections {
			bw += sel.Scan(scanner, b[bw:], args, curArg)
			if x != (nsels - 1) {
				bw += copy(b[bw:], grammar.Symbols[grammar.SYM_COMMA_WS])
			}
		}
		for _, join := range s.joins {
			bw += join.Scan(scanner, b[bw:], args, curArg)
		}
	}
	return bw
}

func (s *From) AddJoin(jc *Join) *From {
	s.joins = append(s.joins, jc)
	return s
}

func (s *From) RemoveSelection(toRemove types.Selection) {
	idx := -1
	for x, sel := range s.selections {
		if sel == toRemove {
			idx = x
			break
		}
	}
	if idx == -1 {
		return
	}
	s.selections = append(s.selections[:idx], s.selections[idx+1:]...)
}

// NewFrom returns a new From struct that scans into a
// FROM clause.
func NewFrom(
	selections []types.Selection,
	joins []*Join,
) *From {
	return &From{
		selections: selections,
		joins:      joins,
	}
}
