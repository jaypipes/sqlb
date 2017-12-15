package sqlb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteStatement(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		s     *deleteStatement
		qs    string
		qargs []interface{}
	}{
		{
			name: "DELETE no WHERE",
			s: &deleteStatement{
				table: users,
			},
			qs: "DELETE FROM users",
		},
		{
			name: "DELETE simple WHERE",
			s: &deleteStatement{
				table: users,
				where: &whereClause{
					filters: []*Expression{
						Equal(colUserName, "foo"),
					},
				},
			},
			qs:    "DELETE FROM users WHERE users.name = ?",
			qargs: []interface{}{"foo"},
		},
	}
	for _, test := range tests {
		expArgc := len(test.qargs)
		argc := test.s.argCount()
		assert.Equal(expArgc, argc)

		expLen := len(test.qs)
		size := test.s.size(defaultScanner)
		size += interpolationLength(DIALECT_MYSQL, argc)
		assert.Equal(expLen, size)

		b := make([]byte, size)
		curArg := 0
		written := test.s.scan(defaultScanner, b, test.qargs, &curArg)

		assert.Equal(written, size)
		assert.Equal(test.qs, string(b))
	}
}
