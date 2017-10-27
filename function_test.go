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
		qs    string
		qargs []interface{}
	}{
		{
			name: "MAX(column)",
			c:    Max(colUserName),
			qs:   "MAX(users.name)",
		},
		{
			name: "aliased function",
			c:    Max(colUserName).As("max_name"),
			qs:   "MAX(users.name) AS max_name",
		},
		{
			name: "MIN(column)",
			c:    Min(colUserName),
			qs:   "MIN(users.name)",
		},
		{
			name: "SUM(column)",
			c:    Sum(colUserName),
			qs:   "SUM(users.name)",
		},
		{
			name: "AVG(column)",
			c:    Avg(colUserName),
			qs:   "AVG(users.name)",
		},
		{
			name: "COUNT(*)",
			c:    Count(users),
			qs:   "COUNT(*)",
		},
		{
			name: "COUNT(DISTINCT column)",
			c:    CountDistinct(colUserName),
			qs:   "COUNT(DISTINCT users.name)",
		},
		{
			name: "Ensure AS alias not in COUNT(DISTINCT column)",
			c:    CountDistinct(colUserName.As("user_name")),
			qs:   "COUNT(DISTINCT users.name)",
		},
		{
			name: "CAST(column AS type)",
			c:    Cast(colUserName, SQL_TYPE_TEXT),
			qs:   "CAST(users.name AS TEXT)",
		},
		{
			name: "TRIM(column)",
			c:    Trim(colUserName),
			qs:   "TRIM(users.name)",
		},
		{
			name: "CHAR_LENGTH(column)",
			c:    CharLength(colUserName),
			qs:   "CHAR_LENGTH(users.name)",
		},
		{
			name: "BIT_LENGTH(column)",
			c:    BitLength(colUserName),
			qs:   "BIT_LENGTH(users.name)",
		},
		{
			name: "ASCII(column)",
			c:    Ascii(colUserName),
			qs:   "ASCII(users.name)",
		},
		{
			name: "REVERSE(column)",
			c:    Reverse(colUserName),
			qs:   "REVERSE(users.name)",
		},
		{
			name: "CONCAT(column, column)",
			c:    Concat(colUserName, colUserName),
			qs:   "CONCAT(users.name, users.name)",
		},
		{
			name:  "CONCAT_WS(string, column, column)",
			c:     ConcatWs("-", colUserName, colUserName),
			qs:    "CONCAT_WS(?, users.name, users.name)",
			qargs: []interface{}{"-"},
		},
		{
			name: "NOW()",
			c:    Now(),
			qs:   "NOW()",
		},
		{
			name: "CURRENT_TIMESTAMP()",
			c:    CurrentTimestamp(),
			qs:   "CURRENT_TIMESTAMP()",
		},
		{
			name: "CURRENT_TIME()",
			c:    CurrentTime(),
			qs:   "CURRENT_TIME()",
		},
		{
			name: "CURRENT_DATE()",
			c:    CurrentDate(),
			qs:   "CURRENT_DATE()",
		},
		{
			name: "EXTRACT(unit FROM column)",
			c:    Extract(colUserName, UNIT_MINUTE_SECOND),
			qs:   "EXTRACT(MINUTE_SECOND FROM users.name)",
		},
	}
	for _, test := range tests {
		expLen := len(test.qs)
		s := test.c.size()
		assert.Equal(expLen, s)

		expArgc := len(test.qargs)
		assert.Equal(expArgc, test.c.argCount())

		b := make([]byte, s)
		written, _ := test.c.scan(b, test.qargs)

		assert.Equal(written, s)
		assert.Equal(test.qs, string(b))
	}
}
