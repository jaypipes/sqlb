//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import "strings"

// Table describes metadata about a table in a database.
type Table struct {
	// Meta is a pointer at the metadata collection for the database
	Meta *Meta
	// Name is the name of the table in the database
	Name string
	// Columns is a map of Column structs, keyed by the column's actual name
	// (not alias)
	Columns map[string]*Column
}

// C returns a pointer to a Column with a name matching the supplied string, or
// nil if no such column is known
//
// The name matching is done using case-insensitive matching, since this is how
// the SQL standard works for identifiers and symbols (even though Microsoft
// SQL Server uses case-sensitive identifier names).
func (t *Table) C(name string) *Column {
	if c, ok := t.Columns[name]; ok {
		return c
	}
	for _, c := range t.Columns {
		if strings.EqualFold(c.Name, name) {
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

// AddColumn returns a new Column that is used to describe metadata about a
// named column.
func (t *Table) AddColumn(name string) *Column {
	c := &Column{
		Table: t,
		Name:  name,
	}
	t.Columns[name] = c
	return c
}
