//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package expr_test

import (
	"testing"

	"github.com/jaypipes/sqlb/core/expr"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	articles := m.T("articles")
	articleStates := m.T("article_states")
	userProfiles := m.T("user_profiles")
	colUserId := users.C("id")
	colUserName := users.C("name")
	colArticleId := articles.C("id")
	colArticleAuthor := articles.C("author")
	colArticleState := articles.C("state")
	colArticleStateId := articleStates.C("id")
	colArticleStateName := articleStates.C("name")
	colUserProfileContent := userProfiles.C("content")
	colUserProfileUser := userProfiles.C("user")

	subq := expr.Select(colUserId).As("users_derived")

	tests := []struct {
		name  string
		q     *expr.Selection
		qs    string
		qargs []interface{}
	}{
		{
			name: "Simple INNER JOIN",
			q: expr.Select(
				colArticleId, colUserName.As("author"),
			).Join(
				users, expr.Equal(
					colArticleAuthor, colUserId,
				),
			),
			qs: "SELECT articles.id, users.name AS author FROM articles JOIN users ON articles.author = users.id",
		},
		{
			name: "Simple LEFT JOIN",
			q: expr.Select(
				colArticleId, colUserName.As("author"),
			).OuterJoin(
				users, expr.Equal(
					colArticleAuthor, colUserId,
				),
			),
			qs: "SELECT articles.id, users.name AS author FROM articles LEFT JOIN users ON articles.author = users.id",
		},
		{
			name: "JOIN A to B and A to C",
			q: expr.Select(
				colArticleId,
				colUserName.As("author"),
				colArticleStateName.As("state"),
			).Join(
				users, expr.Equal(
					colArticleAuthor, colUserId,
				),
			).Join(
				articleStates, expr.Equal(
					colArticleState, colArticleStateId,
				),
			),
			qs: "SELECT articles.id, users.name AS author, article_states.name AS state FROM articles JOIN users ON articles.author = users.id JOIN article_states ON articles.state = article_states.id",
		},
		{
			name: "JOIN A to B and B to C",
			q: expr.Select(
				colArticleId,
				colUserName.As("author"),
				colUserProfileContent.As("author_profile"),
			).Join(
				users, expr.Equal(
					colArticleAuthor, colUserId,
				),
			).Join(
				userProfiles, expr.Equal(
					colUserId, colUserProfileUser,
				),
			),
			qs: "SELECT articles.id, users.name AS author, user_profiles.content AS author_profile FROM articles JOIN users ON articles.author = users.id JOIN user_profiles ON users.id = user_profiles.user",
		},
		{
			name: "LEFT JOIN with WHERE",
			q: expr.Select(
				colArticleId, colUserName.As("author"),
			).OuterJoin(
				users,
				expr.Equal(colArticleAuthor, colUserId),
			).Where(
				expr.IsNull(colArticleAuthor),
			),
			qs: "SELECT articles.id, users.name AS author FROM articles LEFT JOIN users ON articles.author = users.id WHERE articles.author IS NULL",
		},
		{
			name: "JOIN to derived table (subquery in FROM clause)",
			q: expr.Select(
				colUserId,
				colUserName,
			).Join(
				subq, expr.Equal(
					colUserId, subq.C("id"),
				),
			),
			qs: "SELECT users.id, users.name FROM users JOIN (SELECT users.id FROM users) AS users_derived ON users.id = users_derived.id",
		},
		{
			name: "LEFT JOIN to derived table (subquery in FROM clause)",
			q: expr.Select(
				colUserId,
				colUserName,
			).OuterJoin(
				subq, expr.Equal(colUserId, subq.C("id")),
			),
			qs: "SELECT users.id, users.name FROM users LEFT JOIN (SELECT users.id FROM users) AS users_derived ON users.id = users_derived.id",
		},
		{
			name: "LEFT JOIN to derived table with WHERE referencing derived table",
			q: expr.Select(
				colUserId,
				colUserName,
			).OuterJoin(
				subq, expr.Equal(
					colUserId, subq.C("id"),
				),
			).Where(expr.Equal(subq.C("id"), 1)),
			qs:    "SELECT users.id, users.name FROM users LEFT JOIN (SELECT users.id FROM users) AS users_derived ON users.id = users_derived.id WHERE users_derived.id = ?",
			qargs: []interface{}{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			b := builder.New()

			qs, qargs := b.StringArgs(tt.q.Query())
			assert.Equal(len(tt.qargs), len(qargs))
			assert.Equal(tt.qs, qs)
		})
	}
}

