package sqlb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertStatement(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	colUserName := users.C("name")
	colUserId := users.C("id")

	tests := []struct {
		name  string
		s     *insertStatement
		qs    string
		qargs []interface{}
	}{
		{
			name: "Simple INSERT",
			s: &insertStatement{
				table:   users,
				columns: []*Column{colUserId, colUserName},
				values:  []interface{}{nil, "foo"},
			},
			qs:    "INSERT INTO users (id, name) VALUES (?, ?)",
			qargs: []interface{}{nil, "foo"},
		},
		{
			name: "Ensure no aliasing in table names",
			s: &insertStatement{
				table:   users.As("u"),
				columns: []*Column{colUserId, colUserName},
				values:  []interface{}{nil, "foo"},
			},
			qs:    "INSERT INTO users (id, name) VALUES (?, ?)",
			qargs: []interface{}{nil, "foo"},
		},
		{
			name: "Ensure no aliasing in column names",
			s: &insertStatement{
				table:   users,
				columns: []*Column{colUserId.As("user_id"), colUserName},
				values:  []interface{}{nil, "foo"},
			},
			qs:    "INSERT INTO users (id, name) VALUES (?, ?)",
			qargs: []interface{}{nil, "foo"},
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
