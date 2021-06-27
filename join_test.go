//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"testing"

	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/stretchr/testify/assert"
)

type JoinClauseTest struct {
	c     *JoinClause
	qs    string
	qargs []interface{}
}

func TestJoinClause(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	articles := m.Table("articles")
	colUserId := users.C("id")
	colArticleAuthor := articles.C("author")

	auCond := Equal(colArticleAuthor, colUserId)
	uaCond := Equal(colUserId, colArticleAuthor)

	tests := []JoinClauseTest{
		// articles to users table defs
		JoinClauseTest{
			c:  Join(articles, users, auCond),
			qs: " JOIN users ON articles.author = users.id",
		},
		// users to articles table defs
		JoinClauseTest{
			c:  Join(users, articles, uaCond),
			qs: " JOIN articles ON users.id = articles.author",
		},
		// articles to users tables
		JoinClauseTest{
			c:  Join(articles, users, auCond),
			qs: " JOIN users ON articles.author = users.id",
		},
		// join an aliased table to non-aliased table
		JoinClauseTest{
			c: &JoinClause{
				left:  articles.As("a"),
				right: users,
				on:    Equal(articles.As("a").C("author"), colUserId),
			},
			qs: " JOIN users ON a.author = users.id",
		},
		// join a non-aliased table to aliased table
		JoinClauseTest{
			c: &JoinClause{
				left:  articles,
				right: users.As("u"),
				on:    Equal(colArticleAuthor, users.As("u").C("id")),
			},
			qs: " JOIN users AS u ON articles.author = u.id",
		},
		// aliased projections should not include "AS alias" in output
		JoinClauseTest{
			c: &JoinClause{
				left:  articles,
				right: users,
				on:    Equal(colArticleAuthor, colUserId.As("user_id")),
			},
			qs: " JOIN users ON articles.author = users.id",
		},
		// simple outer join manual construction
		JoinClauseTest{
			c: &JoinClause{
				JoinType: JOIN_OUTER,
				left:     articles,
				right:    users,
				on:       Equal(colArticleAuthor, colUserId),
			},
			qs: " LEFT JOIN users ON articles.author = users.id",
		},
		// OuterJoin() function
		JoinClauseTest{
			c:  OuterJoin(articles, users, Equal(colArticleAuthor, colUserId)),
			qs: " LEFT JOIN users ON articles.author = users.id",
		},
		// cross join manual construction
		JoinClauseTest{
			c: &JoinClause{
				JoinType: JOIN_CROSS,
				left:     articles,
				right:    users,
			},
			qs: " CROSS JOIN users",
		},
		// CrossJoin() function
		JoinClauseTest{
			c:  CrossJoin(articles, users),
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
