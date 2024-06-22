//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar/clause"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
	"github.com/jaypipes/sqlb/internal/testutil"
)

func TestFrom(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	articles := sqlb.T(m, "articles")
	articleStates := sqlb.T(m, "article_states")

	colUserID := users.C("id")
	colArticleAuthor := articles.C("author")
	colArticleState := articles.C("state")
	colArticleStateID := articleStates.C("id")

	tests := []struct {
		name  string
		c     *clause.From
		qs    string
		qargs []interface{}
	}{
		{
			name: "Table",
			c: clause.NewFrom(
				[]api.Selection{users},
				[]*clause.Join{},
			),
			qs: "FROM users",
		},
		{
			name: "aliased Table",
			c: clause.NewFrom(
				[]api.Selection{users.As("u")},
				[]*clause.Join{},
			),
			qs: "FROM users AS u",
		},
		{
			name: "Single JOIN",
			c: clause.NewFrom(
				[]api.Selection{articles},
				[]*clause.Join{
					clause.InnerJoin(
						articles,
						users,
						expression.Equal(colArticleAuthor, colUserID),
					),
				},
			),
			qs: "FROM articles JOIN users ON articles.author = users.id",
		},
		{
			name: "Multiple JOINs",
			c: clause.NewFrom(
				[]api.Selection{articles},
				[]*clause.Join{
					clause.InnerJoin(
						articles,
						users,
						expression.Equal(colArticleAuthor, colUserID),
					),
					clause.InnerJoin(
						articles,
						articleStates,
						expression.Equal(colArticleState, colArticleStateID),
					),
				},
			),
			qs: "FROM articles JOIN users ON articles.author = users.id JOIN article_states ON articles.state = article_states.id",
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
