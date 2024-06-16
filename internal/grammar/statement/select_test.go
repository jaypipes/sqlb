//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package statement_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar/clause"
	"github.com/jaypipes/sqlb/internal/grammar/element"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
	"github.com/jaypipes/sqlb/internal/grammar/function"
	"github.com/jaypipes/sqlb/internal/grammar/statement"
	"github.com/jaypipes/sqlb/internal/testutil"
)

func TestSelectStatement(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	articles := sqlb.T(m, "articles")
	article_states := sqlb.T(m, "article_states")
	colUserName := users.C("name")
	colUserId := users.C("id")
	colArticleId := articles.C("id")
	colArticleAuthor := articles.C("author")
	colArticleState := articles.C("state")
	colArticleStateId := article_states.C("id")
	colArticleStateName := article_states.C("name")

	tests := []struct {
		name  string
		s     *statement.Select
		qs    string
		qargs []interface{}
	}{
		{
			name: "A literal value",
			s: statement.NewSelect(
				[]builder.Projection{element.NewValue(nil, 1)},
				nil, nil, nil, nil, nil, nil, nil,
			),
			qs:    "SELECT ?",
			qargs: []interface{}{1},
		},
		{
			name: "A literal value aliased",
			s: statement.NewSelect(
				[]builder.Projection{element.NewValue(nil, 1).As("foo")},
				nil, nil, nil, nil, nil, nil, nil,
			),
			qs:    "SELECT ? AS foo",
			qargs: []interface{}{1},
		},
		{
			name: "Two literal values",
			s: statement.NewSelect(
				[]builder.Projection{
					element.NewValue(nil, 1),
					element.NewValue(nil, 2),
				},
				nil, nil, nil, nil, nil, nil, nil,
			),
			qs:    "SELECT ?, ?",
			qargs: []interface{}{1, 2},
		},
		{
			name: "Table and column",
			s: statement.NewSelect(
				[]builder.Projection{colUserName},
				[]builder.Selection{users},
				nil, nil, nil, nil, nil, nil,
			),
			qs: "SELECT users.name FROM users",
		},
		{
			name: "aliased Table and Column",
			s: statement.NewSelect(
				[]builder.Projection{
					users.As("u").C("name"),
				},
				[]builder.Selection{users.As("u")},
				nil, nil, nil, nil, nil, nil,
			),
			qs: "SELECT u.name FROM users AS u",
		},
		{
			name: "Table and multiple Column",
			s: statement.NewSelect(
				[]builder.Projection{colUserId, colUserName},
				[]builder.Selection{users},
				nil, nil, nil, nil, nil, nil,
			),
			qs: "SELECT users.id, users.name FROM users",
		},
		{
			name: "Simple WHERE",
			s: statement.NewSelect(
				[]builder.Projection{colUserName},
				[]builder.Selection{users},
				nil,
				clause.NewWhere(
					expression.Equal(colUserName, "foo"),
				),
				nil, nil, nil, nil,
			),
			qs:    "SELECT users.name FROM users WHERE users.name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "Simple LIMIT",
			s: statement.NewSelect(
				[]builder.Projection{colUserName},
				[]builder.Selection{users},
				nil, nil, nil, nil, nil,
				clause.NewLimit(10, nil),
			),
			qs:    "SELECT users.name FROM users LIMIT ?",
			qargs: []interface{}{10},
		},
		{
			name: "Simple ORDER BY",
			s: statement.NewSelect(
				[]builder.Projection{colUserName},
				[]builder.Selection{users},
				nil, nil, nil, nil,
				clause.NewOrderBy(colUserName.Desc()),
				nil,
			),
			qs: "SELECT users.name FROM users ORDER BY users.name DESC",
		},
		{
			name: "Simple GROUP BY",
			s: statement.NewSelect(
				[]builder.Projection{colUserName},
				[]builder.Selection{users},
				nil, nil,
				clause.NewGroupBy(colUserName),
				nil, nil, nil,
			),
			qs: "SELECT users.name FROM users GROUP BY users.name",
		},
		{
			name: "GROUP BY, ORDER BY and LIMIT",
			s: statement.NewSelect(
				[]builder.Projection{colUserName},
				[]builder.Selection{users},
				nil, nil,
				clause.NewGroupBy(colUserName),
				nil,
				clause.NewOrderBy(colUserName.Desc()),
				clause.NewLimit(10, nil),
			),
			qs:    "SELECT users.name FROM users GROUP BY users.name ORDER BY users.name DESC LIMIT ?",
			qargs: []interface{}{10},
		},
		{
			name: "Single JOIN",
			s: statement.NewSelect(
				[]builder.Projection{colArticleId, colUserName.As("author")},
				[]builder.Selection{articles},
				[]*clause.Join{
					clause.InnerJoin(
						articles,
						users,
						expression.Equal(colArticleAuthor, colUserId),
					),
				},
				nil, nil, nil, nil, nil,
			),
			qs: "SELECT articles.id, users.name AS author FROM articles JOIN users ON articles.author = users.id",
		},
		{
			name: "Multiple JOINs",
			s: statement.NewSelect(
				[]builder.Projection{colArticleId, colUserName.As("author"), colArticleStateName.As("state")},
				[]builder.Selection{articles},
				[]*clause.Join{
					clause.InnerJoin(
						articles,
						users,
						expression.Equal(colArticleAuthor, colUserId),
					),
					clause.InnerJoin(
						articles,
						article_states,
						expression.Equal(colArticleState, colArticleStateId),
					),
				},
				nil, nil, nil, nil, nil,
			),
			qs: "SELECT articles.id, users.name AS author, article_states.name AS state FROM articles JOIN users ON articles.author = users.id JOIN article_states ON articles.state = article_states.id",
		},
		{
			name: "COUNT(*) on a table",
			s: statement.NewSelect(
				[]builder.Projection{function.Count(users)},
				[]builder.Selection{users},
				nil, nil, nil, nil, nil, nil,
			),
			qs: "SELECT COUNT(*) FROM users",
		},
	}
	for _, test := range tests {
		b := builder.New()

		expArgc := len(test.qargs)
		argc := test.s.ArgCount()
		assert.Equal(expArgc, argc)

		expLen := len(test.qs)
		size := test.s.Size(b)
		size += b.InterpolationLength(argc)
		assert.Equal(expLen, size)

		curArg := 0
		test.s.Scan(b, test.qargs, &curArg)

		assert.Equal(test.qs, b.String())
	}
}
