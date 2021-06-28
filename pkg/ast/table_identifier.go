//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ast

import (
	"sort"

	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/schema"
	"github.com/jaypipes/sqlb/pkg/types"
)

// TableIdentifier identifies a Table in a SQL statement
type TableIdentifier struct {
	st      *schema.Table
	Alias   string
	Name    string
	columns []*ColumnIdentifier
}

// Schema returns a pointer to the underlying Schema
func (t *TableIdentifier) Schema() *schema.Schema {
	return t.st.Schema
}

// C returns a pointer to a ColumnIdentifier with a name or alias matching the supplied
// string, or nil if no such column is known
func (t *TableIdentifier) C(name string) *ColumnIdentifier {
	for _, c := range t.columns {
		if c.Name == name || c.Alias == name {
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
	size := len(t.Name)
	if t.Alias != "" {
		size += len(grammar.Symbols[grammar.SYM_AS]) + len(t.Alias)
	}
	return size
}

func (t *TableIdentifier) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := copy(b, t.Name)
	if t.Alias != "" {
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_AS])
		bw += copy(b[bw:], t.Alias)
	}
	return bw
}

func (t *TableIdentifier) As(alias string) *TableIdentifier {
	tbl := &TableIdentifier{
		st:    t.st,
		Alias: alias,
		Name:  t.Name,
	}
	cols := make([]*ColumnIdentifier, len(t.columns))
	for x, c := range t.columns {
		cols[x] = &ColumnIdentifier{
			Alias: c.Alias,
			Name:  c.Name,
			tbl:   tbl,
		}
	}
	tbl.columns = cols
	return tbl
}

// TableIdentifierFromSchema returns a TableIdentifier of a given name from a
// supplied Schema
func TableIdentifierFromSchema(
	s *schema.Schema,
	name string,
) *TableIdentifier {
	st := s.T(name)
	if st == nil {
		return nil
	}
	ti := &TableIdentifier{
		st:   st,
		Name: name,
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
			Name: cname,
		}
	}

	ti.columns = cols
	return ti
}
