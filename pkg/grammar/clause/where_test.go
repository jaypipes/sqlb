//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause_test

import (
	"strings"
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/pkg/grammar/clause"
	"github.com/jaypipes/sqlb/pkg/grammar/expression"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/jaypipes/sqlb/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestWhere(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := sqlb.T(sc, "users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		c     *clause.Where
		qs    string
		qargs []interface{}
	}{
		{
			name: "Empty WHERE clause",
			c:    clause.NewWhere(),
			qs:   "",
		},
		{
			name: "Single expression",
			c: clause.NewWhere(
				expression.Equal(colUserName, "foo"),
			),
			qs:    " WHERE users.name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "AND expression",
			c: clause.NewWhere(
				expression.And(
					expression.NotEqual(colUserName, "foo"),
					expression.NotEqual(colUserName, "bar"),
				),
			),
			qs:    " WHERE (users.name != ? AND users.name != ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "Multiple unary expressions should be AND'd together",
			c: clause.NewWhere(
				expression.NotEqual(colUserName, "foo"),
				expression.NotEqual(colUserName, "bar"),
			),
			qs:    " WHERE users.name != ? AND users.name != ?",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "OR expression",
			c: clause.NewWhere(
				expression.Or(
					expression.Equal(colUserName, "foo"),
					expression.Equal(colUserName, "bar"),
				),
			),
			qs:    " WHERE (users.name = ? OR users.name = ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "OR and another unary expression",
			c: clause.NewWhere(
				expression.Or(
					expression.Equal(colUserName, "foo"),
					expression.Equal(colUserName, "bar"),
				),
				expression.NotEqual(colUserName, "baz"),
			),
			qs:    " WHERE (users.name = ? OR users.name = ?) AND users.name != ?",
			qargs: []interface{}{"foo", "bar", "baz"},
		},
		{
			name: "Two AND expressions OR'd together",
			c: clause.NewWhere(
				expression.Or(
					expression.And(
						expression.NotEqual(colUserName, "foo"),
						expression.NotEqual(colUserName, "bar"),
					),
					expression.And(
						expression.NotEqual(colUserName, "baz"),
						expression.Equal(colUserId, 1),
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

		var b strings.Builder
		b.Grow(size)
		curArg := 0
		test.c.Scan(scanner.DefaultScanner, &b, test.qargs, &curArg)

		assert.Equal(test.qs, b.String())
	}
}
