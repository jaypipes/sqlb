//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package expression_test

import (
	"strings"
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/pkg/grammar/expression"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/jaypipes/sqlb/pkg/types"
	"github.com/stretchr/testify/assert"
)

type expressionTest struct {
	c     *expression.Expression
	qs    string
	qargs []interface{}
}

func TestExpressions(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := sqlb.T(sc, "users")
	articles := sqlb.T(sc, "articles")
	colUserId := users.C("id")
	colUserName := users.C("name")
	colArticleAuthor := articles.C("author")

	tests := []expressionTest{
		// equal value
		expressionTest{
			c:     expression.Equal(colUserName, "foo"),
			qs:    "users.name = ?",
			qargs: []interface{}{"foo"},
		},
		// reverse args equal
		expressionTest{
			c:     expression.Equal("foo", colUserName),
			qs:    "? = users.name",
			qargs: []interface{}{"foo"},
		},
		// equal columns
		expressionTest{
			c:  expression.Equal(colUserId, colArticleAuthor),
			qs: "users.id = articles.author",
		},
		// not equal value
		expressionTest{
			c:     expression.NotEqual(colUserName, "foo"),
			qs:    "users.name != ?",
			qargs: []interface{}{"foo"},
		},
		// in single value
		expressionTest{
			c:     expression.In(colUserName, "foo"),
			qs:    "users.name IN (?)",
			qargs: []interface{}{"foo"},
		},
		// in multi value
		expressionTest{
			c:     expression.In(colUserName, "foo", "bar", 1),
			qs:    "users.name IN (?, ?, ?)",
			qargs: []interface{}{"foo", "bar", 1},
		},
		// AND expression
		expressionTest{
			c: expression.And(
				expression.NotEqual(colUserName, "foo"),
				expression.NotEqual(colUserName, "bar"),
			),
			qs:    "(users.name != ? AND users.name != ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		// OR expression
		expressionTest{
			c: expression.Or(
				expression.Equal(colUserName, "foo"),
				expression.Equal(colUserName, "bar"),
			),
			qs:    "(users.name = ? OR users.name = ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		// BETWEEN column and two values
		expressionTest{
			c:     expression.Between(colUserName, "foo", "bar"),
			qs:    "users.name BETWEEN ? AND ?",
			qargs: []interface{}{"foo", "bar"},
		},
		// column IS NULL
		expressionTest{
			c:  expression.IsNull(colUserName),
			qs: "users.name IS NULL",
		},
		// column IS NOT NULL
		expressionTest{
			c:  expression.IsNotNull(colUserName),
			qs: "users.name IS NOT NULL",
		},
		// col > value
		expressionTest{
			c:     expression.GreaterThan(colUserName, "foo"),
			qs:    "users.name > ?",
			qargs: []interface{}{"foo"},
		},
		// col >= value
		expressionTest{
			c:     expression.GreaterThanOrEqual(colUserName, "foo"),
			qs:    "users.name >= ?",
			qargs: []interface{}{"foo"},
		},
		// col < value
		expressionTest{
			c:     expression.LessThan(colUserName, "foo"),
			qs:    "users.name < ?",
			qargs: []interface{}{"foo"},
		},
		// col <= value
		expressionTest{
			c:     expression.LessThanOrEqual(colUserName, "foo"),
			qs:    "users.name <= ?",
			qargs: []interface{}{"foo"},
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
