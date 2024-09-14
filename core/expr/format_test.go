//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package expr_test

import (
	"testing"

	"github.com/jaypipes/sqlb/core/expr"
	"github.com/jaypipes/sqlb/core/types"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	articles := m.T("articles")
	colUserName := users.C("name")
	colUserId := users.C("id")
	colArticleId := articles.C("id")
	colArticleAuthor := articles.C("author")

	q := expr.Select(colArticleId, colUserName.As("author"))
	q.Join(articles, expr.Equal(colUserId, colArticleAuthor))
	q.Where(expr.Equal(colUserName, "foo"))
	q.GroupBy(colUserName)
	q.OrderBy(colUserName.Desc())
	q.Limit(10)

	tests := []struct {
		name  string
		b     *builder.Builder
		query *expr.Selection
		qs    string
		qargs []interface{}
	}{
		{
			name:  "default space clause separator",
			b:     builder.New(),
			query: q,
			qs:    "SELECT articles.id, users.name AS author FROM users JOIN articles ON users.id = articles.author WHERE users.name = ? GROUP BY users.name ORDER BY users.name DESC LIMIT ?",
			qargs: []interface{}{"foo", 10},
		},
		{
			name: "newline clause separator",
			b: builder.New(
				types.WithDialect(types.DialectMySQL),
				types.WithFormatSeparateClauseWith("\n"),
			),
			query: q,
			qs: `SELECT articles.id, users.name AS author
FROM users
JOIN articles ON users.id = articles.author
WHERE users.name = ?
GROUP BY users.name
ORDER BY users.name DESC
LIMIT ?`,
			qargs: []interface{}{"foo", 10},
		},
		{
			name: "newline clause separator with prefix newline",
			b: builder.New(
				types.WithDialect(types.DialectMySQL),
				types.WithFormatSeparateClauseWith("\n"),
				types.WithFormatPrefixWith("\n"),
			),
			query: q,
			qs: `
SELECT articles.id, users.name AS author
FROM users
JOIN articles ON users.id = articles.author
WHERE users.name = ?
GROUP BY users.name
ORDER BY users.name DESC
LIMIT ?`,
			qargs: []interface{}{"foo", 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			qs, qargs := tt.b.StringArgs(tt.query.Query())
			assert.Equal(tt.qs, qs)
			assert.Equal(tt.qargs, qargs)
		})
	}
}
