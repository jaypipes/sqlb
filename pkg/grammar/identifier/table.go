//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package identifier

import (
	"sort"

	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/schema"
	"github.com/jaypipes/sqlb/pkg/types"
)

// Table identifies a Table in a SQL statement
type Table struct {
	st      *schema.Table
	Alias   string
	Name    string
	columns []*Column
}

// Schema returns a pointer to the underlying Schema
func (t *Table) Schema() *schema.Schema {
	return t.st.Schema
}

// C returns a pointer to a Column with a name or alias matching the supplied
// string, or nil if no such column is known
func (t *Table) C(name string) *Column {
	for _, c := range t.columns {
		if c.Name == name || c.Alias == name {
			return c
		}
	}
	return nil
}

func (t *Table) Projections() []types.Projection {
	res := make([]types.Projection, len(t.columns))
	for x, c := range t.columns {
		res[x] = c
	}
	return res
}

func (t *Table) ArgCount() int {
	return 0
}

func (t *Table) Size(scanner types.Scanner) int {
	size := len(t.Name)
	if t.Alias != "" {
		size += len(grammar.Symbols[grammar.SYM_AS]) + len(t.Alias)
	}
	return size
}

func (t *Table) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := copy(b, t.Name)
	if t.Alias != "" {
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_AS])
		bw += copy(b[bw:], t.Alias)
	}
	return bw
}

func (t *Table) As(alias string) *Table {
	tbl := &Table{
		st:    t.st,
		Alias: alias,
		Name:  t.Name,
	}
	cols := make([]*Column, len(t.columns))
	for x, c := range t.columns {
		cols[x] = &Column{
			Alias: c.Alias,
			Name:  c.Name,
			tbl:   tbl,
		}
	}
	tbl.columns = cols
	return tbl
}

// TableFromSchema returns a Table of a given name from a
// supplied Schema
func TableFromSchema(
	s *schema.Schema,
	name string,
) *Table {
	st := s.T(name)
	if st == nil {
		return nil
	}
	ti := &Table{
		st:   st,
		Name: name,
	}
	cols := make([]*Column, len(st.Columns))
	colNames := make([]string, len(st.Columns))
	x := 0
	for cname, _ := range st.Columns {
		colNames[x] = cname
		x++
	}
	sort.Strings(colNames)
	for x, cname := range colNames {
		cols[x] = &Column{
			tbl:  ti,
			Name: cname,
		}
	}

	ti.columns = cols
	return ti
}
