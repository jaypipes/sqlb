//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause_test

import (
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/pkg/grammar/clause"
	"github.com/jaypipes/sqlb/pkg/grammar/expression"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

type JoinTest struct {
	c     *clause.Join
	qs    string
	qargs []interface{}
}

func TestJoin(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := sqlb.T(sc, "users")
	articles := sqlb.T(sc, "articles")
	colUserId := users.C("id")
	colArticleAuthor := articles.C("author")

	auCond := expression.Equal(colArticleAuthor, colUserId)
	uaCond := expression.Equal(colUserId, colArticleAuthor)

	tests := []JoinTest{
		// articles to users table defs
		JoinTest{
			c:  clause.InnerJoin(articles, users, auCond),
			qs: " JOIN users ON articles.author = users.id",
		},
		// users to articles table defs
		JoinTest{
			c:  clause.InnerJoin(users, articles, uaCond),
			qs: " JOIN articles ON users.id = articles.author",
		},
		// articles to users tables
		JoinTest{
			c:  clause.InnerJoin(articles, users, auCond),
			qs: " JOIN users ON articles.author = users.id",
		},
		// join an aliased table to non-aliased table
		JoinTest{
			c: clause.InnerJoin(
				articles.As("a"),
				users,
				expression.Equal(articles.As("a").C("author"), colUserId),
			),
			qs: " JOIN users ON a.author = users.id",
		},
		// join a non-aliased table to aliased table
		JoinTest{
			c: clause.InnerJoin(
				articles,
				users.As("u"),
				expression.Equal(colArticleAuthor, users.As("u").C("id")),
			),
			qs: " JOIN users AS u ON articles.author = u.id",
		},
		// aliased projections should not include "AS alias" in output
		JoinTest{
			c: clause.InnerJoin(
				articles,
				users,
				expression.Equal(colArticleAuthor, colUserId.As("user_id")),
			),
			qs: " JOIN users ON articles.author = users.id",
		},
		// OuterJoin() function
		JoinTest{
			c: clause.OuterJoin(
				articles,
				users,
				expression.Equal(colArticleAuthor, colUserId),
			),
			qs: " LEFT JOIN users ON articles.author = users.id",
		},
		// CrossJoin() function
		JoinTest{
			c:  clause.CrossJoin(articles, users),
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
