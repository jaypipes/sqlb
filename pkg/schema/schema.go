//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package schema

import (
	"database/sql"

	"github.com/jaypipes/sqlb/pkg/types"
)

// Schema holds metadata about the tables, columns and views comprising a
// database.
type Schema struct {
	DB      *sql.DB
	Dialect types.Dialect
	Name    string
	Tables  map[string]*Table
}

// T returns a pointer to a Table with a name or alias matching the supplied
// string, or nil if no such table is known
func (s *Schema) T(name string) *Table {
	if t, ok := s.Tables[name]; ok {
		return t
	}
	for _, t := range s.Tables {
		if t.Alias == name {
			return t
		}
	}
	return nil
}

// AddTable returns a new Table that is used to describe metadata about a named
// table.
func (s *Schema) AddTable(name string) *Table {
	t := &Table{
		Schema:  s,
		Name:    name,
		Columns: map[string]*Column{},
	}
	s.Tables[name] = t
	return t
}

// New returns a new Schema that is used to describe metadata about the tables,
// columns and views comprising a database.
func New(dialect types.Dialect, name string) *Schema {
	return &Schema{
		Dialect: dialect,
		Name:    name,
		Tables:  map[string]*Table{},
	}
}
