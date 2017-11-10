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
				DIALECT_MYSQL: "MAX(users.name)",
			},
		},
		{
			name: "aliased function",
			c:    Max(colUserName).As("max_name"),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "MAX(users.name) AS max_name",
			},
		},
		{
			name: "MIN(column)",
			c:    Min(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "MIN(users.name)",
			},
		},
		{
			name: "SUM(column)",
			c:    Sum(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "SUM(users.name)",
			},
		},
		{
			name: "AVG(column)",
			c:    Avg(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "AVG(users.name)",
			},
		},
		{
			name: "COUNT(*)",
			c:    Count(users),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "COUNT(*)",
			},
		},
		{
			name: "COUNT(DISTINCT column)",
			c:    CountDistinct(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "COUNT(DISTINCT users.name)",
			},
		},
		{
			name: "Ensure AS alias not in COUNT(DISTINCT column)",
			c:    CountDistinct(colUserName.As("user_name")),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "COUNT(DISTINCT users.name)",
			},
		},
		{
			name: "CAST(column AS type)",
			c:    Cast(colUserName, SQL_TYPE_TEXT),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "CAST(users.name AS TEXT)",
			},
		},
		{
			name: "TRIM(column)",
			c:    Trim(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "TRIM(users.name)",
			},
		},
		{
			name: "CHAR_LENGTH(column)",
			c:    CharLength(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "CHAR_LENGTH(users.name)",
			},
		},
		{
			name: "BIT_LENGTH(column)",
			c:    BitLength(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "BIT_LENGTH(users.name)",
			},
		},
		{
			name: "ASCII(column)",
			c:    Ascii(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "ASCII(users.name)",
			},
		},
		{
			name: "REVERSE(column)",
			c:    Reverse(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "REVERSE(users.name)",
			},
		},
		{
			name: "CONCAT(column, column)",
			c:    Concat(colUserName, colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "CONCAT(users.name, users.name)",
			},
		},
		{
			name: "CONCAT_WS(string, column, column)",
			c:    ConcatWs("-", colUserName, colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "CONCAT_WS(?, users.name, users.name)",
			},
			qargs: []interface{}{"-"},
		},
		{
			name: "NOW()",
			c:    Now(),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "NOW()",
			},
		},
		{
			name: "CURRENT_TIMESTAMP()",
			c:    CurrentTimestamp(),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "CURRENT_TIMESTAMP()",
			},
		},
		{
			name: "CURRENT_TIME()",
			c:    CurrentTime(),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "CURRENT_TIME()",
			},
		},
		{
			name: "CURRENT_DATE()",
			c:    CurrentDate(),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "CURRENT_DATE()",
			},
		},
		{
			name: "EXTRACT(unit FROM column)",
			c:    Extract(colUserName, UNIT_MINUTE_SECOND),
			qs: map[Dialect]string{
				DIALECT_MYSQL: "EXTRACT(MINUTE_SECOND FROM users.name)",
			},
		},
	}
	for _, test := range tests {
		expArgc := len(test.qargs)
		argc := test.c.argCount()
		assert.Equal(expArgc, argc)

		// Test the MySQL SQL dialect output
		{
			qs := test.qs[DIALECT_MYSQL]
			expLen := len(qs)
			size := test.c.size()
			size += interpolationLength(DIALECT_MYSQL, argc)
			assert.Equal(expLen, size)

			b := make([]byte, size)
			curArg := 0
			written := test.c.scan(b, test.qargs, &curArg)

			assert.Equal(written, size)
			assert.Equal(qs, string(b))
		}
	}
}
