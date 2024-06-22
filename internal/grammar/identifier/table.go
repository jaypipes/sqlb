//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package identifier

import (
	"sort"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar"
)

// Table identifies a Table in a SQL statement
type Table struct {
	Alias   string
	Name    string
	columns []*Column
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
	t *api.Table,
	name string,
) *Table {
	if t == nil {
		return nil
	}
	ti := &Table{
		Name: name,
	}
	cols := make([]*Column, len(t.Columns))
	colNames := make([]string, len(t.Columns))
	x := 0
	for cname := range t.Columns {
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
