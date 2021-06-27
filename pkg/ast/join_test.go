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
	"github.com/stretchr/testify/assert"
)

type JoinClauseTest struct {
	c     *ast.JoinClause
	qs    string
	qargs []interface{}
}

func TestJoinClause(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := sqlb.T(sc, "users")
	articles := sqlb.T(sc, "articles")
	colUserId := users.C("id")
	colArticleAuthor := articles.C("author")

	auCond := ast.Equal(colArticleAuthor, colUserId)
	uaCond := ast.Equal(colUserId, colArticleAuthor)

	tests := []JoinClauseTest{
		// articles to users table defs
		JoinClauseTest{
			c:  ast.Join(articles, users, auCond),
			qs: " JOIN users ON articles.author = users.id",
		},
		// users to articles table defs
		JoinClauseTest{
			c:  ast.Join(users, articles, uaCond),
			qs: " JOIN articles ON users.id = articles.author",
		},
		// articles to users tables
		JoinClauseTest{
			c:  ast.Join(articles, users, auCond),
			qs: " JOIN users ON articles.author = users.id",
		},
		// join an aliased table to non-aliased table
		JoinClauseTest{
			c: ast.Join(
				articles.As("a"),
				users,
				ast.Equal(articles.As("a").C("author"), colUserId),
			),
			qs: " JOIN users ON a.author = users.id",
		},
		// join a non-aliased table to aliased table
		JoinClauseTest{
			c: ast.Join(
				articles,
				users.As("u"),
				ast.Equal(colArticleAuthor, users.As("u").C("id")),
			),
			qs: " JOIN users AS u ON articles.author = u.id",
		},
		// aliased projections should not include "AS alias" in output
		JoinClauseTest{
			c: ast.Join(
				articles,
				users,
				ast.Equal(colArticleAuthor, colUserId.As("user_id")),
			),
			qs: " JOIN users ON articles.author = users.id",
		},
		// OuterJoin() function
		JoinClauseTest{
			c: ast.OuterJoin(
				articles,
				users,
				ast.Equal(colArticleAuthor, colUserId),
			),
			qs: " LEFT JOIN users ON articles.author = users.id",
		},
		// CrossJoin() function
		JoinClauseTest{
			c:  ast.CrossJoin(articles, users),
			qs: " CROSS JOIN users",
		},
	}
	for _, test := range tests {
		expLen := len(test.qs)
		s := test.c.Size(scanner.DefaultScanner)
		assert.Equal(expLen, s)

		expArgc := len(test.qargs)
		assert.Equal(expArgc, test.c.ArgCount())

		b := make([]byte, s)
		curArg := 0
		written := test.c.Scan(scanner.DefaultScanner, b, test.qargs, &curArg)

		assert.Equal(written, s)
		assert.Equal(test.qs, string(b))
	}
}
