//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package scanner_test

import (
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/jaypipes/sqlb/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestFormatOptions(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := sqlb.T(sc, "users")
	articles := sqlb.T(sc, "articles")
	colUserName := users.C("name")
	colUserId := users.C("id")
	colArticleId := articles.C("id")
	colArticleAuthor := articles.C("author")

	q := sqlb.Select(colArticleId, colUserName.As("author"))
	q.Join(articles, ast.Equal(colUserId, colArticleAuthor))
	q.Where(ast.Equal(colUserName, "foo"))
	q.GroupBy(colUserName)
	q.OrderBy(colUserName.Desc())
	q.Limit(10)

	tests := []struct {
		name    string
		scanner types.Scanner
		query   types.Element
		qs      string
		qargs   []interface{}
	}{
		{
			name:  "default space clause separator",
			query: q,
			qs:    "SELECT articles.id, users.name AS author FROM users JOIN articles ON users.id = articles.author WHERE users.name = ? GROUP BY users.name ORDER BY users.name DESC LIMIT ?",
			qargs: []interface{}{"foo", 10},
		},
		{
			name: "newline clause separator ",
			scanner: scanner.New(types.DIALECT_MYSQL).WithFormatOptions(
				&types.FormatOptions{
					SeparateClauseWith: "\n",
				},
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
			scanner: scanner.New(types.DIALECT_MYSQL).WithFormatOptions(
				&types.FormatOptions{
					SeparateClauseWith: "\n",
					PrefixWith:         "\n",
				},
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
	for _, test := range tests {
		sc := test.scanner
		if sc == nil {
			sc = scanner.DefaultScanner
		}
		qs, qargs := sc.StringArgs(test.query)
		assert.Equal(test.qs, qs)
		assert.Equal(test.qargs, qargs)
	}
}
