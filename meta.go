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
)

type Column struct {
    Name string
    Table *Table
}

type Table struct {
    Name string
    Schema string
    Columns map[string]*Column
}

type Meta struct {
    db *sql.DB
    tables map[string]*Table
    schemaName string
}

func Reflect(driver string, db *sql.DB, meta *Meta) error {
    schemaName := getSchemaName(driver, db)
    // Grab information about all tables in the schema
    rows, err := db.Query(IS_TABLES, schemaName)
    if err != nil {
        return err
    }
    tables := make(map[string]*Table, 0)
    for rows.Next() {
        table := &Table{Schema: schemaName}
        err = rows.Scan(&table.Name)
        if err != nil {
            return err
        }
        tables[table.Name] = table
    }
    meta.schemaName = schemaName
    meta.tables = tables
    meta.db = db
    return nil
}

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
