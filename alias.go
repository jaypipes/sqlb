//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sqlb

import (
	"context"
	"database/sql"

	"github.com/jaypipes/sqlb/core/expr"
	"github.com/jaypipes/sqlb/core/fn"
	"github.com/jaypipes/sqlb/core/reflect"
	"github.com/jaypipes/sqlb/core/types"
	"github.com/jaypipes/sqlb/internal/builder"
)

// Dialect is the SQL variant of the underlying RDBMS
type Dialect = types.Dialect

var MySQL = types.DialectMySQL
var PostgreSQL = types.DialectPostgreSQL
var PgSQL = types.DialectPostgreSQL
var PGSQL = types.DialectPostgreSQL
var TSQL = types.DialectTSQL
var MSSSQL = types.DialectTSQL
var MicrosoftSQLServer = types.DialectTSQL

// WithDialect informs sqlb of the Dialect
var WithDialect = types.WithDialect

// WithFormatSeparateClauseWith instructs sqlb to use a supplied string as the
// separator between clauses
var WithFormatSeparateClauseWith = types.WithFormatSeparateClauseWith

// WithFormatPrefixWith instructs sqlb to use a supplied string as a prefix for
// the resulting SQL string
var WithFormatPrefixWith = types.WithFormatPrefixWith

// Reflect examines the supplied database connection and discovers Table
// definitions within that connection's associated database, returning a
// pointer to a Meta struct with the discovered information.
var Reflect = reflect.Reflect

// Equal accepts two things and returns an Element representing an equality
// expression that can be passed to a Join or Where clause.
var Equal = expr.Equal

// NotEqual accepts two things and returns an Element representing an
// inequality expression that can be passed to a Join or Where clause.
var NotEqual = expr.NotEqual

// And accepts two things and returns an Element representing an AND expression
// that can be passed to a Join or Where clause.
var And = expr.And

// Or accepts two things and returns an Element representing an OR expression
// that can be passed to a Join or Where clause.
var Or = expr.Or

// In accepts two things and returns an Element representing an IN expression
// that can be passed to a Join or Where clause.
var In = expr.In

// Between accepts an element and a start and end things and returns an Element
// representing a BETWEEN expression that can be passed to a Join or Where
// clause.
var Between = expr.Between

// IsNull accepts an element and returns an Element representing an IS NULL
// expression that can be passed to a Join or Where clause.
var IsNull = expr.IsNull

// IsNotNull accepts an element and returns an Element representing an IS NOT
// NULL expression that can be passed to a Join or Where clause.
var IsNotNull = expr.IsNotNull

// GreaterThan accepts two things and returns an Element representing a greater
// than expression that can be passed to a Join or Where clause.
var GreaterThan = expr.GreaterThan

// GreaterThanOrEqual accepts two things and returns an Element representing a
// greater than or equality expression that can be passed to a Join or Where
// clause.
var GreaterThanOrEqual = expr.GreaterThanOrEqual

// LessThan accepts two things and returns an Element representing a less than
// expression that can be passed to a Join or Where clause.
var LessThan = expr.LessThan

// LessThanOrEqual accepts two things and returns an Element representing a
// less than or equality expression that can be passed to a Join or Where
// clause.
var LessThanOrEqual = expr.LessThanOrEqual

var InvalidJoinNoSelect = types.InvalidJoinNoSelect
var InvalidJoinUnknownTarget = types.InvalidJoinUnknownTarget
var NoTargetTable = types.NoTargetTable
var NoValues = types.NoValues
var UnknownColumn = types.UnknownColumn
var TableRequired = types.TableRequired

// Max returns a AggregateFunction that can be passed to a Select function to
// create a MAX(<value expression>) SQL function.  The supplied argument should
// be a ValueExpression or something that can be converted into a
// ValueExpression.
var Max = fn.Max

// Min returns a AggregateFunction that can be passed to a Select function to
// create a MIN(<value expression>) SQL function.  The supplied argument should
// be a ValueExpression or something that can be converted into a
// ValueExpression.
var Min = fn.Min

// Sum returns a AggregateFunction that can be passed to a Select function to
// create a SUM(<value expression>) SQL function.  The supplied argument should
// be a ValueExpression or something that can be converted into a
// ValueExpression.
var Sum = fn.Sum

// Avg returns a AggregateFunction that can be passed to a Select function to
// create a AVG(<value expression>) SQL function.  The supplied argument should
// be a ValueExpression or something that can be converted into a
// ValueExpression.
var Avg = fn.Avg

