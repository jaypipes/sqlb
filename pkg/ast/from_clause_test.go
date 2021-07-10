//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ast_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/jaypipes/sqlb/pkg/types"
)

func TestFromClause(t *testing.T) {
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
		s     *ast.FromClause
		qs    string
		qargs []interface{}
	}{
		{
			name: "Table",
			s: ast.NewFromClause(
				[]types.Selection{users},
				[]*ast.JoinClause{},
			),
			qs: "FROM users",
		},
		{
			name: "aliased Table",
			s: ast.NewFromClause(
				[]types.Selection{users.As("u")},
				[]*ast.JoinClause{},
			),
			qs: "FROM users AS u",
		},
		{
			name: "Single JOIN",
			s: ast.NewFromClause(
				[]types.Selection{articles},
				[]*ast.JoinClause{
					ast.Join(
						articles,
						users,
						ast.Equal(colArticleAuthor, colUserID),
					),
				},
			),
			qs: "FROM articles JOIN users ON articles.author = users.id",
		},
		{
			name: "Multiple JOINs",
			s: ast.NewFromClause(
				[]types.Selection{articles},
				[]*ast.JoinClause{
					ast.Join(
						articles,
						users,
						ast.Equal(colArticleAuthor, colUserID),
					),
					ast.Join(
						articles,
						articleStates,
						ast.Equal(colArticleState, colArticleStateID),
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

		b := make([]byte, size)
		curArg := 0
		written := test.s.Scan(scanner.DefaultScanner, b, test.qargs, &curArg)

		assert.Equal(written, size)
		assert.Equal(test.qs, string(b))
	}
}
