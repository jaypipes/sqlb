package sqlb

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"
)

const (
    MYSQL_GET_DB = `SELECT DATABASE()`
    PGSQL_GET_DB = `SELECT CURRENT_DATABASE()`
    IS_TABLES = `
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

type Meta struct {
    db *sql.DB
    schemaName string
    tdefs map[string]*TableDef
}

func NewMeta(driver string, schemaName string) *Meta {
    return &Meta{
        schemaName: schemaName,
        tdefs: make(map[string]*TableDef, 0),
    }
}

// Create and return a new TableDef with the given table name.
func (m *Meta) NewTableDef(tblName string) *TableDef {
    td, exists := m.tdefs[tblName]
    if exists {
        return td
    }
    td = &TableDef{meta: m, name: tblName}
    m.tdefs[tblName] = td
    return td
}

func (m *Meta) TableDef(tblName string) *TableDef {
    return m.tdefs[tblName]
}

func (m *Meta) Table(tblName string) *Table {
    td, found := m.tdefs[tblName]
    if ! found {
        return nil
    }
    return td.Table()
}

func Reflect(driver string, db *sql.DB, meta *Meta) error {
    schemaName := getSchemaName(driver, db)
    // Grab information about all tables in the schema
    rows, err := db.Query(IS_TABLES, schemaName)
    if err != nil {
        return err
    }
    meta.schemaName = schemaName
    tdefs := make(map[string]*TableDef, 0)
    for rows.Next() {
        tdef := &TableDef{meta: meta}
        err = rows.Scan(&tdef.name)
        if err != nil {
            return err
        }
        tdefs[tdef.name] = tdef
    }
    if err = fillTableColumnDefs(db, schemaName, &tdefs); err != nil {
        return err
    }
    meta.tdefs = tdefs
    meta.db = db
    return nil
}

// Grabs column information from the information schema and populates the
// supplied map of TableDef descriptors' columns
func fillTableColumnDefs(db *sql.DB, schemaName string, tdefs *map[string]*TableDef) error {
    rows, err := db.Query(IS_COLUMNS, schemaName)
    if err != nil {
        return err
    }
    var tdef *TableDef
    for rows.Next() {
        var tname string
        var cname string
        err = rows.Scan(&tname, &cname)
        if err != nil {
            return err
        }
        tdef = (*tdefs)[tname]
        if tdef.cdefs == nil {
            tdef.cdefs = make([]*ColumnDef, 0)
        }
        cdef := &ColumnDef{tdef: tdef, name: cname}
        tdef.cdefs = append(tdef.cdefs, cdef)
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
