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
		expLen := len(test.qs)
		s := test.s.size()
		assert.Equal(expLen, s)

		expArgc := len(test.qargs)
		assert.Equal(expArgc, test.s.argCount())

		b := make([]byte, s)
		written, _ := test.s.scan(b, test.qargs)

		assert.Equal(written, s)
		assert.Equal(test.qs, string(b))
	}
}