func TestJoinPanics(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	articles := m.T("articles")
	articleStates := m.T("article_states")
	colUserId := users.C("id")
	colArticleAuthor := articles.C("author")

	tests := []struct {
		name string
		q    func()
	}{
		{
			name: "Bad JOIN. Can't Join() against no selection",
			q:    func() { expr.Select().Join(users, expr.Equal(colArticleAuthor, colUserId)) },
		},
		{
			name: "Bad JOIN. Can't Join() against a selection that isn't in the containing SELECT",
			q:    func() { expr.Select(articleStates).Join(users, expr.Equal(colArticleAuthor, colUserId)) },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			assert.Panics(tt.q)
		})
	}
}

func TestNestedSetQueries(t *testing.T) {
	// ref: https://github.com/jaypipes/sqlb/issues/49
	assert := assert.New(t)

	m := testutil.M()
	orgs := m.T("organizations")

	o1 := orgs.As("o1")
	o2 := orgs.As("o2")

	o1id := o1.C("id")
	o2id := o2.C("id")
	o1rootid := o1.C("root_organization_id")
	o2rootid := o2.C("root_organization_id")
	o1nestedleft := o1.C("nested_set_left")
	o2nestedleft := o2.C("nested_set_left")
	o2nestedright := o2.C("nested_set_right")

	joinCond := expr.And(
		expr.Equal(o1rootid, o2rootid),
		expr.Between(o1nestedleft, o2nestedleft, o2nestedright),
	)
	q := expr.Select(o1id).Join(o2, joinCond)
	q.Where(expr.Equal(o2id, 2))

	b := builder.New()
	qs, qargs := b.StringArgs(q.Query())

	expqs := "SELECT o1.id FROM organizations AS o1 JOIN organizations AS o2 ON o1.root_organization_id = o2.root_organization_id AND o1.nested_set_left BETWEEN o2.nested_set_left AND o2.nested_set_right WHERE o2.id = ?"
	expqargs := []interface{}{2}

	assert.Equal(expqs, qs)
	assert.Equal(expqargs, qargs)
}

func TestNestedSetWithAdditionalJoin(t *testing.T) {
	// ref: https://github.com/jaypipes/sqlb/issues/60
	assert := assert.New(t)

	m := testutil.M()
	orgs := m.T("organizations")
	orgUsers := m.T("organization_users")

	o1 := orgs.As("o1")
	o2 := orgs.As("o2")
	ou := orgUsers.As("ou")

	o1id := o1.C("id")
	o2id := o2.C("id")
	o1rootid := o1.C("root_organization_id")
	o2rootid := o2.C("root_organization_id")
	o1nestedleft := o1.C("nested_set_left")
	o2nestedleft := o2.C("nested_set_left")
	o2nestedright := o2.C("nested_set_right")
	ouUserId := ou.C("user_id")
	ouOrgId := ou.C("organization_id")

	nestedJoinCond := expr.And(
		expr.Equal(o1rootid, o2rootid),
		expr.Between(o1nestedleft, o2nestedleft, o2nestedright),
	)
	ouJoinCond := expr.And(
		expr.Equal(o2id, ouOrgId),
		expr.Equal(ouUserId, 1),
	)
	q := expr.Select(o1id)
	q.Join(o2, nestedJoinCond)
	q.Join(ou, ouJoinCond)

	b := builder.New()
	qs, qargs := b.StringArgs(q.Query())

	expqs := "SELECT o1.id FROM organizations AS o1 JOIN organizations AS o2 ON o1.root_organization_id = o2.root_organization_id AND o1.nested_set_left BETWEEN o2.nested_set_left AND o2.nested_set_right JOIN organization_users AS ou ON o2.id = ou.organization_id AND ou.user_id = ?"
	expqargs := []interface{}{1}

	assert.Equal(expqs, qs)
	assert.Equal(expqargs, qargs)
}

