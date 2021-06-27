//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"testing"

	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/jaypipes/sqlb/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestHavingClause(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := T(sc, "users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		c     *HavingClause
		qs    string
		qargs []interface{}
	}{
		{
			name: "Empty HAVING clause",
			c:    &HavingClause{},
			qs:   "",
		},
		{
			name: "Single expression",
			c: &HavingClause{
				conditions: []*ast.Expression{
					ast.Equal(colUserName, "foo"),
				},
			},
			qs:    " HAVING users.name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "AND expression",
			c: &HavingClause{
				conditions: []*ast.Expression{
					ast.And(
						ast.NotEqual(colUserName, "foo"),
						ast.NotEqual(colUserName, "bar"),
					),
				},
			},
			qs:    " HAVING (users.name != ? AND users.name != ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "Multiple unary expressions should be AND'd together",
			c: &HavingClause{
				conditions: []*ast.Expression{
					ast.NotEqual(colUserName, "foo"),
					ast.NotEqual(colUserName, "bar"),
				},
			},
			qs:    " HAVING users.name != ? AND users.name != ?",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "OR expression",
			c: &HavingClause{
				conditions: []*ast.Expression{
					ast.Or(
						ast.Equal(colUserName, "foo"),
						ast.Equal(colUserName, "bar"),
					),
				},
			},
			qs:    " HAVING (users.name = ? OR users.name = ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "OR and another unary expression",
			c: &HavingClause{
				conditions: []*ast.Expression{
					ast.Or(
						ast.Equal(colUserName, "foo"),
						ast.Equal(colUserName, "bar"),
					),
					ast.NotEqual(colUserName, "baz"),
				},
			},
			qs:    " HAVING (users.name = ? OR users.name = ?) AND users.name != ?",
			qargs: []interface{}{"foo", "bar", "baz"},
		},
		{
			name: "Two AND expressions OR'd together",
			c: &HavingClause{
				conditions: []*ast.Expression{
					ast.Or(
						ast.And(
							ast.NotEqual(colUserName, "foo"),
							ast.NotEqual(colUserName, "bar"),
						),
						ast.And(
							ast.NotEqual(colUserName, "baz"),
							ast.Equal(colUserId, 1),
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
