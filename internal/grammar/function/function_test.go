//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package function_test

import (
	"strings"
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/function"
	"github.com/jaypipes/sqlb/internal/scanner"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/jaypipes/sqlb/types"
	"github.com/stretchr/testify/assert"
)

func TestFunctions(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		c     scanner.Projection
		qs    map[types.Dialect]string
		qargs []interface{}
	}{
		{
			name: "MAX(column)",
			c:    function.Max(colUserName),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "MAX(users.name)",
				types.DialectPostgreSQL: "MAX(users.name)",
			},
		},
		{
			name: "aliased function",
			c:    function.Max(colUserName).As("max_name"),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "MAX(users.name) AS max_name",
				types.DialectPostgreSQL: "MAX(users.name) AS max_name",
			},
		},
		{
			name: "MIN(column)",
			c:    function.Min(colUserName),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "MIN(users.name)",
				types.DialectPostgreSQL: "MIN(users.name)",
			},
		},
		{
			name: "SUM(column)",
			c:    function.Sum(colUserName),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "SUM(users.name)",
				types.DialectPostgreSQL: "SUM(users.name)",
			},
		},
		{
			name: "AVG(column)",
			c:    function.Avg(colUserName),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "AVG(users.name)",
				types.DialectPostgreSQL: "AVG(users.name)",
			},
		},
		{
			name: "COUNT(*)",
			c:    function.Count(users),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "COUNT(*)",
				types.DialectPostgreSQL: "COUNT(*)",
			},
		},
		{
			name: "COUNT(DISTINCT column)",
			c:    function.CountDistinct(colUserName),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "COUNT(DISTINCT users.name)",
				types.DialectPostgreSQL: "COUNT(DISTINCT users.name)",
			},
		},
		{
			name: "Ensure AS alias not in COUNT(DISTINCT column)",
			c:    function.CountDistinct(colUserName.As("user_name")),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "COUNT(DISTINCT users.name)",
				types.DialectPostgreSQL: "COUNT(DISTINCT users.name)",
			},
		},
		{
			name: "CAST(column AS type)",
			c:    function.Cast(colUserName, grammar.SQL_TYPE_TEXT),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "CAST(users.name AS TEXT)",
				types.DialectPostgreSQL: "CAST(users.name AS TEXT)",
			},
		},
		{
			name: "CHAR_LENGTH(column)",
			c:    function.CharLength(colUserName),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "CHAR_LENGTH(users.name)",
				types.DialectPostgreSQL: "CHAR_LENGTH(users.name)",
			},
		},
		{
			name: "BIT_LENGTH(column)",
			c:    function.BitLength(colUserName),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "BIT_LENGTH(users.name)",
				types.DialectPostgreSQL: "BIT_LENGTH(users.name)",
			},
		},
		{
			name: "ASCII(column)",
			c:    function.Ascii(colUserName),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "ASCII(users.name)",
				types.DialectPostgreSQL: "ASCII(users.name)",
			},
		},
		{
			name: "REVERSE(column)",
			c:    function.Reverse(colUserName),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "REVERSE(users.name)",
				types.DialectPostgreSQL: "REVERSE(users.name)",
			},
		},
		{
			name: "CONCAT(column, column)",
			c:    function.Concat(colUserName, colUserName),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "CONCAT(users.name, users.name)",
				types.DialectPostgreSQL: "CONCAT(users.name, users.name)",
			},
		},
		{
			name: "CONCAT_WS(string, column, column)",
			c:    function.ConcatWs("-", colUserName, colUserName),
			qs: map[types.Dialect]string{
				types.DialectMySQL: "CONCAT_WS(?, users.name, users.name)",
				// Should be:
				// types.DialectPostgreSQL: "CONCAT_WS($1, users.name, users.name)",
			},
			qargs: []interface{}{"-"},
		},
		{
			name: "NOW()",
			c:    function.Now(),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "NOW()",
				types.DialectPostgreSQL: "NOW()",
			},
		},
		{
			name: "CURRENT_TIMESTAMP()",
			c:    function.CurrentTimestamp(),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "CURRENT_TIMESTAMP()",
				types.DialectPostgreSQL: "CURRENT_TIMESTAMP()",
			},
		},
		{
			name: "CURRENT_TIME()",
			c:    function.CurrentTime(),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "CURRENT_TIME()",
				types.DialectPostgreSQL: "CURRENT_TIME()",
			},
		},
		{
			name: "CURRENT_DATE()",
			c:    function.CurrentDate(),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "CURRENT_DATE()",
				types.DialectPostgreSQL: "CURRENT_DATE()",
			},
		},
		{
			name: "EXTRACT(unit FROM column)",
			c:    function.Extract(colUserName, grammar.UNIT_MINUTE_SECOND),
			qs: map[types.Dialect]string{
				types.DialectMySQL: "EXTRACT(MINUTE_SECOND FROM users.name)",
				// Should be:
				// types.DialectPostgreSQL: "EXTRACT(MINUTE_SECOND FROM TIMESTAMP users.name)",
				types.DialectPostgreSQL: "EXTRACT(MINUTE_SECOND FROM users.name)",
			},
		},
	}
	for _, test := range tests {
		expArgc := len(test.qargs)
		argc := test.c.ArgCount()
		assert.Equal(expArgc, argc)

		// Test each SQL dialect output
		for dialect, qs := range test.qs {
			sc := &scanner.Scanner{
				Dialect: dialect,
			}
			expLen := len(qs)
			size := test.c.Size(sc)
			size += scanner.InterpolationLength(dialect, argc)
			assert.Equal(expLen, size)

			var b strings.Builder
			b.Grow(size)
			curArg := 0
			test.c.Scan(sc, &b, test.qargs, &curArg)

			assert.Equal(qs, b.String())
		}
	}
}
