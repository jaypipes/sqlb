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
	ERR_DELETE_NO_TARGET = errors.New("No target table supplied.")
)

type DeleteQuery struct {
	e       error
	b       []byte
	args    []interface{}
	stmt    *DeleteStatement
	scanner types.Scanner
}

func (q *DeleteQuery) IsValid() bool {
	return q.e == nil && q.stmt != nil
}

func (q *DeleteQuery) Error() error {
	return q.e
}

func (q *DeleteQuery) String() string {
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

func (q *DeleteQuery) StringArgs() (string, []interface{}) {
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

func (q *DeleteQuery) Where(e *ast.Expression) *DeleteQuery {
	q.stmt.AddWhere(e)
	return q
}

// Delete returns a DeleteQuery given a table that will produce a DELETE SQL
// statement
func Delete(t *ast.TableIdentifier) *DeleteQuery {
	if t == nil {
		return &DeleteQuery{e: ERR_DELETE_NO_TARGET}
	}

	scanner := scanner.New(t.Schema().Dialect)
	stmt := &DeleteStatement{
		table: t,
	}
	return &DeleteQuery{
		stmt:    stmt,
		scanner: scanner,
	}
}

//func (t *ast.TableIdentifier) Delete() *DeleteQuery {
//	return Delete(t)
//}
