package sqlb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatOptions(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	colUserName := users.C("name")

	tests := []struct {
		name    string
		scanner *sqlScanner
		s       *selectStatement
		qs      string
		qargs   []interface{}
	}{
		{
			name: "default space clause separator",
			s: &selectStatement{
				selections: []selection{users},
				projs:      []projection{colUserName},
				where: &whereClause{
					filters: []*Expression{
						Equal(colUserName, "foo"),
					},
				},
			},
			qs:    "SELECT users.name FROM users WHERE users.name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "newline clause separator ",
			s: &selectStatement{
				selections: []selection{users},
				projs:      []projection{colUserName},
				where: &whereClause{
					filters: []*Expression{
						Equal(colUserName, "foo"),
					},
				},
			},
			scanner: &sqlScanner{
				dialect: DIALECT_MYSQL,
				format: &FormatOptions{
					SeparateClauseWith: "\n",
				},
			},
			qs:    "SELECT users.name\nFROM users\nWHERE users.name = ?",
			qargs: []interface{}{"foo"},
		},
	}
	for _, test := range tests {
		scanner := test.scanner
		if scanner == nil {
			scanner = defaultScanner
		}
		expArgc := len(test.qargs)
		argc := test.s.argCount()
		assert.Equal(expArgc, argc)

		expLen := len(test.qs)
		size := test.s.size(scanner)
		size += interpolationLength(DIALECT_MYSQL, argc)
		assert.Equal(expLen, size)

		b := make([]byte, size)
		curArg := 0
		written := test.s.scan(scanner, b, test.qargs, &curArg)

		assert.Equal(written, size)
		assert.Equal(test.qs, string(b))
	}
}
