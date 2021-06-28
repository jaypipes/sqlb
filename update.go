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
	ERR_UPDATE_NO_TARGET      = errors.New("No target table supplied.")
	ERR_UPDATE_NO_VALUES      = errors.New("No values supplied.")
	ERR_UPDATE_UNKNOWN_COLUMN = errors.New("Received an unknown column.")
)

type UpdateQuery struct {
	e       error
	b       []byte
	args    []interface{}
	stmt    *ast.UpdateStatement
	scanner types.Scanner
}

func (q *UpdateQuery) IsValid() bool {
	return q.e == nil && q.stmt != nil
}

func (q *UpdateQuery) Error() error {
	return q.e
}

func (q *UpdateQuery) String() string {
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

func (q *UpdateQuery) StringArgs() (string, []interface{}) {
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

func (q *UpdateQuery) Where(e *ast.Expression) *UpdateQuery {
	q.stmt.AddWhere(e)
	return q
}

// Given a table and a map of column name to value for that column to update,
// returns an UpdateQuery that will produce an UPDATE SQL statement
func Update(t *ast.TableIdentifier, values map[string]interface{}) *UpdateQuery {
	if t == nil {
		return &UpdateQuery{e: ERR_UPDATE_NO_TARGET}
	}
	if len(values) == 0 {
		return &UpdateQuery{e: ERR_UPDATE_NO_VALUES}
	}

	// Make sure all keys in the map point to actual columns in the target
	// table.
	cols := make([]*ast.ColumnIdentifier, len(values))
	vals := make([]interface{}, len(values))
	x := 0
	for k, v := range values {
		c := t.C(k)
		if c == nil {
			return &UpdateQuery{e: ERR_UPDATE_UNKNOWN_COLUMN}
		}
		cols[x] = c
		vals[x] = v
		x++
	}

	scanner := scanner.New(t.Schema().Dialect)
	stmt := ast.NewUpdateStatement(t, cols, vals, nil)
	return &UpdateQuery{
		stmt:    stmt,
		scanner: scanner,
	}
}

//func (t *TableIdentifier) Update(values map[string]interface{}) *UpdateQuery {
//	return Update(t, values)
//}
