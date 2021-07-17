//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ast_test

import (
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/grammar"
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
			c:    ast.Max(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "MAX(users.name)",
				types.DIALECT_POSTGRESQL: "MAX(users.name)",
			},
		},
		{
			name: "aliased function",
			c:    ast.Max(colUserName).As("max_name"),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "MAX(users.name) AS max_name",
				types.DIALECT_POSTGRESQL: "MAX(users.name) AS max_name",
			},
		},
		{
			name: "MIN(column)",
			c:    ast.Min(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "MIN(users.name)",
				types.DIALECT_POSTGRESQL: "MIN(users.name)",
			},
		},
		{
			name: "SUM(column)",
			c:    ast.Sum(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "SUM(users.name)",
				types.DIALECT_POSTGRESQL: "SUM(users.name)",
			},
		},
		{
			name: "AVG(column)",
			c:    ast.Avg(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "AVG(users.name)",
				types.DIALECT_POSTGRESQL: "AVG(users.name)",
			},
		},
		{
			name: "COUNT(*)",
			c:    ast.Count(users),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "COUNT(*)",
				types.DIALECT_POSTGRESQL: "COUNT(*)",
			},
		},
		{
			name: "COUNT(DISTINCT column)",
			c:    ast.CountDistinct(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "COUNT(DISTINCT users.name)",
				types.DIALECT_POSTGRESQL: "COUNT(DISTINCT users.name)",
			},
		},
		{
			name: "Ensure AS alias not in COUNT(DISTINCT column)",
			c:    ast.CountDistinct(colUserName.As("user_name")),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "COUNT(DISTINCT users.name)",
				types.DIALECT_POSTGRESQL: "COUNT(DISTINCT users.name)",
			},
		},
		{
			name: "CAST(column AS type)",
			c:    ast.Cast(colUserName, grammar.SQL_TYPE_TEXT),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "CAST(users.name AS TEXT)",
				types.DIALECT_POSTGRESQL: "CAST(users.name AS TEXT)",
			},
		},
		{
			name: "CHAR_LENGTH(column)",
			c:    ast.CharLength(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "CHAR_LENGTH(users.name)",
				types.DIALECT_POSTGRESQL: "CHAR_LENGTH(users.name)",
			},
		},
		{
			name: "BIT_LENGTH(column)",
			c:    ast.BitLength(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "BIT_LENGTH(users.name)",
				types.DIALECT_POSTGRESQL: "BIT_LENGTH(users.name)",
			},
		},
		{
			name: "ASCII(column)",
			c:    ast.Ascii(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "ASCII(users.name)",
				types.DIALECT_POSTGRESQL: "ASCII(users.name)",
			},
		},
		{
			name: "REVERSE(column)",
			c:    ast.Reverse(colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "REVERSE(users.name)",
				types.DIALECT_POSTGRESQL: "REVERSE(users.name)",
			},
		},
		{
			name: "CONCAT(column, column)",
			c:    ast.Concat(colUserName, colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "CONCAT(users.name, users.name)",
				types.DIALECT_POSTGRESQL: "CONCAT(users.name, users.name)",
			},
		},
		{
			name: "CONCAT_WS(string, column, column)",
			c:    ast.ConcatWs("-", colUserName, colUserName),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL: "CONCAT_WS(?, users.name, users.name)",
				// Should be:
				// types.DIALECT_POSTGRESQL: "CONCAT_WS($1, users.name, users.name)",
			},
			qargs: []interface{}{"-"},
		},
		{
			name: "NOW()",
			c:    ast.Now(),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "NOW()",
				types.DIALECT_POSTGRESQL: "NOW()",
			},
		},
		{
			name: "CURRENT_TIMESTAMP()",
			c:    ast.CurrentTimestamp(),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "CURRENT_TIMESTAMP()",
				types.DIALECT_POSTGRESQL: "CURRENT_TIMESTAMP()",
			},
		},
		{
			name: "CURRENT_TIME()",
			c:    ast.CurrentTime(),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "CURRENT_TIME()",
				types.DIALECT_POSTGRESQL: "CURRENT_TIME()",
			},
		},
		{
			name: "CURRENT_DATE()",
			c:    ast.CurrentDate(),
			qs: map[types.Dialect]string{
				types.DIALECT_MYSQL:      "CURRENT_DATE()",
				types.DIALECT_POSTGRESQL: "CURRENT_DATE()",
			},
		},
		{
			name: "EXTRACT(unit FROM column)",
			c:    ast.Extract(colUserName, grammar.UNIT_MINUTE_SECOND),
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

			b := make([]byte, size)
			curArg := 0
			written := test.c.Scan(sc, b, test.qargs, &curArg)

			assert.Equal(written, size)
			assert.Equal(qs, string(b))
		}
	}
}
