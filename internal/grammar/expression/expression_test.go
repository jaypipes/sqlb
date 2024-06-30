//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package expression_test

import (
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

type expressionTest struct {
	el    api.Element
	qs    string
	qargs []interface{}
}

func TestExpressions(t *testing.T) {
	assert := assert.New(t)

	m := testutil.M()
	users := sqlb.T(m, "users")
	articles := sqlb.T(m, "articles")
	colUserId := users.C("id")
	colUserName := users.C("name")
	colArticleAuthor := articles.C("author")

	tests := []expressionTest{
		// equal value
		{
			el:    expression.Equal(colUserName, "foo"),
			qs:    "users.name = ?",
			qargs: []interface{}{"foo"},
		},
		// reverse args equal
		{
			el:    expression.Equal("foo", colUserName),
			qs:    "? = users.name",
			qargs: []interface{}{"foo"},
		},
		// equal columns
		{
			el: expression.Equal(colUserId, colArticleAuthor),
			qs: "users.id = articles.author",
		},
		// not equal value
		{
			el:    expression.NotEqual(colUserName, "foo"),
			qs:    "users.name != ?",
			qargs: []interface{}{"foo"},
		},
		// in single value
		{
			el:    expression.In(colUserName, "foo"),
			qs:    "users.name IN (?)",
			qargs: []interface{}{"foo"},
		},
		// in multi value
		{
			el:    expression.In(colUserName, "foo", "bar", 1),
			qs:    "users.name IN (?, ?, ?)",
			qargs: []interface{}{"foo", "bar", 1},
		},
		// AND expression
		{
			el: expression.And(
				expression.NotEqual(colUserName, "foo"),
				expression.NotEqual(colUserName, "bar"),
			),
			qs:    "(users.name != ? AND users.name != ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		// OR expression
		{
			el: expression.Or(
				expression.Equal(colUserName, "foo"),
				expression.Equal(colUserName, "bar"),
			),
			qs:    "(users.name = ? OR users.name = ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		// BETWEEN column and two values
		{
			el:    expression.Between(colUserName, "foo", "bar"),
			qs:    "users.name BETWEEN ? AND ?",
			qargs: []interface{}{"foo", "bar"},
		},
		// column IS NULL
		{
			el: expression.IsNull(colUserName),
			qs: "users.name IS NULL",
		},
		// column IS NOT NULL
		{
			el: expression.IsNotNull(colUserName),
			qs: "users.name IS NOT NULL",
		},
		// col > value
		{
			el:    expression.GreaterThan(colUserName, "foo"),
			qs:    "users.name > ?",
			qargs: []interface{}{"foo"},
		},
		// col >= value
		{
			el:    expression.GreaterThanOrEqual(colUserName, "foo"),
			qs:    "users.name >= ?",
			qargs: []interface{}{"foo"},
		},
		// col < value
		{
			el:    expression.LessThan(colUserName, "foo"),
			qs:    "users.name < ?",
			qargs: []interface{}{"foo"},
		},
		// col <= value
		{
			el:    expression.LessThanOrEqual(colUserName, "foo"),
			qs:    "users.name <= ?",
			qargs: []interface{}{"foo"},
		},
	}
	for _, test := range tests {
		expArgc := len(test.qargs)
		argc := test.el.ArgCount()
		assert.Equal(expArgc, argc)

		b := builder.New()

		qs, args := b.StringArgs(test.el)

		assert.Equal(test.qs, qs)
		if len(test.qargs) > 0 {
			assert.Equal(test.qargs, args)
		}
	}
}
