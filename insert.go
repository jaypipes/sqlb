package sqlb

import (
    "errors"
)

var (
    ERR_INSERT_NO_VALUES = errors.New("No values supplied.")
    ERR_INSERT_UNKNOWN_COLUMN = errors.New("Received an unknown column.")
)

type InsertQuery struct {
    e error
    b []byte
    args []interface{}
    stmt *insertStatement
}

func (q *InsertQuery) IsValid() bool {
    return q.e == nil &&  q.stmt != nil
}

func (q *InsertQuery) Error() error {
    return q.e
}

func (q *InsertQuery) String() string {
    size := q.stmt.size()
    argc := q.stmt.argCount()
    if len(q.args) != argc  {
        q.args = make([]interface{}, argc)
    }
    if len(q.b) != size {
        q.b = make([]byte, size)
    }
    q.stmt.scan(q.b, q.args)
    return string(q.b)
}

func (q *InsertQuery) StringArgs() (string, []interface{}) {
    size := q.stmt.size()
    argc := q.stmt.argCount()
    if len(q.args) != argc  {
        q.args = make([]interface{}, argc)
    }
    if len(q.b) != size {
        q.b = make([]byte, size)
    }
    q.stmt.scan(q.b, q.args)
    return string(q.b), q.args
}

// Given a table and a map of column name to value for that column to insert,
// returns an InsertQuery that will produce an INSERT SQL statement 
func Insert(t *Table, values map[string]interface{}) *InsertQuery {
    if len(values) == 0 {
        return &InsertQuery{e: ERR_INSERT_NO_VALUES}
    }

    // Make sure all keys in the map point to actual columns in the target
    // table.
    cols := make([]*Column, len(values))
    vals := make([]interface{}, len(values))
    x := 0
    for k, v := range values {
        c := t.C(k)
        if c == nil {
            return &InsertQuery{e: ERR_INSERT_UNKNOWN_COLUMN}
        }
        cols[x] = c
        vals[x] = v
        x++
    }

    stmt := &insertStatement{
        table: t,
        columns: cols,
        values: vals,
    }
    return &InsertQuery{stmt: stmt}
}

func (t *Table) Insert(values map[string]interface{}) *InsertQuery {
    return Insert(t, values)
}
