package sqlb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctions(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		c     *sqlFunc
		qs    map[Dialect]string
		qargs []interface{}
	}{
		{
			name: "MAX(column)",
			c:    Max(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "MAX(users.name)",
				DIALECT_POSTGRESQL: "MAX(users.name)",
			},
		},
		{
			name: "aliased function",
			c:    Max(colUserName).As("max_name"),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "MAX(users.name) AS max_name",
				DIALECT_POSTGRESQL: "MAX(users.name) AS max_name",
			},
		},
		{
			name: "MIN(column)",
			c:    Min(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "MIN(users.name)",
				DIALECT_POSTGRESQL: "MIN(users.name)",
			},
		},
		{
			name: "SUM(column)",
			c:    Sum(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "SUM(users.name)",
				DIALECT_POSTGRESQL: "SUM(users.name)",
			},
		},
		{
			name: "AVG(column)",
			c:    Avg(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "AVG(users.name)",
				DIALECT_POSTGRESQL: "AVG(users.name)",
			},
		},
		{
			name: "COUNT(*)",
			c:    Count(users),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "COUNT(*)",
				DIALECT_POSTGRESQL: "COUNT(*)",
			},
		},
		{
			name: "COUNT(DISTINCT column)",
			c:    CountDistinct(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "COUNT(DISTINCT users.name)",
				DIALECT_POSTGRESQL: "COUNT(DISTINCT users.name)",
			},
		},
		{
			name: "Ensure AS alias not in COUNT(DISTINCT column)",
			c:    CountDistinct(colUserName.As("user_name")),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "COUNT(DISTINCT users.name)",
				DIALECT_POSTGRESQL: "COUNT(DISTINCT users.name)",
			},
		},
		{
			name: "CAST(column AS type)",
			c:    Cast(colUserName, SQL_TYPE_TEXT),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "CAST(users.name AS TEXT)",
				DIALECT_POSTGRESQL: "CAST(users.name AS TEXT)",
			},
		},
		{
			name: "TRIM(column)",
			c:    Trim(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "TRIM(users.name)",
				// Should be:
				// DIALECT_POSTGRESQL: "BTRIM(users.name)",
				DIALECT_POSTGRESQL: "TRIM(users.name)",
			},
		},
		{
			name: "CHAR_LENGTH(column)",
			c:    CharLength(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "CHAR_LENGTH(users.name)",
				DIALECT_POSTGRESQL: "CHAR_LENGTH(users.name)",
			},
		},
		{
			name: "BIT_LENGTH(column)",
			c:    BitLength(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "BIT_LENGTH(users.name)",
				DIALECT_POSTGRESQL: "BIT_LENGTH(users.name)",
			},
		},
		{
			name: "ASCII(column)",
			c:    Ascii(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "ASCII(users.name)",
				DIALECT_POSTGRESQL: "ASCII(users.name)",
			},
		},
		{
			name: "REVERSE(column)",
			c:    Reverse(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "REVERSE(users.name)",
				DIALECT_POSTGRESQL: "REVERSE(users.name)",
			},
		},
		{
			name: "CONCAT(column, column)",
			c:    Concat(colUserName, colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "CONCAT(users.name, users.name)",
				DIALECT_POSTGRESQL: "CONCAT(users.name, users.name)",
			},
		},
		{
			name: "CONCAT_WS(string, column, column)",
			c:    ConcatWs("-", colUserName, colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "CONCAT_WS(?, users.name, users.name)",
				// Should be:
				// DIALECT_POSTGRESQL: "CONCAT_WS($1, users.name, users.name)",
			},
			qargs: []interface{}{"-"},
		},
		{
			name: "NOW()",
			c:    Now(),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "NOW()",
				DIALECT_POSTGRESQL: "NOW()",
			},
		},
		{
			name: "CURRENT_TIMESTAMP()",
			c:    CurrentTimestamp(),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "CURRENT_TIMESTAMP()",
				DIALECT_POSTGRESQL: "CURRENT_TIMESTAMP()",
			},
		},
		{
			name: "CURRENT_TIME()",
			c:    CurrentTime(),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "CURRENT_TIME()",
				DIALECT_POSTGRESQL: "CURRENT_TIME()",
			},
		},
		{
			name: "CURRENT_DATE()",
			c:    CurrentDate(),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "CURRENT_DATE()",
				DIALECT_POSTGRESQL: "CURRENT_DATE()",
			},
		},
		{
			name: "EXTRACT(unit FROM column)",
			c:    Extract(colUserName, UNIT_MINUTE_SECOND),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "EXTRACT(MINUTE_SECOND FROM users.name)",
				// Should be:
				// DIALECT_POSTGRESQL: "EXTRACT(MINUTE_SECOND FROM TIMESTAMP users.name)",
				DIALECT_POSTGRESQL: "EXTRACT(MINUTE_SECOND FROM users.name)",
			},
		},
	}
	for _, test := range tests {
		expArgc := len(test.qargs)
		argc := test.c.argCount()
		assert.Equal(expArgc, argc)

		// Test each SQL dialect output
		for dialect, qs := range test.qs {
			test.c.setDialect(dialect)
			expLen := len(qs)
			size := test.c.size()
			size += interpolationLength(dialect, argc)
			assert.Equal(expLen, size)

			b := make([]byte, size)
			curArg := 0
			written := test.c.scan(b, test.qargs, &curArg)

			assert.Equal(written, size)
			assert.Equal(qs, string(b))
		}
	}
}
