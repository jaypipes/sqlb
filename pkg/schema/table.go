//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package schema

// Table describes metadata about a table in a database.
type Table struct {
	Schema *Schema
	Alias  string
	Name   string
	// Columns is a map of Column structs, keyed by the column's actual name
	// (not alias)
	Columns map[string]*Column
}

// C returns a pointer to a Column with a name or alias matching the supplied
// string, or nil if no such column is known
func (t *Table) C(name string) *Column {
	if c, ok := t.Columns[name]; ok {
		return c
	}
	for _, c := range t.Columns {
		if c.Alias == name {
			return c
		}
	}
	return nil
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
