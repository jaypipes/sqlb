//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jaypipes/sqlb/pkg/types"
	_ "github.com/lib/pq"
)

var (
	ERR_NO_META_STRUCT = errors.New("Please pass a pointer to a sqlb.Meta struct")
)

type Meta struct {
	db         *sql.DB
	dialect    types.Dialect
	schemaName string
	tables     map[string]*Table
}

func NewMeta(dialect types.Dialect, schemaName string) *Meta {
	return &Meta{
		dialect:    dialect,
		schemaName: schemaName,
		tables:     make(map[string]*Table, 0),
	}
}

// Create and return a new Table with the given table name.
func (m *Meta) NewTable(name string) *Table {
	t, exists := m.tables[name]
	if exists {
		return t
	}
	t = &Table{meta: m, name: name}
	m.tables[name] = t
	return t
}

func (m *Meta) Table(name string) *Table {
	t, found := m.tables[name]
	if !found {
		return nil
	}
	return t
}

func Reflect(dialect types.Dialect, db *sql.DB, meta *Meta) error {
	if meta == nil {
		return ERR_NO_META_STRUCT
	}
	schemaName := getSchemaName(dialect, db)
	var qs string
	switch dialect {
	case types.DIALECT_MYSQL:
		qs = `
SELECT t.TABLE_NAME
FROM INFORMATION_SCHEMA.TABLES AS t
WHERE t.TABLE_TYPE = 'BASE TABLE'
AND t.TABLE_SCHEMA = ?
ORDER BY t.TABLE_NAME
`
	case types.DIALECT_POSTGRESQL:
		qs = `
SELECT t.TABLE_NAME
FROM INFORMATION_SCHEMA.TABLES AS t
WHERE t.TABLE_SCHEMA = 'public'
AND t.TABLE_CATALOG = $1
AND t.TABLE_TYPE = 'BASE TABLE'
ORDER BY t.TABLE_NAME
`
	}
	// Grab information about all tables in the schema
	rows, err := db.Query(qs, schemaName)
	if err != nil {
		return err
	}
	meta.schemaName = schemaName
	tables := make(map[string]*Table, 0)
	for rows.Next() {
		t := &Table{meta: meta}
		err = rows.Scan(&t.name)
		if err != nil {
			return err
		}
		tables[t.name] = t
	}
	if err = fillTableColumns(db, dialect, schemaName, &tables); err != nil {
		return err
	}
	meta.tables = tables
	meta.db = db
	meta.dialect = dialect
	return nil
}

// Grabs column information from the information schema and populates the
// supplied map of TableDef descriptors' columns
func fillTableColumns(db *sql.DB, dialect types.Dialect, schemaName string, tables *map[string]*Table) error {
	var qs string
	switch dialect {
	case types.DIALECT_MYSQL:
		qs = `
SELECT c.TABLE_NAME, c.COLUMN_NAME
FROM INFORMATION_SCHEMA.COLUMNS AS c
JOIN INFORMATION_SCHEMA.TABLES AS t
 ON t.TABLE_SCHEMA = c.TABLE_SCHEMA
 AND t.TABLE_NAME = c.TABLE_NAME
WHERE c.TABLE_SCHEMA = ?
AND t.TABLE_TYPE = 'BASE TABLE'
ORDER BY c.TABLE_NAME, c.COLUMN_NAME
`
	case types.DIALECT_POSTGRESQL:
		qs = `
SELECT c.TABLE_NAME, c.COLUMN_NAME
FROM INFORMATION_SCHEMA.COLUMNS AS c
JOIN INFORMATION_SCHEMA.TABLES AS t
 ON t.TABLE_SCHEMA = c.TABLE_SCHEMA
 AND t.TABLE_NAME = c.TABLE_NAME
WHERE c.TABLE_SCHEMA = 'public'
AND c.TABLE_CATALOG = $1
AND t.TABLE_TYPE = 'BASE TABLE'
ORDER BY c.TABLE_NAME, c.COLUMN_NAME
`
	}
	rows, err := db.Query(qs, schemaName)
	if err != nil {
		return err
	}
	var t *Table
	for rows.Next() {
		var tname string
		var cname string
		err = rows.Scan(&tname, &cname)
		if err != nil {
			return err
		}
		t = (*tables)[tname]
		if t.columns == nil {
			t.columns = make([]*Column, 0)
		}
		c := &Column{tbl: t, name: cname}
		t.columns = append(t.columns, c)
	}
	return nil
}

// Returns the database schema name given a driver name and a sql.DB handle
func getSchemaName(dialect types.Dialect, db *sql.DB) string {
	var qs string
	switch dialect {
	case types.DIALECT_MYSQL:
		qs = "SELECT DATABASE()"
	case types.DIALECT_POSTGRESQL:
		qs = "SELECT CURRENT_DATABASE()"
	}
	var schemaName string
	err := db.QueryRow(qs).Scan(&schemaName)
	switch {
	case err != nil:
		return ""
	default:
		return schemaName
	}
}
