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
