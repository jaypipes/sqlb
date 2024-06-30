//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sqlb

import (
	"context"
	"database/sql"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
	"github.com/jaypipes/sqlb/internal/grammar/function"
	"github.com/jaypipes/sqlb/internal/grammar/identifier"
	"github.com/jaypipes/sqlb/internal/grammar/statement"
	"github.com/jaypipes/sqlb/internal/query"
	"github.com/jaypipes/sqlb/meta"
)

// Dialect is the SQL variant of the underlying RDBMS
type Dialect = api.Dialect

// WithDialect informs sqlb of the Dialect
var WithDialect = api.WithDialect

// WithFormatSeparateClauseWith instructs sqlb to use a supplied string as the
// separator between clauses
var WithFormatSeparateClauseWith = api.WithFormatSeparateClauseWith

// WithFormatPrefixWith instructs sqlb to use a supplied string as a prefix for
// the resulting SQL string
var WithFormatPrefixWith = api.WithFormatPrefixWith

// Reflect examines the supplied database connection and discovers Table
// definitions within that connection's associated database, returning a
// pointer to a Meta struct with the discovered information.
var Reflect = meta.Reflect

// Meta holds metadata about the tables, columns and views comprising a
// database.
type Meta api.Meta

// Table describes metadata about a table in a database.
type Table api.Table

// Column describes a column in a Table
type Column api.Column

// Equal accepts two things and returns an Element representing an equality
// expression that can be passed to a Join or Where clause.
var Equal = expression.Equal

// NotEqual accepts two things and returns an Element representing an
// inequality expression that can be passed to a Join or Where clause.
var NotEqual = expression.NotEqual

// And accepts two things and returns an Element representing an AND expression
// that can be passed to a Join or Where clause.
var And = expression.And

// Or accepts two things and returns an Element representing an OR expression
// that can be passed to a Join or Where clause.
var Or = expression.Or

// In accepts two things and returns an Element representing an IN expression
// that can be passed to a Join or Where clause.
var In = expression.In

// Between accepts an element and a start and end things and returns an Element
// representing a BETWEEN expression that can be passed to a Join or Where
// clause.
var Between = expression.Between

// IsNull accepts an element and returns an Element representing an IS NULL
// expression that can be passed to a Join or Where clause.
var IsNull = expression.IsNull

// IsNotNull accepts an element and returns an Element representing an IS NOT
// NULL expression that can be passed to a Join or Where clause.
var IsNotNull = expression.IsNotNull

// GreaterThan accepts two things and returns an Element representing a greater
// than expression that can be passed to a Join or Where clause.
var GreaterThan = expression.GreaterThan

// GreaterThanOrEqual accepts two things and returns an Element representing a
// greater than or equality expression that can be passed to a Join or Where
// clause.
var GreaterThanOrEqual = expression.GreaterThanOrEqual

// LessThan accepts two things and returns an Element representing a less than
// expression that can be passed to a Join or Where clause.
var LessThan = expression.LessThan

// LessThanOrEqual accepts two things and returns an Element representing a
// less than or equality expression that can be passed to a Join or Where
// clause.
var LessThanOrEqual = expression.LessThanOrEqual

// Max returns a Projection that contains the MAX() SQL function
var Max = function.Max

// Min returns a Projection that contains the MIN() SQL function
var Min = function.Min

// Sum returns a Projection that contains the SUM() SQL function
var Sum = function.Sum

// Avg returns a Projection that contains the AVG() SQL function
var Avg = function.Avg

// Count returns a Projection that contains the COUNT() SQL function
var Count = function.Count

// CountDistint returns a Projection that contains the COUNT(x DISTINCT) SQL
// function
var CountDistinct = function.CountDistinct

// Cast returns a Projection that contains the CAST() SQL function
var Cast = function.Cast

// CharLength returns a Projection that contains the CHAR_LENGTH() SQL function
var CharLength = function.CharLength

// BitLength returns a Projection that contains the BIT_LENGTH() SQL function
var BitLength = function.BitLength

// Ascii returns a Projection that contains the ASCII() SQL function
var Ascii = function.Ascii

// Reverse returns a Projection that contains the REVERSE() SQL function
var Reverse = function.Reverse

// Concat returns a Projection that contains the CONCAT() SQL function
var Concat = function.Concat

// ConcatWs returns a Projection that contains the CONCAT_WS() SQL function
var ConcatWs = function.ConcatWs

// Now returns a Projection that contains the NOW() SQL function
var Now = function.Now

// CurrentTimestamp returns a Projection that contains the CURRENT_TIMESTAMP() SQL function
var CurrentTimestamp = function.CurrentTimestamp

// CurrentTime returns a Projection that contains the CURRENT_TIME() SQL function
var CurrentTime = function.CurrentTime

// CurrentDate returns a Projection that contains the CURRENT_DATE() SQL function
var CurrentDate = function.CurrentDate

// Extract returns a Projection that contains the EXTRACT() SQL function
var Extract = function.Extract

// T returns a TableIdentifier of a given name from a supplied Meta
func T(m *api.Meta, name string) *identifier.Table {
	t := m.Table(name)
	return identifier.TableFromMeta(t, name)
}

// Query accepts a `database/sql` `DB` handle and a queryable object (returned
// from Select(), Insert(), Update(), or Delete()) and calls the
// `databases/sql.DB.Query` method on the SQL string produced by that queryable
// object.
func Query(
	db *sql.DB,
	target interface{},
	opts ...api.Option,
) (*sql.Rows, error) {
	return QueryContext(context.TODO(), db, target, opts...)
}

// QueryContext accepts a `database/sql` `DB` handle and a queryable object
// (returned from Select(), Insert(), Update(), or Delete()) and calls the
// `databases/sql.DB.QueryContext` method on the SQL string produced by that
// queryable object.
func QueryContext(
	ctx context.Context,
	db *sql.DB,
	target interface{},
	opts ...api.Option,
) (*sql.Rows, error) {
	b := builder.New(opts...)
	var el api.Element
	switch target := target.(type) {
	case *query.SelectQuery:
		el = target.Element()
	case api.Element:
		el = target
	default:
		panic("expected either api.Element or *query.SelectQuery")
	}

	qs, qargs := b.StringArgs(el)
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
	t *api.Table,
	values map[string]interface{},
) (interface{}, error) {
	if len(values) == 0 {
		return nil, api.NoValues
	}
	if t == nil {
		return nil, api.TableRequired
	}
	return t.Insert(values)
}

// Delete returns a Queryable that produces a DELETE SQL statement for a given
// table
func Delete(
	t *identifier.Table,
) (api.Element, error) {
	if t == nil {
		return nil, api.TableRequired
	}

	return statement.NewDelete(t, nil), nil
}

// Update returns a Queryable that produces an UPDATE SQL statement for a given
// table and map of column name to value for that column to update.
func Update(
	t *identifier.Table,
	values map[string]interface{},
) (api.Element, error) {
	if t == nil {
		return nil, api.TableRequired
	}
	if len(values) == 0 {
		return nil, api.NoValues
	}

	// Make sure all keys in the map point to actual columns in the target
	// table.
	cols := make([]*identifier.Column, len(values))
	vals := make([]interface{}, len(values))
	x := 0
	for k, v := range values {
		c := t.C(k)
		if c == nil {
			return nil, api.UnknownColumn
		}
		cols[x] = c
		vals[x] = v
		x++
	}

	return statement.NewUpdate(t, cols, vals, nil), nil
}
