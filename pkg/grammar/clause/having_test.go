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

func TestHaving(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := sqlb.T(sc, "users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		c     *clause.Having
		qs    string
		qargs []interface{}
	}{
		{
			name: "Empty HAVING clause",
			c:    clause.NewHaving(),
			qs:   "",
		},
		{
			name: "Single expression",
			c: clause.NewHaving(
				expression.Equal(colUserName, "foo"),
			),
			qs:    " HAVING users.name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "AND expression",
			c: clause.NewHaving(
				expression.And(
					expression.NotEqual(colUserName, "foo"),
					expression.NotEqual(colUserName, "bar"),
				),
			),
			qs:    " HAVING (users.name != ? AND users.name != ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "Multiple unary expressions should be AND'd together",
			c: clause.NewHaving(
				expression.NotEqual(colUserName, "foo"),
				expression.NotEqual(colUserName, "bar"),
			),
			qs:    " HAVING users.name != ? AND users.name != ?",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "OR expression",
			c: clause.NewHaving(
				expression.Or(
					expression.Equal(colUserName, "foo"),
					expression.Equal(colUserName, "bar"),
				),
			),
			qs:    " HAVING (users.name = ? OR users.name = ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "OR and another unary expression",
			c: clause.NewHaving(
				expression.Or(
					expression.Equal(colUserName, "foo"),
					expression.Equal(colUserName, "bar"),
				),
				expression.NotEqual(colUserName, "baz"),
			),
			qs:    " HAVING (users.name = ? OR users.name = ?) AND users.name != ?",
			qargs: []interface{}{"foo", "bar", "baz"},
		},
		{
			name: "Two AND expressions OR'd together",
			c: clause.NewHaving(
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

		var b strings.Builder
		b.Grow(size)
		curArg := 0
		test.c.Scan(scanner.DefaultScanner, &b, test.qargs, &curArg)

		assert.Equal(test.qs, b.String())
	}
}
