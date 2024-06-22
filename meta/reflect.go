//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package meta

import (
	"database/sql"

	"github.com/jaypipes/sqlb/api"
)

const (
	selTablesMySQL = `
SELECT t.TABLE_NAME
FROM INFORMATION_SCHEMA.TABLES AS t
WHERE t.TABLE_TYPE = 'BASE TABLE'
AND t.TABLE_SCHEMA = ?
ORDER BY t.TABLE_NAME
`
	selTablesPostgreSQL = `
SELECT t.TABLE_NAME
FROM INFORMATION_SCHEMA.TABLES AS t
WHERE t.TABLE_SCHEMA = 'public'
AND t.TABLE_CATALOG = $1
AND t.TABLE_TYPE = 'BASE TABLE'
ORDER BY t.TABLE_NAME
`
)

// Reflect examines the supplied database connection and discovers Table
// definitions within that connection's associated database, returning a
// pointer to a Meta struct with the discovered information.
func Reflect(
	db *sql.DB,
	mods ...api.OptionModifier,
) (*api.Meta, error) {
	var err error
	opts := api.MergeOptions(mods)
	if opts.Dialect != nil {
		if *opts.Dialect == api.DialectUnknown {
			d := Dialect(db)
			opts.Dialect = &d
		}
	}
	if opts.Dialect == nil {
		d := Dialect(db)
		opts.Dialect = &d
	}
	dbName := DatabaseName(db, api.WithDialect(*opts.Dialect))
	m := &api.Meta{
		DB:      db,
		Dialect: *opts.Dialect,
		Name:    dbName,
		Tables:  map[string]*api.Table{},
	}
	if err = fillTables(db, m); err != nil {
		return nil, err
	}
	if err = fillTableColumns(db, m); err != nil {
		return nil, err
	}
	return m, nil
}

// fillTables populates the supplied `Meta`'s Tables collection by inspecting
// the INFORMATION_SCHEMA in the associated database.
func fillTables(
	db *sql.DB,
	m *api.Meta,
) error {
	var qs string
	switch m.Dialect {
	case api.DialectMySQL:
		qs = selTablesMySQL
	case api.DialectPostgreSQL:
		qs = selTablesPostgreSQL
	}
	// Grab information about all tables in the schema
	rows, err := db.Query(qs, m.Name)
	if err != nil {
		return err
	}
	for rows.Next() {
		var tName string
		if err = rows.Scan(&tName); err != nil {
			return err
		}
		m.AddTable(tName)
	}
	return nil
}

// fillTableColumns grabs column information from the information schema and
// populates the supplied `Meta`'s map of Table's columns
func fillTableColumns(
	db *sql.DB,
	m *api.Meta,
) error {
	var qs string
	switch m.Dialect {
	case api.DialectMySQL:
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
	case api.DialectPostgreSQL:
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
	rows, err := db.Query(qs, m.Name)
	if err != nil {
		return err
	}
	var t *api.Table
	for rows.Next() {
		var tname string
		var cname string
		err = rows.Scan(&tname, &cname)
		if err != nil {
			return err
		}
		t = m.T(tname)
		t.AddColumn(cname)
	}
	return nil
}