func TestJoinDerivedWithMultipleSelections(t *testing.T) {
	// ref: https://github.com/jaypipes/sqlb/issues/68
	assert := assert.New(t)

	// The SQL we want to generate looks like this:
	//
	// SELECT
	//   o.uuid
	// FROM organizations AS o
	// LEFT JOIN organizations AS po
	//   ON o.parent_organization_id = po.id
	// LEFT JOIN (
	//   SELECT o1.id
	//   FROM organizations AS o1
	//   JOIN organizations AS o2
	//     ON o1.root_organization_id = o2.root_organization_id
	//     AND o1.nested_set_left BETWEEN o2.nested_set_left AND o2.nested_set_right
	//   JOIN organization_users AS ou
	//     ON o2.id = ou.organization_id
	//     AND ou.user_id = ?
	// ) AS private_orgs
	//   ON o.id = private_orgs.id
	// WHERE
	//   o.visibility = 1
	//   OR (o.visibility = 0 AND private_orgs.id IS NOT NULL)

	m := testutil.M()
	orgs := m.T("organizations")
	orgUsers := m.T("organization_users")

	o1 := orgs.As("o1")
	o2 := orgs.As("o2")
	ou := orgUsers.As("ou")

	o1id := o1.C("id")
	o2id := o2.C("id")
	o1rootid := o1.C("root_organization_id")
	o2rootid := o2.C("root_organization_id")
	o1nestedleft := o1.C("nested_set_left")
	o2nestedleft := o2.C("nested_set_left")
	o2nestedright := o2.C("nested_set_right")
	ouUserId := ou.C("user_id")
	ouOrgId := ou.C("organization_id")

	nestedJoinCond := expr.And(
		expr.Equal(o1rootid, o2rootid),
		expr.Between(o1nestedleft, o2nestedleft, o2nestedright),
	)
	ouJoin := expr.And(
		expr.Equal(o2id, ouOrgId),
		expr.Equal(ouUserId, 1),
	)
	subq := expr.Select(o1id).Join(o2, nestedJoinCond).Join(ou, ouJoin).As("derived")

	subqOrgId := subq.C("id")
	q := expr.Select(subqOrgId)

	b := builder.New()
	qs, qargs := b.StringArgs(q.Query())

	expqs := "SELECT derived.id FROM (SELECT o1.id FROM organizations AS o1 JOIN organizations AS o2 ON o1.root_organization_id = o2.root_organization_id AND o1.nested_set_left BETWEEN o2.nested_set_left AND o2.nested_set_right JOIN organization_users AS ou ON o2.id = ou.organization_id AND ou.user_id = ?) AS derived"
	expqargs := []interface{}{1}

	assert.Equal(expqs, qs)
	assert.Equal(expqargs, qargs)

	b = builder.New()
	q = expr.Select(
		orgs.C("uuid"),
	).OuterJoin(
		subq,
		expr.Equal(
			orgs.C("id"),
			subqOrgId,
		),
	).Where(expr.IsNotNull(subqOrgId))

	qs, qargs = b.StringArgs(q.Query())

	expqs = "SELECT organizations.uuid FROM organizations LEFT JOIN (SELECT o1.id FROM organizations AS o1 JOIN organizations AS o2 ON o1.root_organization_id = o2.root_organization_id AND o1.nested_set_left BETWEEN o2.nested_set_left AND o2.nested_set_right JOIN organization_users AS ou ON o2.id = ou.organization_id AND ou.user_id = ?) AS derived ON organizations.id = derived.id WHERE derived.id IS NOT NULL"
	expqargs = []interface{}{1}

	assert.Equal(expqs, qs)
	assert.Equal(expqargs, qargs)
}
