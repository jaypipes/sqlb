//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jaypipes/sqlb/pkg/types"
)

func TestSelectClause(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	articles := m.Table("articles")
	article_states := m.Table("article_states")
	colUserName := users.C("name")
	colUserId := users.C("id")
	colArticleId := articles.C("id")
	colArticleAuthor := articles.C("author")
	colArticleState := articles.C("state")
	colArticleStateId := article_states.C("id")
	colArticleStateName := article_states.C("name")

	tests := []struct {
		name  string
		s     *SelectStatement
		qs    string
		qargs []interface{}
	}{
		{
			name: "A literal value",
			s: &SelectStatement{
				projs: []types.Projection{&value{val: 1}},
			},
			qs:    "SELECT ?",
			qargs: []interface{}{1},
		},
		{
			name: "A literal value aliased",
			s: &SelectStatement{
				projs: []types.Projection{
					&value{alias: "foo", val: 1},
				},
			},
			qs:    "SELECT ? AS foo",
			qargs: []interface{}{1},
		},
		{
			name: "Two literal values",
			s: &SelectStatement{
				projs: []types.Projection{
					&value{val: 1},
					&value{val: 1},
				},
			},
			qs:    "SELECT ?, ?",
			qargs: []interface{}{1, 2},
		},
		{
			name: "Table and column",
			s: &SelectStatement{
				selections: []types.Selection{users},
				projs:      []types.Projection{colUserName},
			},
			qs: "SELECT users.name FROM users",
		},
		{
			name: "aliased Table and Column",
			s: &SelectStatement{
				selections: []types.Selection{users.As("u")},
				projs: []types.Projection{
					users.As("u").C("name"),
				},
			},
			qs: "SELECT u.name FROM users AS u",
		},
		{
			name: "Table and multiple Column",
			s: &SelectStatement{
				selections: []types.Selection{users},
				projs:      []types.Projection{colUserId, colUserName},
			},
			qs: "SELECT users.id, users.name FROM users",
		},
		{
			name: "Simple WHERE",
			s: &SelectStatement{
				selections: []types.Selection{users},
				projs:      []types.Projection{colUserName},
				where: &WhereClause{
					filters: []*Expression{
						Equal(colUserName, "foo"),
					},
				},
			},
			qs:    "SELECT users.name FROM users WHERE users.name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "Simple LIMIT",
			s: &SelectStatement{
				selections: []types.Selection{users},
				projs:      []types.Projection{colUserName},
				limit:      &LimitClause{limit: 10},
			},
			qs:    "SELECT users.name FROM users LIMIT ?",
			qargs: []interface{}{10},
		},
		{
			name: "Simple ORDER BY",
			s: &SelectStatement{
				selections: []types.Selection{users},
				projs:      []types.Projection{colUserName},
				orderBy: &OrderByClause{
					scols: []*SortColumn{colUserName.Desc()},
				},
			},
			qs: "SELECT users.name FROM users ORDER BY users.name DESC",
		},
		{
			name: "Simple GROUP BY",
			s: &SelectStatement{
				selections: []types.Selection{users},
				projs:      []types.Projection{colUserName},
				groupBy: &GroupByClause{
					cols: []types.Projection{colUserName},
				},
			},
			qs: "SELECT users.name FROM users GROUP BY users.name",
		},
		{
			name: "GROUP BY, ORDER BY and LIMIT",
			s: &SelectStatement{
				selections: []types.Selection{users},
				projs:      []types.Projection{colUserName},
				groupBy: &GroupByClause{
					cols: []types.Projection{colUserName},
				},
				orderBy: &OrderByClause{
					scols: []*SortColumn{colUserName.Desc()},
				},
				limit: &LimitClause{limit: 10},
			},
			qs:    "SELECT users.name FROM users GROUP BY users.name ORDER BY users.name DESC LIMIT ?",
			qargs: []interface{}{10},
		},
		{
			name: "Single JOIN",
			s: &SelectStatement{
				selections: []types.Selection{articles},
				projs:      []types.Projection{colArticleId, colUserName.As("author")},
				joins: []*JoinClause{
					&JoinClause{
						left:  articles,
						right: users,
						on:    Equal(colArticleAuthor, colUserId),
					},
				},
			},
			qs: "SELECT articles.id, users.name AS author FROM articles JOIN users ON articles.author = users.id",
		},
		{
			name: "Multiple JOINs",
			s: &SelectStatement{
				selections: []types.Selection{articles},
				projs:      []types.Projection{colArticleId, colUserName.As("author"), colArticleStateName.As("state")},
				joins: []*JoinClause{
					&JoinClause{
						left:  articles,
						right: users,
						on:    Equal(colArticleAuthor, colUserId),
					},
					&JoinClause{
						left:  articles,
						right: article_states,
						on:    Equal(colArticleState, colArticleStateId),
					},
				},
			},
			qs: "SELECT articles.id, users.name AS author, article_states.name AS state FROM articles JOIN users ON articles.author = users.id JOIN article_states ON articles.state = article_states.id",
		},
		{
			name: "COUNT(*) on a table",
			s: &SelectStatement{
				selections: []types.Selection{users},
				projs:      []types.Projection{Count(users)},
			},
			qs: "SELECT COUNT(*) FROM users",
		},
	}
	for _, test := range tests {
		expArgc := len(test.qargs)
		argc := test.s.ArgCount()
		assert.Equal(expArgc, argc)

		expLen := len(test.qs)
		size := test.s.Size(defaultScanner)
		size += interpolationLength(types.DIALECT_MYSQL, argc)
		assert.Equal(expLen, size)

		b := make([]byte, size)
		curArg := 0
		written := test.s.Scan(defaultScanner, b, test.qargs, &curArg)

		assert.Equal(written, size)
		assert.Equal(test.qs, string(b))
	}
}
