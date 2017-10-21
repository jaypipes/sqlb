package sqlb

import (
	"errors"
)

var (
	ERR_DELETE_NO_TARGET = errors.New("No target table supplied.")
)

type DeleteQuery struct {
	e    error
	b    []byte
	args []interface{}
	stmt *deleteStatement
}

func (q *DeleteQuery) IsValid() bool {
	return q.e == nil && q.stmt != nil
}

func (q *DeleteQuery) Error() error {
	return q.e
}

func (q *DeleteQuery) String() string {
	size := q.stmt.size()
	argc := q.stmt.argCount()
	if len(q.args) != argc {
		q.args = make([]interface{}, argc)
	}
	if len(q.b) != size {
		q.b = make([]byte, size)
	}
	q.stmt.scan(q.b, q.args)
	return string(q.b)
}

func (q *DeleteQuery) StringArgs() (string, []interface{}) {
	size := q.stmt.size()
	argc := q.stmt.argCount()
	if len(q.args) != argc {
		q.args = make([]interface{}, argc)
	}
	if len(q.b) != size {
		q.b = make([]byte, size)
	}
	q.stmt.scan(q.b, q.args)
	return string(q.b), q.args
}

func (q *DeleteQuery) Where(e *Expression) *DeleteQuery {
	q.stmt.addWhere(e)
	return q
}

// Given a table and a map of column name to value for that column to insert,
// returns an DeleteQuery that will produce an INSERT SQL statement
func Delete(t *Table) *DeleteQuery {
	if t == nil {
		return &DeleteQuery{e: ERR_DELETE_NO_TARGET}
	}

	stmt := &deleteStatement{
		table: t,
	}
	return &DeleteQuery{stmt: stmt}
}

func (t *Table) Delete() *DeleteQuery {
	return Delete(t)
}
