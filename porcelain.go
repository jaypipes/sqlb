//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sqlb

import (
	"context"
	"database/sql"

	"github.com/jaypipes/sqlb/errors"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar/identifier"
	"github.com/jaypipes/sqlb/internal/grammar/statement"
	"github.com/jaypipes/sqlb/meta"
	"github.com/jaypipes/sqlb/query"
)

// Reflect examines the supplied database connection and discovers Table
// definitions within that connection's associated database, returning a
// pointer to a Meta struct with the discovered information.
var Reflect = meta.Reflect

// T returns a TableIdentifier of a given name from a supplied Meta
func T(m *meta.Meta, name string) *identifier.Table {
	return identifier.TableFromMeta(m, name)
}

// Query accepts a `database/sql` `DB` handle and a `pkg/builder.Element` and
// calls the `databases/sql.DB.Query` method on the SQL string produced by the
// `Element`.
func Query(
	db *sql.DB,
	el builder.Element,
) (*sql.Rows, error) {
	return QueryContext(context.TODO(), db, el)
}

// QueryContext accepts a `database/sql` `DB` handle and a `pkg/builder.Element`
// and calls the `database/sql.DB.QueryContext` method on the SQL string
// produced by the `Element`.
func QueryContext(
	ctx context.Context,
	db *sql.DB,
	el builder.Element,
) (*sql.Rows, error) {
	s := builder.New()
	qs, qargs := s.StringArgs(el)
	return db.QueryContext(ctx, qs, qargs...)
}

// Select returns a Queryable that produces a SELECT SQL statement for one or
// more items. Items can be a Table, a Column, a Function, another SELECT
// query, or even a literal value.
//
// Select panics if sqlb cannot compile the supplied arguments into a valid
// SELECT SQL query. This is intentional, as we want compile-time failures for
// invalid SQL construction.
func Select(
	items ...interface{},
) *query.SelectQuery {
	return query.Select(items...)
}

// Insert returns a Queryable that produces an INSERT SQL statement for a given
// table and map of column name to value for that column to insert,
func Insert(
	t *identifier.Table,
	values map[string]interface{},
) (builder.Element, error) {
	if t == nil {
		return nil, errors.TableRequired
	}
	if len(values) == 0 {
		return nil, errors.NoValues
	}

	// Make sure all keys in the map point to actual columns in the target
	// table.
	cols := make([]*identifier.Column, len(values))
	vals := make([]interface{}, len(values))
	x := 0
	for k, v := range values {
		c := t.C(k)
		if c == nil {
			return nil, errors.UnknownColumn
		}
		cols[x] = c
		vals[x] = v
		x++
	}

	return statement.NewInsert(t, cols, vals), nil
}

// Delete returns a Queryable that produces a DELETE SQL statement for a given
// table
func Delete(
	t *identifier.Table,
) (builder.Element, error) {
	if t == nil {
		return nil, errors.TableRequired
	}

	return statement.NewDelete(t, nil), nil
}

// Update returns a Queryable that produces an UPDATE SQL statement for a given
// table and map of column name to value for that column to update.
func Update(
	t *identifier.Table,
	values map[string]interface{},
) (builder.Element, error) {
	if t == nil {
		return nil, errors.TableRequired
	}
	if len(values) == 0 {
		return nil, errors.NoValues
	}

	// Make sure all keys in the map point to actual columns in the target
	// table.
	cols := make([]*identifier.Column, len(values))
	vals := make([]interface{}, len(values))
	x := 0
	for k, v := range values {
		c := t.C(k)
		if c == nil {
			return nil, errors.UnknownColumn
		}
		cols[x] = c
		vals[x] = v
		x++
	}

	return statement.NewUpdate(t, cols, vals, nil), nil
}
