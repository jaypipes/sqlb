//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"testing"

	"github.com/jaypipes/sqlb/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestHavingClause(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		c     *havingClause
		qs    string
		qargs []interface{}
	}{
		{
			name: "Empty HAVING clause",
			c:    &havingClause{},
			qs:   "",
		},
		{
			name: "Single expression",
			c: &havingClause{
				conditions: []*Expression{
					Equal(colUserName, "foo"),
				},
			},
			qs:    " HAVING users.name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "AND expression",
			c: &havingClause{
				conditions: []*Expression{
					And(
						NotEqual(colUserName, "foo"),
						NotEqual(colUserName, "bar"),
					),
				},
			},
			qs:    " HAVING (users.name != ? AND users.name != ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "Multiple unary expressions should be AND'd together",
			c: &havingClause{
				conditions: []*Expression{
					NotEqual(colUserName, "foo"),
					NotEqual(colUserName, "bar"),
				},
			},
			qs:    " HAVING users.name != ? AND users.name != ?",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "OR expression",
			c: &havingClause{
				conditions: []*Expression{
					Or(
						Equal(colUserName, "foo"),
						Equal(colUserName, "bar"),
					),
				},
			},
			qs:    " HAVING (users.name = ? OR users.name = ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "OR and another unary expression",
			c: &havingClause{
				conditions: []*Expression{
					Or(
						Equal(colUserName, "foo"),
						Equal(colUserName, "bar"),
					),
					NotEqual(colUserName, "baz"),
				},
			},
			qs:    " HAVING (users.name = ? OR users.name = ?) AND users.name != ?",
			qargs: []interface{}{"foo", "bar", "baz"},
		},
		{
			name: "Two AND expressions OR'd together",
			c: &havingClause{
				conditions: []*Expression{
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
			qs:    " HAVING ((users.name != ? AND users.name != ?) OR (users.name != ? AND users.id = ?))",
			qargs: []interface{}{"foo", "bar", "baz", 1},
		},
	}
	for _, test := range tests {
		expArgc := len(test.qargs)
		argc := test.c.ArgCount()
		assert.Equal(expArgc, argc)

		expLen := len(test.qs)
		size := test.c.Size(defaultScanner)
		size += interpolationLength(types.DIALECT_MYSQL, argc)
		assert.Equal(expLen, size)

		b := make([]byte, size)
		curArg := 0
		written := test.c.Scan(defaultScanner, b, test.qargs, &curArg)

		assert.Equal(written, size)
		assert.Equal(test.qs, string(b))
	}
}
