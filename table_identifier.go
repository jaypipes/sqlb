//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"sort"

	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/schema"
	"github.com/jaypipes/sqlb/pkg/types"
)

type TableIdentifier struct {
	st      *schema.Table
	alias   string
	name    string
	columns []*ColumnIdentifier
}

// Return a pointer to a ColumnIdentifier with a name or alias matching the supplied
// string, or nil if no such column is known
func (t *TableIdentifier) C(name string) *ColumnIdentifier {
	for _, c := range t.columns {
		if c.name == name || c.alias == name {
			return c
		}
	}
	return nil
}

func (t *TableIdentifier) Projections() []types.Projection {
	res := make([]types.Projection, len(t.columns))
	for x, c := range t.columns {
		res[x] = c
	}
	return res
}

func (t *TableIdentifier) ArgCount() int {
	return 0
}

func (t *TableIdentifier) Size(scanner types.Scanner) int {
	size := len(t.name)
	if t.alias != "" {
		size += len(grammar.Symbols[grammar.SYM_AS]) + len(t.alias)
	}
	return size
}

func (t *TableIdentifier) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := copy(b, t.name)
	if t.alias != "" {
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_AS])
		bw += copy(b[bw:], t.alias)
	}
	return bw
}

func (t *TableIdentifier) As(alias string) *TableIdentifier {
	cols := make([]*ColumnIdentifier, len(t.columns))
	tbl := &TableIdentifier{
		st:    t.st,
		alias: alias,
		name:  t.name,
	}
	for x, c := range t.columns {
		cols[x] = &ColumnIdentifier{
			alias: c.alias,
			name:  c.name,
			tbl:   tbl,
		}
	}
	tbl.columns = cols
	return tbl
}

// T returns a TableIdentifier of a given name from a supplied Schema
func T(s *schema.Schema, name string) *TableIdentifier {
	st := s.T(name)
	if st == nil {
		return nil
	}
	ti := &TableIdentifier{
		st:   st,
		name: name,
	}
	cols := make([]*ColumnIdentifier, len(st.Columns))
	colNames := make([]string, len(st.Columns))
	x := 0
	for cname, _ := range st.Columns {
		colNames[x] = cname
		x++
	}
	sort.Strings(colNames)
	for x, cname := range colNames {
		cols[x] = &ColumnIdentifier{
			tbl:  ti,
			name: cname,
		}
	}

	ti.columns = cols
	return ti
}
