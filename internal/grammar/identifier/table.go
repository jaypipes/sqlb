//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package identifier

import (
	"sort"

	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/meta"
)

// Table identifies a Table in a SQL statement
type Table struct {
	st      *meta.Table
	Alias   string
	Name    string
	columns []*Column
}

// Meta returns a pointer to the underlying Meta
func (t *Table) Meta() *meta.Meta {
	return t.st.Meta
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

func (t *Table) Projections() []builder.Projection {
	res := make([]builder.Projection, len(t.columns))
	for x, c := range t.columns {
		res[x] = c
	}
	return res
}

func (t *Table) ArgCount() int {
	return 0
}

func (t *Table) Size(b *builder.Builder) int {
	size := len(t.Name)
	if t.Alias != "" {
		size += len(grammar.Symbols[grammar.SYM_AS]) + len(t.Alias)
	}
	return size
}

func (t *Table) Scan(b *builder.Builder, args []interface{}, curArg *int) {
	b.WriteString(t.Name)
	if t.Alias != "" {
		b.Write(grammar.Symbols[grammar.SYM_AS])
		b.WriteString(t.Alias)
	}
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

// TableFromMeta returns a Table of a given name from a
// supplied Meta
func TableFromMeta(
	m *meta.Meta,
	name string,
) *Table {
	st := m.T(name)
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
