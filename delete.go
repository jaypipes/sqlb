// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
package sqlb

import (
	"strings"

	"github.com/jaypipes/sqlb/errors"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
	"github.com/jaypipes/sqlb/internal/grammar/identifier"
	"github.com/jaypipes/sqlb/internal/grammar/statement"
	"github.com/jaypipes/sqlb/internal/scanner"
)

type DeleteQuery struct {
	e    error
	b    []byte
	args []interface{}
	stmt *statement.Delete
}

func (q *DeleteQuery) IsValid() bool {
	return q.e == nil && q.stmt != nil
}

func (q *DeleteQuery) Error() error {
	return q.e
}

func (q *DeleteQuery) Scan(s *scanner.Scanner, b *strings.Builder, qargs []interface{}, idx *int) {
	q.stmt.Scan(s, b, qargs, idx)
}

func (q *DeleteQuery) ArgCount() int {
	return q.stmt.ArgCount()
}

func (q *DeleteQuery) Size(s *scanner.Scanner) int {
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
		return &DeleteQuery{e: errors.NoTargetTable}
	}

	return &DeleteQuery{
		stmt: statement.NewDelete(t, nil),
	}
}

//func (t *ast.Table) Delete() *DeleteQuery {
//	return Delete(t)
//}
