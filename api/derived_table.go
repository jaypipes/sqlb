//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import (
	"slices"
	"strings"

	"github.com/jaypipes/sqlb/grammar"
)

// NewDerivedTable returns a new DerivedTable from the supplied
// Selection
func NewDerivedTable(
	name string,
	sel *Selection,
) *DerivedTable {
	if sel.qs == nil {
		panic("cannot create a DerivedTable from an empty Selection")
	}
	cols := map[string]*Column{}
	dt := &DerivedTable{
		name: name,
	}
	// We need to project all columns from the supplied Selection's
	// QuerySpecification to the outer QuerySpecification.
	for _, c := range sel.ColumnsSorted() {
		outerCol := &Column{
			t:    dt,
			name: c.name,
		}
		cols[c.name] = outerCol
	}
	dt.columns = cols
	dt.qs = sel.qs
	return dt
}

// DerivedTable describes a subquery in the FROM clause of a SELECT statement.
type DerivedTable struct {
	// name is the name of the subquery in the FROM clause
	name string
	// columns is a map of Column structs, keyed by the column's actual name
	// (not alias)
	columns map[string]*Column
	// qs is the QuerySpecification the derived table encapsulates
	qs *grammar.QuerySpecification
}

// Query returns the derived table's underlying QuerySpecification which can be
// passed to a Builder
func (t *DerivedTable) Query() *grammar.QuerySpecification {
	return t.qs
}

// Name returns the name of the DerivedTable
func (t *DerivedTable) Name() string {
	return t.name
}

// ColumnsSorted returns a slice of the Table's Columns sorted by Column name.
func (t *DerivedTable) ColumnsSorted() []*Column {
	cols := make([]*Column, 0, len(t.columns))
	for _, c := range t.columns {
		cols = append(cols, c)
	}
	slices.SortFunc(cols, func(a, b *Column) int {
		return strings.Compare(a.Name(), b.Name())
	})
	return cols
}

// C returns a pointer to a Column with a name matching the supplied string, or
// nil if no such column is known
//
// The name matching is done using case-insensitive matching, since this is how
// the SQL standard works for identifiers and symbols (even though Microsoft
// SQL Server uses case-sensitive identifier names).
func (t *DerivedTable) C(name string) *Column {
	if c, ok := t.columns[name]; ok {
		return c
	}
	for _, c := range t.columns {
		if strings.EqualFold(c.name, name) {
			return c
		}
	}
	return nil
}

// Column returns a pointer to a Column with a name or alias matching the
// supplied string, or nil if no such column is known
func (t *DerivedTable) Column(name string) *Column {
	return t.C(name)
}