// Count returns a AggregateFunction that can be passed to a Select function.
// It accepts zero or one parameter. If no parameters are passed, the
// AggregateFunction returned represents a COUNT(*) SQL function. If a
// parameter is passed, it should be a ValueExpression or something that can be
// converted into a ValueExpression.
var Count = fn.Count

// Substring returns a SubstringFunction that produces a SUBSTRING() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The first argument is the subject of the SUBSTRING function and must be
// coercible to a character value expression. The second argument is the FROM
// portion of the SUBSTRING function, which is the index in the subject from
// which to return a substring. The second argument must be coercible to a
// numeric value expression.
var Substring = fn.Substring

// RegexSubstring returns a RegexSubstringFunction that produces a SUBSTRING()
// SQL function of the Regular Expression subtype that can be passed to sqlb
// constructs and functions like Select()
//
// The first argument is the subject of the SUBSTRING function and must be
// coercible to a character value expression. The second argument is the
// SIMILAR portion of the SUBSTRING function, which is the regular expression
// pattern to evaluate against the subject. The second argument must be
// coercible to a character value expression. The third argument is the ESCAPE
// portion of the SUBSTRING function, which is the characters that should be
// used as an escape sequence for the regular expression. The third argument
// must be coercible to a character value expression.
var RegexSubstring = fn.RegexSubstring

// Upper returns a FoldFunction that produces an UPPER() SQL function that can
// be passed to sqlb constructs and functions like Select()
//
// The only argument is the subject of the UPPER function and must be coercible
// to a character value expression.
var Upper = fn.Upper

// Lower returns a FoldFunction that produces a LOWER() SQL function that can
// be passed to sqlb constructs and functions like Select()
//
// The only argument is the subject of the LOWER function and must be coercible
// to a character value expression.
var Lower = fn.Lower

// Convert returns a TranscodingFunction that produces a CONVERT() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The first argument is the subject of the CONVERT function and must be
// coercible to a character value expression. The second argument is the USING
// portion of the CONVERT function.
var Convert = fn.Convert

// Translate returns a TransliterationFunction that produces a TRANSLATE() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The first argument is the subject of the TRANSLATE function and must be
// coercible to a character value expression. The second argument is the USING
// portion of the TRANSLATE function.
var Translate = fn.Translate

// Trim returns a TrimFunction that produces a TRIM() SQL function that can be
// passed to sqlb constructs and functions like Select()
//
// The first argument is the subject of the TRIM function and must be coercible
// to a character value expression. The second argument is the character(s) you
// wish to trim from the subject. The second argument must be coercible to a
// character value expression. The third argument specifies whether the
// leading, trailing or both sides of the subject should be trimmed.
var Trim = fn.Trim

// TrimSpace returns a TrimFunction that produces a TRIM(col) SQL
// function that can be passed to sqlb constructs and functions like Select()
var TrimSpace = fn.TrimSpace

// LTrimSpace returns a TrimFunction that produces a TRIM(LEADING col) SQL
// function that can be passed to sqlb constructs and functions like Select()
var LTrimSpace = fn.LTrimSpace

// RTrimSpace returns a TrimFunction that produces a TRIM(TRAILING col) SQL
// function that can be passed to sqlb constructs and functions like Select()
var RTrimSpace = fn.RTrimSpace

// LTrim returns a TrimFunction that produces a TRIM(LEADING char FROM col) SQL
// function that can be passed to sqlb constructs and functions like Select()
var LTrim = fn.LTrim
var TrimPrefix = fn.TrimPrefix

// RTrim returns a TrimFunction that produces a TRIM(TRAILING char FROM col) SQL
// function that can be passed to sqlb constructs and functions like Select()
var RTrim = fn.RTrim
var TrimSuffix = fn.RTrim

// Overlay returns a OverlayFunction that produces an OVERLAY() SQL function
// that can be passed to sqlb constructs and functions like Select()
//
// The first argument is the subject of the OVERLAY function and must be
// coercible to a character value expression. THe second argument is the
// PLACING portion of the OVERLAY function, which is the string that is
// replacing the characters in the subject string. The second argument must be
// coercible to a character value expression. The third argument is the FROM
// portion of the OVERLAY function, which is the index in the subject from
// which to insert the PLACING argument string. The third argument must be
// coercible to a numeric value expression.
var Overlay = fn.Overlay

