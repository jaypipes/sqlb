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

// Reflect examines the supplied database connection and discovers Table
// definitions within that connection's associated database, returning a
// pointer to a Schema struct with the discovered information.
func Reflect(
	dialect types.Dialect,
	db *sql.DB,
) (*Schema, error) {
	var err error
	dbName := DatabaseName(dialect, db)
	schema := New(dialect, dbName)
	schema.DB = db
	if err = fillTables(db, schema); err != nil {
		return nil, err
	}
	if err = fillTableColumns(db, schema); err != nil {
		return nil, err
	}
	return schema, nil
}

// fillTables populates the supplied schema's Tables collection by inspecting
// the INFORMATION_SCHEMA in the associated database.
func fillTables(
	db *sql.DB,
	schema *Schema,
) error {
	var qs string
	switch schema.Dialect {
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
	rows, err := db.Query(qs, schema.Name)
	if err != nil {
		return err
	}
	for rows.Next() {
		var tName string
		if err = rows.Scan(&tName); err != nil {
			return err
		}
		schema.AddTable(tName)
	}
	return nil
}

// fillTableColumns grabs column information from the information schema and
// populates the supplied Schema's map of Table's columns
func fillTableColumns(
	db *sql.DB,
	schema *Schema,
) error {
	var qs string
	switch schema.Dialect {
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
	rows, err := db.Query(qs, schema.Name)
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
		t = schema.T(tname)
		t.AddColumn(cname)
	}
	return nil
}

// DatabaseName returns the database schema name given a dialect and a sql.DB
// handle
func DatabaseName(dialect types.Dialect, db *sql.DB) string {
	var qs string
	switch dialect {
	case types.DIALECT_MYSQL:
		qs = "SELECT DATABASE()"
	case types.DIALECT_POSTGRESQL:
		qs = "SELECT CURRENT_DATABASE()"
	}
	var dbName string
	if err := db.QueryRow(qs).Scan(&dbName); err != nil {
		return ""
	}
	return dbName
}
