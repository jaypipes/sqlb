//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package function_test

import (
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/function"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFunctions(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		c     builder.Projection
		qs    map[api.Dialect]string
		qargs []interface{}
	}{
		{
			name: "MAX(column)",
			c:    function.Max(colUserName),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "MAX(users.name)",
				api.DialectPostgreSQL: "MAX(users.name)",
			},
		},
		{
			name: "aliased function",
			c:    function.Max(colUserName).As("max_name"),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "MAX(users.name) AS max_name",
				api.DialectPostgreSQL: "MAX(users.name) AS max_name",
			},
		},
		{
			name: "MIN(column)",
			c:    function.Min(colUserName),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "MIN(users.name)",
				api.DialectPostgreSQL: "MIN(users.name)",
			},
		},
		{
			name: "SUM(column)",
			c:    function.Sum(colUserName),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "SUM(users.name)",
				api.DialectPostgreSQL: "SUM(users.name)",
			},
		},
		{
			name: "AVG(column)",
			c:    function.Avg(colUserName),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "AVG(users.name)",
				api.DialectPostgreSQL: "AVG(users.name)",
			},
		},
		{
			name: "COUNT(*)",
			c:    function.Count(users),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "COUNT(*)",
				api.DialectPostgreSQL: "COUNT(*)",
			},
		},
		{
			name: "COUNT(DISTINCT column)",
			c:    function.CountDistinct(colUserName),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "COUNT(DISTINCT users.name)",
				api.DialectPostgreSQL: "COUNT(DISTINCT users.name)",
			},
		},
		{
			name: "Ensure AS alias not in COUNT(DISTINCT column)",
			c:    function.CountDistinct(colUserName.As("user_name")),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "COUNT(DISTINCT users.name)",
				api.DialectPostgreSQL: "COUNT(DISTINCT users.name)",
			},
		},
		{
			name: "CAST(column AS type)",
			c:    function.Cast(colUserName, grammar.SQL_TYPE_TEXT),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "CAST(users.name AS TEXT)",
				api.DialectPostgreSQL: "CAST(users.name AS TEXT)",
			},
		},
		{
			name: "CHAR_LENGTH(column)",
			c:    function.CharLength(colUserName),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "CHAR_LENGTH(users.name)",
				api.DialectPostgreSQL: "CHAR_LENGTH(users.name)",
			},
		},
		{
			name: "BIT_LENGTH(column)",
			c:    function.BitLength(colUserName),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "BIT_LENGTH(users.name)",
				api.DialectPostgreSQL: "BIT_LENGTH(users.name)",
			},
		},
		{
			name: "ASCII(column)",
			c:    function.Ascii(colUserName),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "ASCII(users.name)",
				api.DialectPostgreSQL: "ASCII(users.name)",
			},
		},
		{
			name: "REVERSE(column)",
			c:    function.Reverse(colUserName),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "REVERSE(users.name)",
				api.DialectPostgreSQL: "REVERSE(users.name)",
			},
		},
		{
			name: "CONCAT(column, column)",
			c:    function.Concat(colUserName, colUserName),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "CONCAT(users.name, users.name)",
				api.DialectPostgreSQL: "CONCAT(users.name, users.name)",
			},
		},
		{
			name: "CONCAT_WS(string, column, column)",
			c:    function.ConcatWs("-", colUserName, colUserName),
			qs: map[api.Dialect]string{
				api.DialectMySQL: "CONCAT_WS(?, users.name, users.name)",
				// Should be:
				// api.DialectPostgreSQL: "CONCAT_WS($1, users.name, users.name)",
			},
			qargs: []interface{}{"-"},
		},
		{
			name: "NOW()",
			c:    function.Now(),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "NOW()",
				api.DialectPostgreSQL: "NOW()",
			},
		},
		{
			name: "CURRENT_TIMESTAMP()",
			c:    function.CurrentTimestamp(),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "CURRENT_TIMESTAMP()",
				api.DialectPostgreSQL: "CURRENT_TIMESTAMP()",
			},
		},
		{
			name: "CURRENT_TIME()",
			c:    function.CurrentTime(),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "CURRENT_TIME()",
				api.DialectPostgreSQL: "CURRENT_TIME()",
			},
		},
		{
			name: "CURRENT_DATE()",
			c:    function.CurrentDate(),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "CURRENT_DATE()",
				api.DialectPostgreSQL: "CURRENT_DATE()",
			},
		},
		{
			name: "EXTRACT(unit FROM column)",
			c:    function.Extract(colUserName, grammar.UNIT_MINUTE_SECOND),
			qs: map[api.Dialect]string{
				api.DialectMySQL: "EXTRACT(MINUTE_SECOND FROM users.name)",
				// Should be:
				// api.DialectPostgreSQL: "EXTRACT(MINUTE_SECOND FROM TIMESTAMP users.name)",
				api.DialectPostgreSQL: "EXTRACT(MINUTE_SECOND FROM users.name)",
			},
		},
	}
	for _, test := range tests {
		expArgc := len(test.qargs)
		argc := test.c.ArgCount()
		assert.Equal(expArgc, argc)

		// Test each SQL dialect output
		for dialect, qs := range test.qs {
			b := builder.New(api.WithDialect(dialect))
			expLen := len(qs)
			size := test.c.Size(b)
			size += b.InterpolationLength(argc)
			assert.Equal(expLen, size)

			curArg := 0
			test.c.Scan(b, test.qargs, &curArg)

			assert.Equal(qs, b.String())
		}
	}
}
