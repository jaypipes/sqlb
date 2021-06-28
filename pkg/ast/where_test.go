//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ast_test

import (
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/jaypipes/sqlb/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestWhereClause(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := sqlb.T(sc, "users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		c     *ast.WhereClause
		qs    string
		qargs []interface{}
	}{
		{
			name: "Empty WHERE clause",
			c:    &ast.WhereClause{},
			qs:   "",
		},
		{
			name: "Single expression",
			c: ast.NewWhereClause(
				ast.Equal(colUserName, "foo"),
			),
			qs:    " WHERE users.name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "AND expression",
			c: ast.NewWhereClause(
				ast.And(
					ast.NotEqual(colUserName, "foo"),
					ast.NotEqual(colUserName, "bar"),
				),
			),
			qs:    " WHERE (users.name != ? AND users.name != ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "Multiple unary expressions should be AND'd together",
			c: ast.NewWhereClause(
				ast.NotEqual(colUserName, "foo"),
				ast.NotEqual(colUserName, "bar"),
			),
			qs:    " WHERE users.name != ? AND users.name != ?",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "OR expression",
			c: ast.NewWhereClause(
				ast.Or(
					ast.Equal(colUserName, "foo"),
					ast.Equal(colUserName, "bar"),
				),
			),
			qs:    " WHERE (users.name = ? OR users.name = ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "OR and another unary expression",
			c: ast.NewWhereClause(
				ast.Or(
					ast.Equal(colUserName, "foo"),
					ast.Equal(colUserName, "bar"),
				),
				ast.NotEqual(colUserName, "baz"),
			),
			qs:    " WHERE (users.name = ? OR users.name = ?) AND users.name != ?",
			qargs: []interface{}{"foo", "bar", "baz"},
		},
		{
			name: "Two AND expressions OR'd together",
			c: ast.NewWhereClause(
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
			),
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