// CurrentDate returns a CurrentDateFunction that produces a CURRENT_DATE() SQL
// function that can be passed to sqlb constructs and functions like Select()
var CurrentDate = fn.CurrentDate

// CurrentTime returns a CurrentTimeFunction that produces a CURRENT_TIME() SQL
// function that can be passed to sqlb constructs and functions like Select()
var CurrentTime = fn.CurrentTime

// CurrentTimestamp returns a CurrentTimestampFunction that produces a
// CURRENT_TIMESTAMP() SQL function that can be passed to sqlb constructs and
// functions like Select()
var CurrentTimestamp = fn.CurrentTimestamp

// LocalTime returns a LocalTimeFunction that produces a LOCALTIME() SQL
// function that can be passed to sqlb constructs and functions like Select()
var LocalTime = fn.LocalTime

// LocalTimestamp returns a LocalTimestampFunction that produces a
// LOCALTIMESTAMP() SQL function that can be passed to sqlb constructs and
// functions like Select()
var LocalTimestamp = fn.LocalTimestamp

// CharacterLength returns a LengthExpression that produces a CHAR_LENGTH() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The first argument is the subject of the CHAR_LENGTH function and must be
// coercible to a string value expression.
var CharacterLength = fn.CharacterLength
var CharLength = fn.CharacterLength

// OctetLength returns a LengthExpression that produces a OCTET_LENGTH() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The first argument is the subject of the OCTET_LENGTH function and must be
// coercible to a string value expression.
var OctetLength = fn.OctetLength

// Position returns a PositionExpression that produces a POSITION() SQL function
// that can be passed to sqlb constructs and functions like Select()
//
// The first argument is the subject of the POSITION function and must be
// coercible to either a blob value expression or a string value expression.
// The second argument is the thing to search for the presence of the subject.
// The second argument must also be coercible to either a blob value expression
// or a string value expression.
var Position = fn.Position

// Extract returns a ExtractExpression that produces a EXTRACT() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The first argument is the subject of the EXTRACT function and must be
// coercible to a datetime value expression or interval value expression. The
// second argument specifies which datetime or timezone field to extract from
// the value expression identified in the first parameter.
var Extract = fn.Extract

// NaturalLogarithm returns a NumericValueFunction that produces a LN() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The argument is the subject of the LN function and must be coercible to a
// numeric value expression.
var NaturalLogarithm = fn.NaturalLogarithm
var Ln = fn.NaturalLogarithm

// Absolute returns a NumericValueFunction that produces a ABS() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The argument is the subject of the ABS function and must be coercible to a
// numeric value expression.
var Absolute = fn.Absolute
var Abs = fn.Absolute

// Exponential returns a NumericValueFunction that produces a EXP() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The argument is the subject of the EXP function and must be coercible to a
// numeric value expression.
var Exponential = fn.Exponential
var Exp = fn.Exponential

// SquareRoot returns a NumericUnaryfunction that produces a SQRT() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The argument is the subject of the SQRT function and must be coercible to a
// numeric value expression.
var SquareRoot = fn.SquareRoot
var SqRt = SquareRoot

// Ceiling returns a NumericUnaryfunction that produces a CEIL() SQL function
// that can be passed to sqlb constructs and functions like Select()
//
// The argument is the subject of the CEIL function and must be coercible to a
// numeric value expression.
var Ceiling = fn.Ceiling
var Ceil = fn.Ceiling

// Floor returns a NumericUnaryfunction that produces a FLOOR() SQL function
// that can be passed to sqlb constructs and functions like Select()
//
// The argument is the subject of the FLOOR function and must be coercible to a
// numeric value expression.
var Floor = fn.Floor

/*
// Cast returns a Projection that contains the CAST() SQL function
var Cast = function.Cast

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
*/

// Query accepts a `database/sql` `DB` handle and a queryable object (returned
// from Select(), Insert(), Update(), or Delete()) and calls the
// `databases/sql.DB.Query` method on the SQL string produced by that queryable
// object.
func Query(
	db *sql.DB,
	target interface{},
	opts ...types.Option,
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
	opts ...types.Option,
) (*sql.Rows, error) {
	b := builder.New(opts...)
	var el interface{}
	switch target := target.(type) {
	case *expr.Selection:
		el = target.Query()
	default:
		el = target
	}

	qs, qargs := b.StringArgs(el)
	return db.QueryContext(ctx, qs, qargs...)
}

var Select = expr.Select
