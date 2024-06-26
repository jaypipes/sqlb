//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import (
	"database/sql"
	"strings"
)

// Meta holds metadata about the tables, columns and views comprising a
// database.
type Meta struct {
	DB      *sql.DB
	Dialect Dialect
	Name    string
	Tables  map[string]*Table
}

// T returns a pointer to a Table with a name matching the supplied string, or
// nil if no such table is known
//
// The name matching is done using case-insensitive matching, since this is how
// the SQL standard works for identifiers and symbols (even though Microsoft
// SQL Server uses case-sensitive identifier names).
func (s *Meta) T(name string) *Table {
	if t, ok := s.Tables[name]; ok {
		return t
	}
	for _, t := range s.Tables {
		if strings.EqualFold(t.Name, name) {
			return t
		}
	}
	return nil
}

// Table returns a pointer to a Table with a name or alias matching the
// supplied string, or nil if no such table is known
func (m *Meta) Table(name string) *Table {
	return m.T(name)
}

// AddTable returns a new Table that is used to describe metadata about a named
// table.
func (s *Meta) AddTable(name string) *Table {
	t := &Table{
		Meta:    s,
		Name:    name,
		Columns: map[string]*Column{},
	}
	if s.Tables == nil {
		s.Tables = map[string]*Table{}
	}
	s.Tables[name] = t
	return t
}
