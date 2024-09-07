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
	"github.com/jaypipes/sqlb/api/meta"
	"github.com/jaypipes/sqlb/internal/builder"
)

// Dialect is the SQL variant of the underlying RDBMS
type Dialect = api.Dialect

var MySQL = api.DialectMySQL
var PostgreSQL = api.DialectPostgreSQL
var PgSQL = api.DialectPostgreSQL
var PGSQL = api.DialectPostgreSQL
var TSQL = api.DialectTSQL
var MSSSQL = api.DialectTSQL
var MicrosoftSQLServer = api.DialectTSQL

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

// Equal accepts two things and returns an Element representing an equality
// expression that can be passed to a Join or Where clause.
var Equal = api.Equal

// NotEqual accepts two things and returns an Element representing an
// inequality expression that can be passed to a Join or Where clause.
var NotEqual = api.NotEqual

// And accepts two things and returns an Element representing an AND expression
// that can be passed to a Join or Where clause.
var And = api.And

// Or accepts two things and returns an Element representing an OR expression
// that can be passed to a Join or Where clause.
var Or = api.Or

// In accepts two things and returns an Element representing an IN expression
// that can be passed to a Join or Where clause.
var In = api.In

// Between accepts an element and a start and end things and returns an Element
// representing a BETWEEN expression that can be passed to a Join or Where
// clause.
var Between = api.Between

// IsNull accepts an element and returns an Element representing an IS NULL
// expression that can be passed to a Join or Where clause.
var IsNull = api.IsNull

// IsNotNull accepts an element and returns an Element representing an IS NOT
// NULL expression that can be passed to a Join or Where clause.
var IsNotNull = api.IsNotNull

// GreaterThan accepts two things and returns an Element representing a greater
// than expression that can be passed to a Join or Where clause.
var GreaterThan = api.GreaterThan

// GreaterThanOrEqual accepts two things and returns an Element representing a
// greater than or equality expression that can be passed to a Join or Where
// clause.
var GreaterThanOrEqual = api.GreaterThanOrEqual

// LessThan accepts two things and returns an Element representing a less than
// expression that can be passed to a Join or Where clause.
var LessThan = api.LessThan

// LessThanOrEqual accepts two things and returns an Element representing a
// less than or equality expression that can be passed to a Join or Where
// clause.
var LessThanOrEqual = api.LessThanOrEqual

var InvalidJoinNoSelect = api.InvalidJoinNoSelect
var InvalidJoinUnknownTarget = api.InvalidJoinUnknownTarget
var NoTargetTable = api.NoTargetTable
var NoValues = api.NoValues
var UnknownColumn = api.UnknownColumn
var TableRequired = api.TableRequired

/*
// Max returns a Projection that contains the MAX() SQL function
var Max = function.Max

// Min returns a Projection that contains the MIN() SQL function
var Min = function.Min

// Sum returns a Projection that contains the SUM() SQL function
var Sum = function.Sum

// Avg returns a Projection that contains the AVG() SQL function
var Avg = function.Avg

// Count returns a Projection that contains the COUNT() SQL function
var Count = api.Count

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
*/

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
	var el interface{}
	switch target := target.(type) {
	case *api.Selection:
		el = target.Query()
	default:
		el = target
	}

	qs, qargs := b.StringArgs(el)
	return db.QueryContext(ctx, qs, qargs...)
}

var Select = api.Select
var Update = api.Update
var Insert = api.Insert
var Delete = api.Delete
