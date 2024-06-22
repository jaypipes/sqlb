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

func TestWhere(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
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
		b := builder.New()

		expArgc := len(test.qargs)
		argc := test.c.ArgCount()
		assert.Equal(expArgc, argc)

		qs, _ := b.StringArgs(test.c)

		assert.Equal(test.qs, qs)
	}
}
