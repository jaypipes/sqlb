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

type expressionTest struct {
	c     *ast.Expression
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
			c:     ast.Equal(colUserName, "foo"),
			qs:    "users.name = ?",
			qargs: []interface{}{"foo"},
		},
		// reverse args equal
		expressionTest{
			c:     ast.Equal("foo", colUserName),
			qs:    "? = users.name",
			qargs: []interface{}{"foo"},
		},
		// equal columns
		expressionTest{
			c:  ast.Equal(colUserId, colArticleAuthor),
			qs: "users.id = articles.author",
		},
		// not equal value
		expressionTest{
			c:     ast.NotEqual(colUserName, "foo"),
			qs:    "users.name != ?",
			qargs: []interface{}{"foo"},
		},
		// in single value
		expressionTest{
			c:     ast.In(colUserName, "foo"),
			qs:    "users.name IN (?)",
			qargs: []interface{}{"foo"},
		},
		// in multi value
		expressionTest{
			c:     ast.In(colUserName, "foo", "bar", 1),
			qs:    "users.name IN (?, ?, ?)",
			qargs: []interface{}{"foo", "bar", 1},
		},
		// AND expression
		expressionTest{
			c: ast.And(
				ast.NotEqual(colUserName, "foo"),
				ast.NotEqual(colUserName, "bar"),
			),
			qs:    "(users.name != ? AND users.name != ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		// OR expression
		expressionTest{
			c: ast.Or(
				ast.Equal(colUserName, "foo"),
				ast.Equal(colUserName, "bar"),
			),
			qs:    "(users.name = ? OR users.name = ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		// BETWEEN column and two values
		expressionTest{
			c:     ast.Between(colUserName, "foo", "bar"),
			qs:    "users.name BETWEEN ? AND ?",
			qargs: []interface{}{"foo", "bar"},
		},
		// column IS NULL
		expressionTest{
			c:  ast.IsNull(colUserName),
			qs: "users.name IS NULL",
		},
		// column IS NOT NULL
		expressionTest{
			c:  ast.IsNotNull(colUserName),
			qs: "users.name IS NOT NULL",
		},
		// col > value
		expressionTest{
			c:     ast.GreaterThan(colUserName, "foo"),
			qs:    "users.name > ?",
			qargs: []interface{}{"foo"},
		},
		// col >= value
		expressionTest{
			c:     ast.GreaterThanOrEqual(colUserName, "foo"),
			qs:    "users.name >= ?",
			qargs: []interface{}{"foo"},
		},
		// col < value
		expressionTest{
			c:     ast.LessThan(colUserName, "foo"),
			qs:    "users.name < ?",
			qargs: []interface{}{"foo"},
		},
		// col <= value
		expressionTest{
			c:     ast.LessThanOrEqual(colUserName, "foo"),
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

		b := make([]byte, size)
		curArg := 0
		written := test.c.Scan(scanner.DefaultScanner, b, test.qargs, &curArg)

		assert.Equal(written, size)
		assert.Equal(test.qs, string(b))
	}
}
