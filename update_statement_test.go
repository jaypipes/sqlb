package sqlb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateStatement(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		s     *updateStatement
		qs    string
		qargs []interface{}
	}{
		{
			name: "UPDATE no WHERE",
			s: &updateStatement{
				table:   users,
				columns: []*Column{colUserName},
				values:  []interface{}{"foo"},
			},
			qs:    "UPDATE users SET name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "UPDATE simple WHERE",
			s: &updateStatement{
				table:   users,
				columns: []*Column{colUserName},
				values:  []interface{}{"foo"},
				where: &whereClause{
					filters: []*Expression{
						Equal(colUserName, "bar"),
					},
				},
			},
			qs:    "UPDATE users SET name = ? WHERE users.name = ?",
			qargs: []interface{}{"foo", "bar"},
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
