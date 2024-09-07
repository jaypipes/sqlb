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

// NewTable returns a new Table with the supplied properties
func NewTable(
	meta *Meta, // The metadata from the RDBMS
	name string, // The true name of the table
	cnames ...string, // column names
) *Table {
	t := &Table{
		m:    meta,
		name: name,
	}
	cols := make(map[string]*Column, len(cnames))
	for _, cname := range cnames {
		c := &Column{
			t:    t,
			name: cname,
		}
		cols[cname] = c
	}
	t.columns = cols
	return t
}

// Table describes metadata about a table in a database.
type Table struct {
	// Meta is a pointer at the metadata collection for the database
	m *Meta
	// Name is the name of the table in the database
	name string
	// Columns is a map of Column structs, keyed by the column's actual name
	// (not alias)
	columns map[string]*Column
	// Alias is any alias/correlation name given to this Table for use in a
	// SELECT statement
	alias string
}

// Meta returns the metadata associated with the underlying RDBMS
func (t *Table) Meta() *Meta {
	return t.m
}

// AddColumn adds a new Column to the Table. The supplied argument can be
// either a *Column or a string. If the argument is a string, a new Column with
// that name is created. If a same-named Column already existed for the Table,
// it is overwritten with the supplied Column.
func (t *Table) AddColumn(c interface{}) {
	if t.columns == nil {
		t.columns = map[string]*Column{}
	}
	if col, ok := c.(*Column); ok {
		t.columns[col.name] = col
	} else {
		cname := c.(string)
		t.columns[cname] = &Column{name: cname, t: t}
	}
}

// Name returns the true name of the Table, no alias
func (t *Table) Name() string {
	return t.name
}

// Alias returns the aliased name of the Table
func (t *Table) Alias() string {
	return t.alias
}

// ColumnsSorted returns a slice of the Table's Columns sorted by Column name.
func (t *Table) ColumnsSorted() []*Column {
	cols := make([]*Column, 0, len(t.columns))
	for _, c := range t.columns {
		cols = append(cols, c)
	}
	slices.SortFunc(cols, func(a, b *Column) int {
		return strings.Compare(a.Name(), b.Name())
	})
	return cols
}

// As returns a copy of the Table, aliased to the supplied name
func (t *Table) As(alias string) *Table {
	at := &Table{
		m:     t.m,
		name:  t.name,
		alias: alias,
	}
	// Build a copy of the table's columns and point those columns to the new
	// aliased table
	atCols := make(map[string]*Column, len(t.columns))
	for k, c := range t.columns {
		atc := &Column{
			name: c.name,
			t:    at,
		}
		atCols[k] = atc
	}
	at.columns = atCols
	return at
}

// C returns a pointer to a Column with a name matching the supplied string, or
// nil if no such column is known
//
// The name matching is done using case-insensitive matching, since this is how
// the SQL standard works for identifiers and symbols (even though Microsoft
// SQL Server uses case-sensitive identifier names).
func (t *Table) C(name string) *Column {
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
func (t *Table) Column(name string) *Column {
	return t.C(name)
}

// Insert returns an InstanceStatement that produces an INSERT SQL statement
// for the table and map of column name to value for that column to insert,
func (t *Table) Insert(
	values map[string]interface{},
) (*grammar.InsertStatement, error) {
	if t == nil {
		return nil, TableRequired
	}
	if len(values) == 0 {
		return nil, NoValues
	}

	// Make sure all keys in the map point to actual columns in the target
	// table.
	cols := make([]string, len(values))
	vals := make([]interface{}, len(values))
	x := 0
	for k, v := range values {
		c := t.C(k)
		if c == nil {
			return nil, UnknownColumn
		}
		cols[x] = c.name
		vals[x] = v
		x++
	}

	return &grammar.InsertStatement{
		TableName: t.name,
		Columns:   cols,
		Values:    vals,
	}, nil
}

// Delete returns a `grammar.DeleteStatement` that will produce an
// DELETE SQL statement when passed to Query or QueryContext.
func (t *Table) Delete() *DeleteStatement {
	return &DeleteStatement{
		ds: &grammar.DeleteStatementSearched{
			TableName: t.name,
		},
	}
}

// Update returns a `grammar.UpdateStatement` that will produce an
// UPDATE SQL statement when passed to Query or QueryContext.
//
// The supplied map of values is keyed by the column name the value will be
// updated to.
func (t *Table) Update(
	values map[string]interface{},
) (*UpdateStatement, error) {
	if len(values) == 0 {
		return nil, NoValues
	}
	cols := make([]string, len(values))
	vals := make([]interface{}, len(values))
	x := 0
	for k, v := range values {
		c := t.C(k)
		if c == nil {
			return nil, UnknownColumn
		}
		cols[x] = k
		vals[x] = v
		x++
	}
	return &UpdateStatement{
		us: &grammar.UpdateStatementSearched{
			TableName: t.name,
			Columns:   cols,
			Values:    vals,
		},
	}, nil
}
