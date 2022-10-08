//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/pkg/grammar/clause"
	"github.com/jaypipes/sqlb/pkg/grammar/expression"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/jaypipes/sqlb/pkg/types"
)

func TestFrom(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := sqlb.T(sc, "users")
	articles := sqlb.T(sc, "articles")
	articleStates := sqlb.T(sc, "article_states")

	colUserID := users.C("id")
	colArticleAuthor := articles.C("author")
	colArticleState := articles.C("state")
	colArticleStateID := articleStates.C("id")

	tests := []struct {
		name  string
		s     *clause.From
		qs    string
		qargs []interface{}
	}{
		{
			name: "Table",
			s: clause.NewFrom(
				[]types.Selection{users},
				[]*clause.Join{},
			),
			qs: "FROM users",
		},
		{
			name: "aliased Table",
			s: clause.NewFrom(
				[]types.Selection{users.As("u")},
				[]*clause.Join{},
			),
			qs: "FROM users AS u",
		},
		{
			name: "Single JOIN",
			s: clause.NewFrom(
				[]types.Selection{articles},
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
			s: clause.NewFrom(
				[]types.Selection{articles},
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
		expArgc := len(test.qargs)
		argc := test.s.ArgCount()
		assert.Equal(expArgc, argc)

		expLen := len(test.qs)
		size := test.s.Size(scanner.DefaultScanner)
		size += scanner.InterpolationLength(types.DIALECT_MYSQL, argc)
		assert.Equal(expLen, size)

		var b strings.Builder
		b.Grow(size)
		curArg := 0
		test.s.Scan(scanner.DefaultScanner, &b, test.qargs, &curArg)

		assert.Equal(test.qs, b.String())
	}
}
