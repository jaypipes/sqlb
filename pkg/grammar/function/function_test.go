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
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/grammar/function"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/jaypipes/sqlb/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestFunctions(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := sqlb.T(sc, "users")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		c     types.Projection
		qs    map[types.Dialect]string
		qargs []interface{}
	}{
		{
			name: "MAX(column)",
			c:    function.Max(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "MAX(users.name)",
				types.DIALECT_POSTGRESQL: "MAX(users.name)",
			},
		},
		{
			name: "aliased function",
			c:    function.Max(colUserName).As("max_name"),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "MAX(users.name) AS max_name",
				types.DIALECT_POSTGRESQL: "MAX(users.name) AS max_name",
			},
		},
		{
			name: "MIN(column)",
			c:    function.Min(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "MIN(users.name)",
				types.DIALECT_POSTGRESQL: "MIN(users.name)",
			},
		},
		{
			name: "SUM(column)",
			c:    function.Sum(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "SUM(users.name)",
				types.DIALECT_POSTGRESQL: "SUM(users.name)",
			},
		},
		{
			name: "AVG(column)",
			c:    function.Avg(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "AVG(users.name)",
				types.DIALECT_POSTGRESQL: "AVG(users.name)",
			},
		},
		{
			name: "COUNT(*)",
			c:    function.Count(users),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "COUNT(*)",
				types.DIALECT_POSTGRESQL: "COUNT(*)",
			},
		},
		{
			name: "COUNT(DISTINCT column)",
			c:    function.CountDistinct(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "COUNT(DISTINCT users.name)",
				types.DIALECT_POSTGRESQL: "COUNT(DISTINCT users.name)",
			},
		},
		{
			name: "Ensure AS alias not in COUNT(DISTINCT column)",
			c:    function.CountDistinct(colUserName.As("user_name")),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "COUNT(DISTINCT users.name)",
				types.DIALECT_POSTGRESQL: "COUNT(DISTINCT users.name)",
			},
		},
		{
			name: "CAST(column AS type)",
			c:    function.Cast(colUserName, grammar.SQL_TYPE_TEXT),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "CAST(users.name AS TEXT)",
				types.DIALECT_POSTGRESQL: "CAST(users.name AS TEXT)",
			},
		},
		{
			name: "CHAR_LENGTH(column)",
			c:    function.CharLength(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "CHAR_LENGTH(users.name)",
				types.DIALECT_POSTGRESQL: "CHAR_LENGTH(users.name)",
			},
		},
		{
			name: "BIT_LENGTH(column)",
			c:    function.BitLength(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "BIT_LENGTH(users.name)",
				types.DIALECT_POSTGRESQL: "BIT_LENGTH(users.name)",
			},
		},
		{
			name: "ASCII(column)",
			c:    function.Ascii(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "ASCII(users.name)",
				types.DIALECT_POSTGRESQL: "ASCII(users.name)",
			},
		},
		{
			name: "REVERSE(column)",
			c:    function.Reverse(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "REVERSE(users.name)",
				types.DIALECT_POSTGRESQL: "REVERSE(users.name)",
			},
		},
		{
			name: "CONCAT(column, column)",
			c:    function.Concat(colUserName, colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "CONCAT(users.name, users.name)",
				types.DIALECT_POSTGRESQL: "CONCAT(users.name, users.name)",
			},
		},
		{
			name: "CONCAT_WS(string, column, column)",
			c:    function.ConcatWs("-", colUserName, colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL: "CONCAT_WS(?, users.name, users.name)",
				// Should be:
				// types.DIALECT_POSTGRESQL: "CONCAT_WS($1, users.name, users.name)",
			},
			qargs: []interface{}{"-"},
		},
		{
			name: "NOW()",
			c:    function.Now(),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "NOW()",
				types.DIALECT_POSTGRESQL: "NOW()",
			},
		},
		{
			name: "CURRENT_TIMESTAMP()",
			c:    function.CurrentTimestamp(),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "CURRENT_TIMESTAMP()",
				types.DIALECT_POSTGRESQL: "CURRENT_TIMESTAMP()",
			},
		},
		{
			name: "CURRENT_TIME()",
			c:    function.CurrentTime(),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "CURRENT_TIME()",
				types.DIALECT_POSTGRESQL: "CURRENT_TIME()",
			},
		},
		{
			name: "CURRENT_DATE()",
			c:    function.CurrentDate(),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "CURRENT_DATE()",
				types.DIALECT_POSTGRESQL: "CURRENT_DATE()",
			},
		},
		{
			name: "EXTRACT(unit FROM column)",
			c:    function.Extract(colUserName, grammar.UNIT_MINUTE_SECOND),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL: "EXTRACT(MINUTE_SECOND FROM users.name)",
				// Should be:
				// types.DIALECT_POSTGRESQL: "EXTRACT(MINUTE_SECOND FROM TIMESTAMP users.name)",
				types.DIALECT_POSTGRESQL: "EXTRACT(MINUTE_SECOND FROM users.name)",
			},
		},
	}
	for _, test := range tests {
		expArgc := len(test.qargs)
		argc := test.c.ArgCount()
		assert.Equal(expArgc, argc)

		// Test each SQL dialect output
		for dialect, qs := range test.qs {
			sc := scanner.New(dialect)
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
