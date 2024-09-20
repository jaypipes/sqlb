//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package meta

import (
	"database/sql"
	"strings"

	"github.com/jaypipes/sqlb/core/types"
)

// Meta holds metadata about the tables, columns and views comprising a
// database.
type Meta struct {
	// DB is a pointer to the underlying Go SQL database connection handle, or
	// nil if no underlying connection has yet been established.
	DB *sql.DB
	// Dialect indicates the flavour of SQL language that the database
	// supports.
	Dialect types.Dialect
	// Name is the actual name of the table within the database.
	Name string
	// Tables is a map, keyed by the table name, of pointers to Table structs
	// describing tables in this database.
	Tables map[string]*Table
}

// Table returns a pointer to a Table with a name matching the supplied string,
// or nil if no such table is known
//
// The name matching is done using case-insensitive matching, since this is how
// the SQL standard works for identifiers and symbols (even though Microsoft
// SQL Server uses case-sensitive identifier names).
func (m *Meta) Table(name string) *Table {
	if t, ok := m.Tables[name]; ok {
		return t
	}
	for _, t := range m.Tables {
		if strings.EqualFold(t.name, name) {
			return t
		}
	}
	return nil
}

// T returns a pointer to a Table with a name matching the supplied string, or
// nil if no such table is known
//
// The name matching is done using case-insensitive matching, since this is how
// the SQL standard works for identifiers and symbols (even though Microsoft
// SQL Server uses case-sensitive identifier names).
func (m *Meta) T(name string) *Table {
	return m.Table(name)
}
