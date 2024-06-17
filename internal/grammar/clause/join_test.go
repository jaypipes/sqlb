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

type JoinTest struct {
	c     *clause.Join
	qs    string
	qargs []interface{}
}

func TestJoin(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	articles := sqlb.T(m, "articles")
	colUserId := users.C("id")
	colArticleAuthor := articles.C("author")

	auCond := expression.Equal(colArticleAuthor, colUserId)
	uaCond := expression.Equal(colUserId, colArticleAuthor)

	tests := []JoinTest{
		// articles to users table defs
		{
			c:  clause.InnerJoin(articles, users, auCond),
			qs: " JOIN users ON articles.author = users.id",
		},
		// users to articles table defs
		{
			c:  clause.InnerJoin(users, articles, uaCond),
			qs: " JOIN articles ON users.id = articles.author",
		},
		// articles to users tables
		{
			c:  clause.InnerJoin(articles, users, auCond),
			qs: " JOIN users ON articles.author = users.id",
		},
		// join an aliased table to non-aliased table
		{
			c: clause.InnerJoin(
				articles.As("a"),
				users,
				expression.Equal(articles.As("a").C("author"), colUserId),
			),
			qs: " JOIN users ON a.author = users.id",
		},
		// join a non-aliased table to aliased table
		{
			c: clause.InnerJoin(
				articles,
				users.As("u"),
				expression.Equal(colArticleAuthor, users.As("u").C("id")),
			),
			qs: " JOIN users AS u ON articles.author = u.id",
		},
		// aliased projections should not include "AS alias" in output
		{
			c: clause.InnerJoin(
				articles,
				users,
				expression.Equal(colArticleAuthor, colUserId.As("user_id")),
			),
			qs: " JOIN users ON articles.author = users.id",
		},
		// OuterJoin() function
		{
			c: clause.OuterJoin(
				articles,
				users,
				expression.Equal(colArticleAuthor, colUserId),
			),
			qs: " LEFT JOIN users ON articles.author = users.id",
		},
		// CrossJoin() function
		{
			c:  clause.CrossJoin(articles, users),
			qs: " CROSS JOIN users",
		},
	}
	for _, test := range tests {
		b := builder.New()

		expLen := len(test.qs)
		s := test.c.Size(b)
		assert.Equal(expLen, s)

		expArgc := len(test.qargs)
		assert.Equal(expArgc, test.c.ArgCount())

		curArg := 0
		test.c.Scan(b, test.qargs, &curArg)

		assert.Equal(test.qs, b.String())
	}
}
