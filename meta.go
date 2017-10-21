package sqlb

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	MYSQL_GET_DB = `SELECT DATABASE()`
	PGSQL_GET_DB = `SELECT CURRENT_DATABASE()`
	IS_TABLES    = `
SELECT t.TABLE_NAME
FROM INFORMATION_SCHEMA.TABLES AS t
WHERE t.TABLE_TYPE = 'BASE TABLE'
AND t.TABLE_SCHEMA = ?
ORDER BY t.TABLE_NAME
`
	IS_COLUMNS = `
SELECT c.TABLE_NAME, c.COLUMN_NAME
FROM INFORMATION_SCHEMA.COLUMNS AS c
JOIN INFORMATION_SCHEMA.TABLES AS t
 ON t.TABLE_SCHEMA = c.TABLE_SCHEMA
 AND t.TABLE_NAME = c.TABLE_NAME
WHERE c.TABLE_SCHEMA = ?
AND t.TABLE_TYPE = 'BASE TABLE'
ORDER BY c.TABLE_NAME, c.COLUMN_NAME
`
)

var (
	ERR_NO_META_STRUCT = errors.New("Please pass a pointer to a sqlb.Meta struct")
)

type Meta struct {
	db         *sql.DB
	schemaName string
	tables     map[string]*Table
}

func NewMeta(driver string, schemaName string) *Meta {
	return &Meta{
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

func Reflect(driver string, db *sql.DB, meta *Meta) error {
	if meta == nil {
		return ERR_NO_META_STRUCT
	}
	schemaName := getSchemaName(driver, db)
	// Grab information about all tables in the schema
	rows, err := db.Query(IS_TABLES, schemaName)
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
	if err = fillTableColumns(db, schemaName, &tables); err != nil {
		return err
	}
	meta.tables = tables
	meta.db = db
	return nil
}

// Grabs column information from the information schema and populates the
// supplied map of TableDef descriptors' columns
func fillTableColumns(db *sql.DB, schemaName string, tables *map[string]*Table) error {
	rows, err := db.Query(IS_COLUMNS, schemaName)
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
func getSchemaName(driver string, db *sql.DB) string {
	var qs string
	switch driver {
	case "mysql":
		qs = MYSQL_GET_DB
	case "pgsql":
		qs = PGSQL_GET_DB
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
