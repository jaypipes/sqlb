//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"errors"
	"strings"

	"github.com/jaypipes/sqlb/pkg/grammar/expression"
	"github.com/jaypipes/sqlb/pkg/grammar/identifier"
	"github.com/jaypipes/sqlb/pkg/grammar/statement"
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
	stmt    *statement.Delete
	scanner types.Scanner
}

func (q *DeleteQuery) IsValid() bool {
	return q.e == nil && q.stmt != nil
}

func (q *DeleteQuery) Error() error {
	return q.e
}

func (q *DeleteQuery) Scan(s types.Scanner, b *strings.Builder, qargs []interface{}, idx *int) {
	q.stmt.Scan(s, b, qargs, idx)
}

func (q *DeleteQuery) ArgCount() int {
	return q.stmt.ArgCount()
}

func (q *DeleteQuery) Size(s types.Scanner) int {
	return q.stmt.Size(s)
}

func (q *DeleteQuery) Where(e *expression.Expression) *DeleteQuery {
	q.stmt.AddWhere(e)
	return q
}

// Delete returns a DeleteQuery given a table that will produce a DELETE SQL
// statement
func Delete(t *identifier.Table) *DeleteQuery {
	if t == nil {
		return &DeleteQuery{e: ERR_DELETE_NO_TARGET}
	}

	return &DeleteQuery{
		scanner: scanner.New(t.Schema().Dialect),
		stmt:    statement.NewDelete(t, nil),
	}
}

//func (t *ast.Table) Delete() *DeleteQuery {
//	return Delete(t)
//}
