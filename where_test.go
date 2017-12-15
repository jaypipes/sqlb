package sqlb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type whereClauseTest struct {
}

func TestWhereClause(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		c     *whereClause
		qs    string
		qargs []interface{}
	}{
		{
			name: "Empty WHERE clause",
			c:    &whereClause{},
			qs:   "",
		},
		{
			name: "Single expression",
			c: &whereClause{
				filters: []*Expression{
					Equal(colUserName, "foo"),
				},
			},
			qs:    " WHERE users.name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "AND expression",
			c: &whereClause{
				filters: []*Expression{
					And(
						NotEqual(colUserName, "foo"),
						NotEqual(colUserName, "bar"),
					),
				},
			},
			qs:    " WHERE (users.name != ? AND users.name != ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "Multiple unary expressions should be AND'd together",
			c: &whereClause{
				filters: []*Expression{
					NotEqual(colUserName, "foo"),
					NotEqual(colUserName, "bar"),
				},
			},
			qs:    " WHERE users.name != ? AND users.name != ?",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "OR expression",
			c: &whereClause{
				filters: []*Expression{
					Or(
						Equal(colUserName, "foo"),
						Equal(colUserName, "bar"),
					),
				},
			},
			qs:    " WHERE (users.name = ? OR users.name = ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "OR and another unary expression",
			c: &whereClause{
				filters: []*Expression{
					Or(
						Equal(colUserName, "foo"),
						Equal(colUserName, "bar"),
					),
					NotEqual(colUserName, "baz"),
				},
			},
			qs:    " WHERE (users.name = ? OR users.name = ?) AND users.name != ?",
			qargs: []interface{}{"foo", "bar", "baz"},
		},
		{
			name: "Two AND expressions OR'd together",
			c: &whereClause{
				filters: []*Expression{
					Or(
						And(
							NotEqual(colUserName, "foo"),
							NotEqual(colUserName, "bar"),
						),
						And(
							NotEqual(colUserName, "baz"),
							Equal(colUserId, 1),
						),
					),
				},
			},
			qs:    " WHERE ((users.name != ? AND users.name != ?) OR (users.name != ? AND users.id = ?))",
			qargs: []interface{}{"foo", "bar", "baz", 1},
		},
	}
	for _, test := range tests {
		expArgc := len(test.qargs)
		argc := test.c.argCount()
		assert.Equal(expArgc, argc)

		expLen := len(test.qs)
		size := test.c.size(defaultScanner)
		size += interpolationLength(DIALECT_MYSQL, argc)
		assert.Equal(expLen, size)

		b := make([]byte, size)
		curArg := 0
		written := test.c.scan(defaultScanner, b, test.qargs, &curArg)

		assert.Equal(written, size)
		assert.Equal(test.qs, string(b))
	}
}
