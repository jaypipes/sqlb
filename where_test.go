//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"testing"

	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/jaypipes/sqlb/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestWhereClause(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := T(sc, "users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		c     *WhereClause
		qs    string
		qargs []interface{}
	}{
		{
			name: "Empty WHERE clause",
			c:    &WhereClause{},
			qs:   "",
		},
		{
			name: "Single expression",
			c: &WhereClause{
				filters: []*Expression{
					Equal(colUserName, "foo"),
				},
			},
			qs:    " WHERE users.name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "AND expression",
			c: &WhereClause{
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
			c: &WhereClause{
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
			c: &WhereClause{
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
			c: &WhereClause{
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
			c: &WhereClause{
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
		argc := test.c.ArgCount()
		assert.Equal(expArgc, argc)

		expLen := len(test.qs)
		size := test.c.Size(scanner.DefaultScanner)
		size += scanner.InterpolationLength(types.DIALECT_MYSQL, argc)
		assert.Equal(expLen, size)

		b := make([]byte, size)
		curArg := 0
		written := test.c.Scan(scanner.DefaultScanner, b, test.qargs, &curArg)

		assert.Equal(written, size)
		assert.Equal(test.qs, string(b))
	}
}
