//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"errors"

	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/types"
)

var (
	ERR_INSERT_NO_VALUES      = errors.New("No values supplied.")
	ERR_INSERT_UNKNOWN_COLUMN = errors.New("Received an unknown column.")
)

type InsertQuery struct {
	e       error
	b       []byte
	args    []interface{}
	stmt    *InsertStatement
	scanner types.Scanner
}

func (q *InsertQuery) IsValid() bool {
	return q.e == nil && q.stmt != nil
}

func (q *InsertQuery) Error() error {
	return q.e
}

func (q *InsertQuery) String() string {
	sizes := q.scanner.Size(q.stmt)
	if len(q.args) != sizes.ArgCount {
		q.args = make([]interface{}, sizes.ArgCount)
	}
	if len(q.b) != sizes.BufferSize {
		q.b = make([]byte, sizes.BufferSize)
	}
	q.scanner.Scan(q.b, q.args, q.stmt)
	return string(q.b)
}

func (q *InsertQuery) StringArgs() (string, []interface{}) {
	sizes := q.scanner.Size(q.stmt)
	if len(q.args) != sizes.ArgCount {
		q.args = make([]interface{}, sizes.ArgCount)
	}
	if len(q.b) != sizes.BufferSize {
		q.b = make([]byte, sizes.BufferSize)
	}
	q.scanner.Scan(q.b, q.args, q.stmt)
	return string(q.b), q.args
}

// Given a table and a map of column name to value for that column to insert,
// returns an InsertQuery that will produce an INSERT SQL statement
func Insert(t *ast.TableIdentifier, values map[string]interface{}) *InsertQuery {
	if len(values) == 0 {
		return &InsertQuery{e: ERR_INSERT_NO_VALUES}
	}

	// Make sure all keys in the map point to actual columns in the target
	// table.
	cols := make([]*ast.ColumnIdentifier, len(values))
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

	scanner := scanner.New(t.Schema().Dialect)
	stmt := &InsertStatement{
		table:   t,
		columns: cols,
		values:  vals,
	}
	return &InsertQuery{
		stmt:    stmt,
		scanner: scanner,
	}
}

//func (t *TableIdentifier) Insert(values map[string]interface{}) *InsertQuery {
//	return Insert(t, values)
//}
