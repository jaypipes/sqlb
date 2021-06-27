//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"testing"

	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/jaypipes/sqlb/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestFormatOptions(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := T(sc, "users")
	articles := T(sc, "articles")
	colUserName := users.C("name")
	colUserId := users.C("id")
	colArticleId := articles.C("id")
	colArticleAuthor := articles.C("author")

	stmt := &SelectStatement{
		selections: []types.Selection{articles},
		projs:      []types.Projection{colArticleId, colUserName.As("author")},
		joins: []*ast.JoinClause{
			ast.Join(
				articles,
				users,
				ast.Equal(colArticleAuthor, colUserId),
			),
		},
		where: ast.NewWhereClause(
			ast.Equal(colUserName, "foo"),
		),
		groupBy: ast.NewGroupByClause(colUserName),
		orderBy: ast.NewOrderByClause(colUserName.Desc()),
		limit:   ast.NewLimitClause(10, nil),
	}

	tests := []struct {
		name    string
		scanner types.Scanner
		s       *SelectStatement
		qs      string
		qargs   []interface{}
	}{
		{
			name:  "default space clause separator",
			s:     stmt,
			qs:    "SELECT articles.id, users.name AS author FROM articles JOIN users ON articles.author = users.id WHERE users.name = ? GROUP BY users.name ORDER BY users.name DESC LIMIT ?",
			qargs: []interface{}{"foo", 10},
		},
		{
			name: "newline clause separator ",
			scanner: scanner.New(types.DIALECT_MYSQL).WithFormatOptions(
				&types.FormatOptions{
					SeparateClauseWith: "\n",
				},
			),
			s: stmt,
			qs: `SELECT articles.id, users.name AS author
FROM articles
JOIN users ON articles.author = users.id
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
			s: stmt,
			qs: `
SELECT articles.id, users.name AS author
FROM articles
JOIN users ON articles.author = users.id
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
		sel := &SelectQuery{
			sel:     test.s,
			scanner: sc,
		}
		qs, qargs := sel.StringArgs()
		assert.Equal(qs, test.qs)
		assert.Equal(qargs, test.qargs)
	}
}
