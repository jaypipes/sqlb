// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package sqlb

import (
	"strings"

	"github.com/jaypipes/sqlb/pkg/errors"
	"github.com/jaypipes/sqlb/pkg/grammar/identifier"
	"github.com/jaypipes/sqlb/pkg/grammar/statement"
	"github.com/jaypipes/sqlb/pkg/types"
)

type InsertQuery struct {
	e    error
	stmt *statement.Insert
}

func (q *InsertQuery) IsValid() bool {
	return q.e == nil && q.stmt != nil
}

func (q *InsertQuery) Error() error {
	return q.e
}

func (q *InsertQuery) Scan(s types.Scanner, b *strings.Builder, qargs []interface{}, idx *int) {
	q.stmt.Scan(s, b, qargs, idx)
}

func (q *InsertQuery) ArgCount() int {
	return q.stmt.ArgCount()
}

func (q *InsertQuery) Size(s types.Scanner) int {
	return q.stmt.Size(s)
}

// Given a table and a map of column name to value for that column to insert,
// returns an InsertQuery that will produce an INSERT SQL statement
func Insert(t *identifier.Table, values map[string]interface{}) *InsertQuery {
	if len(values) == 0 {
		return &InsertQuery{e: errors.NoValues}
	}

	// Make sure all keys in the map point to actual columns in the target
	// table.
	cols := make([]*identifier.Column, len(values))
	vals := make([]interface{}, len(values))
	x := 0
	for k, v := range values {
		c := t.C(k)
		if c == nil {
			return &InsertQuery{e: errors.UnknownColumn}
		}
		cols[x] = c
		vals[x] = v
		x++
	}

	return &InsertQuery{
		stmt: statement.NewInsert(t, cols, vals),
	}
}

//func (t *Table) Insert(values map[string]interface{}) *InsertQuery {
//	return Insert(t, values)
//}
