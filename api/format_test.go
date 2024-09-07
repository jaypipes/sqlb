//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api_test

import (
	"testing"

	"github.com/jaypipes/sqlb/api"
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

	q := api.Select(colArticleId, colUserName.As("author"))
	q.Join(articles, api.Equal(colUserId, colArticleAuthor))
	q.Where(api.Equal(colUserName, "foo"))
	q.GroupBy(colUserName)
	q.OrderBy(colUserName.Desc())
	q.Limit(10)

	tests := []struct {
		name  string
		b     *builder.Builder
		query *api.Selection
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
				api.WithDialect(api.DialectMySQL),
				api.WithFormatSeparateClauseWith("\n"),
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
				api.WithDialect(api.DialectMySQL),
				api.WithFormatSeparateClauseWith("\n"),
				api.WithFormatPrefixWith("\n"),
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
