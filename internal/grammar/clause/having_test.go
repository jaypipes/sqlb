//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause_test

import (
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar/clause"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestHaving(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
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
		b := builder.New()

		expArgc := len(test.qargs)
		argc := test.c.ArgCount()
		assert.Equal(expArgc, argc)

		expLen := len(test.qs)
		size := test.c.Size(b)
		size += b.InterpolationLength(argc)
		assert.Equal(expLen, size)

		curArg := 0
		test.c.Scan(b, test.qargs, &curArg)

		assert.Equal(test.qs, b.String())
	}
}
