//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import "github.com/jaypipes/sqlb/pkg/types"

type Table struct {
	alias   string
	meta    *Meta
	name    string
	columns []*Column
}

// Return a pointer to a Column with a name or alias matching the supplied
// string, or nil if no such column is known
func (t *Table) C(name string) *Column {
	for _, c := range t.columns {
		if c.name == name || c.alias == name {
			return c
		}
	}
	return nil
}

func (t *Table) NewColumn(name string) *Column {
	c := t.C(name)
	if c != nil {
		return c
	}
	c = &Column{
		name: name,
		tbl:  t,
	}
	t.columns = append(t.columns, c)
	return c
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
	size := len(t.name)
	if t.alias != "" {
		size += len(Symbols[SYM_AS]) + len(t.alias)
	}
	return size
}

func (t *Table) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := copy(b, t.name)
	if t.alias != "" {
		bw += copy(b[bw:], Symbols[SYM_AS])
		bw += copy(b[bw:], t.alias)
	}
	return bw
}

func (t *Table) As(alias string) *Table {
	cols := make([]*Column, len(t.columns))
	tbl := &Table{
		alias: alias,
		name:  t.name,
		meta:  t.meta,
	}
	for x, c := range t.columns {
		cols[x] = &Column{
			alias: c.alias,
			name:  c.name,
			tbl:   tbl,
		}
	}
	tbl.columns = cols
	return tbl
}
