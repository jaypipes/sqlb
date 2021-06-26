//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"testing"

	"github.com/jaypipes/sqlb/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestFormatOptions(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	articles := m.Table("articles")
	colUserName := users.C("name")
	colUserId := users.C("id")
	colArticleId := articles.C("id")
	colArticleAuthor := articles.C("author")

	stmt := &selectStatement{
		selections: []types.Selection{articles},
		projs:      []types.Projection{colArticleId, colUserName.As("author")},
		joins: []*joinClause{
			&joinClause{
				left:  articles,
				right: users,
				on:    Equal(colArticleAuthor, colUserId),
			},
		},
		where: &whereClause{
			filters: []*Expression{
				Equal(colUserName, "foo"),
			},
		},
		groupBy: &groupByClause{
			cols: []types.Projection{colUserName},
		},
		orderBy: &orderByClause{
			scols: []*sortColumn{colUserName.Desc()},
		},
		limit: &limitClause{limit: 10},
	}

	tests := []struct {
		name    string
		scanner *sqlScanner
		s       *selectStatement
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
			scanner: &sqlScanner{
				dialect: types.DIALECT_MYSQL,
				format: &types.FormatOptions{
					SeparateClauseWith: "\n",
				},
			},
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
			scanner: &sqlScanner{
				dialect: types.DIALECT_MYSQL,
				format: &types.FormatOptions{
					SeparateClauseWith: "\n",
					PrefixWith:         "\n",
				},
			},
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
		scanner := test.scanner
		if scanner == nil {
			scanner = defaultScanner
		}
		sel := &SelectQuery{
			sel:     test.s,
			scanner: scanner,
		}
		qs, qargs := sel.StringArgs()
		assert.Equal(qs, test.qs)
		assert.Equal(qargs, test.qargs)
	}
}
