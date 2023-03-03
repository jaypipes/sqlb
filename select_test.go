// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package sqlb

import (
	"fmt"
	"testing"

	"github.com/jaypipes/sqlb/pkg/errors"
	"github.com/jaypipes/sqlb/pkg/grammar/expression"
	"github.com/jaypipes/sqlb/pkg/grammar/function"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestSelectQuery(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := T(sc, "users")
	articles := T(sc, "articles")
	articleStates := T(sc, "article_states")
	userProfiles := T(sc, "user_profiles")
	colUserId := users.C("id")
	colUserName := users.C("name")
	colArticleId := articles.C("id")
	colArticleAuthor := articles.C("author")
	colArticleState := articles.C("state")
	colArticleStateId := articleStates.C("id")
	colArticleStateName := articleStates.C("name")
	colUserProfileContent := userProfiles.C("content")
	colUserProfileUser := userProfiles.C("user")

	subq := Select(colUserId).As("users_derived")

	tests := []struct {
		name  string
		q     *SelectQuery
		qs    string
		qargs []interface{}
		qe    error
	}{
		{
			name: "Simple FROM",
			q:    Select(users),
			qs:   "SELECT users.id, users.name FROM users",
		},
		{
			name: "Simple SELECT COUNT(*) FROM",
			q:    Select(function.Count(users)),
			qs:   "SELECT COUNT(*) FROM users",
		},
		{
			name:  "Simple WHERE",
			q:     Select(users).Where(expression.Equal(colUserName, "foo")),
			qs:    "SELECT users.id, users.name FROM users WHERE users.name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "WHERE with an OR expression",
			q: Select(users).Where(
				expression.Or(
					expression.Equal(colUserName, "foo"),
					expression.Equal(colUserName, "bar"),
				),
			),
			qs:    "SELECT users.id, users.name FROM users WHERE (users.name = ? OR users.name = ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "Simple GROUP BY",
			q:    Select(users).GroupBy(colUserName),
			qs:   "SELECT users.id, users.name FROM users GROUP BY users.name",
		},
		{
			name:  "Simple HAVING",
			q:     Select(users).Having(expression.Equal(colUserName, "foo")),
			qs:    "SELECT users.id, users.name FROM users HAVING users.name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "Simple ORDER BY",
			q:    Select(users).OrderBy(colUserName.Desc()),
			qs:   "SELECT users.id, users.name FROM users ORDER BY users.name DESC",
		},
		{
			name:  "Simple LIMIT",
			q:     Select(users).Limit(10),
			qs:    "SELECT users.id, users.name FROM users LIMIT ?",
			qargs: []interface{}{10},
		},
		{
			name:  "Simple LIMIT with OFFSET",
			q:     Select(users).LimitWithOffset(10, 20),
			qs:    "SELECT users.id, users.name FROM users LIMIT ? OFFSET ?",
			qargs: []interface{}{10, 20},
		},
		{
			name: "Simple named derived table",
			q:    Select(Select(users).As("u")),
			qs:   "SELECT u.id, u.name FROM (SELECT users.id, users.name FROM users) AS u",
		},
		{
			name: "Simple un-named derived table",
			q:    Select(Select(users)),
			qs:   "SELECT derived0.id, derived0.name FROM (SELECT users.id, users.name FROM users) AS derived0",
		},
		{
			name: "Bad JOIN. Can't Join() against no selection",
			q:    Select().Join(users, expression.Equal(colArticleAuthor, colUserId)),
			qe:   errors.InvalidJoinNoSelect,
		},
		{
			name: "Bad JOIN. Can't Join() against a selection that isn't in the containing SELECT",
			q:    Select(articleStates).Join(users, expression.Equal(colArticleAuthor, colUserId)),
			qe:   errors.InvalidJoinUnknownTarget,
		},
		{
			name: "Simple INNER JOIN",
			q:    Select(colArticleId, colUserName.As("author")).Join(users, expression.Equal(colArticleAuthor, colUserId)),
			qs:   "SELECT articles.id, users.name AS author FROM articles JOIN users ON articles.author = users.id",
		},
		{
			name: "Simple LEFT JOIN",
			q:    Select(colArticleId, colUserName.As("author")).OuterJoin(users, expression.Equal(colArticleAuthor, colUserId)),
			qs:   "SELECT articles.id, users.name AS author FROM articles LEFT JOIN users ON articles.author = users.id",
		},
		{
			name: "JOIN A to B and A to C",
			q: Select(
				colArticleId,
				colUserName.As("author"),
				colArticleStateName.As("state"),
			).Join(users, expression.Equal(colArticleAuthor, colUserId)).Join(articleStates, expression.Equal(colArticleState, colArticleStateId)),
			qs: "SELECT articles.id, users.name AS author, article_states.name AS state FROM articles JOIN users ON articles.author = users.id JOIN article_states ON articles.state = article_states.id",
		},
		{
			name: "LEFT JOIN with WHERE",
			q: Select(
				colArticleId, colUserName.As("author"),
			).OuterJoin(
				users,
				expression.Equal(colArticleAuthor, colUserId),
			).Where(
				expression.IsNull(colArticleAuthor),
			),
			qs: "SELECT articles.id, users.name AS author FROM articles LEFT JOIN users ON articles.author = users.id WHERE articles.author IS NULL",
		},
		{
			name: "LEFT JOIN to derived table (subquery in FROM clause)",
			q: Select(
				colUserId,
				colUserName,
			).OuterJoin(subq, expression.Equal(colUserId, subq.C("id"))),
			qs: "SELECT users.id, users.name FROM users LEFT JOIN (SELECT users.id FROM users) AS users_derived ON users.id = users_derived.id",
		},
		{
			name: "JOIN to derived table (subquery in FROM clause)",
			q: Select(
				colUserId,
				colUserName,
			).Join(subq, expression.Equal(colUserId, subq.C("id"))),
			qs: "SELECT users.id, users.name FROM users JOIN (SELECT users.id FROM users) AS users_derived ON users.id = users_derived.id",
		},
		{
			name: "JOIN A to B and B to C",
			q: Select(
				colArticleId,
				colUserName.As("author"),
				colUserProfileContent.As("author_profile"),
			).Join(users, expression.Equal(colArticleAuthor, colUserId)).Join(userProfiles, expression.Equal(colUserId, colUserProfileUser)),
			qs: "SELECT articles.id, users.name AS author, user_profiles.content AS author_profile FROM articles JOIN users ON articles.author = users.id JOIN user_profiles ON users.id = user_profiles.user",
		},
		{
			name: "LEFT JOIN to derived table with WHERE referencing derived table",
			q: Select(
				colUserId,
				colUserName,
			).OuterJoin(subq, expression.Equal(colUserId, subq.C("id"))).Where(expression.Equal(subq.C("id"), 1)),
			qs:    "SELECT users.id, users.name FROM users LEFT JOIN (SELECT users.id FROM users) AS users_derived ON users.id = users_derived.id WHERE users_derived.id = ?",
			qargs: []interface{}{1},
		},
	}
	for _, test := range tests {
		if test.qe != nil {
			assert.Equal(test.qe, test.q.Error())
			continue
		} else if test.q.Error() != nil {
			qe := test.q.Error()
			assert.Fail(qe.Error())
			continue
		}
		sc := scanner.DefaultScanner
		qs, qargs := sc.StringArgs(test.q)
		assert.Equal(len(test.qargs), len(qargs))
		assert.Equal(test.qs, qs)
	}
}

func TestNestedSetQueries(t *testing.T) {
	// ref: https://github.com/jaypipes/sqlb/issues/49
	assert := assert.New(t)

	sc := testutil.Schema()
	orgs := T(sc, "organizations")

	o1 := orgs.As("o1")
	o2 := orgs.As("o2")

	o1id := o1.C("id")
	o2id := o2.C("id")
	o1rootid := o1.C("root_organization_id")
	o2rootid := o2.C("root_organization_id")
	o1nestedleft := o1.C("nested_set_left")
	o2nestedleft := o2.C("nested_set_left")
	o2nestedright := o2.C("nested_set_right")

	joinCond := expression.And(
		expression.Equal(o1rootid, o2rootid),
		expression.Between(o1nestedleft, o2nestedleft, o2nestedright),
	)
	q := Select(o1id).Join(o2, joinCond)
	q.Where(expression.Equal(o2id, 2))

	scan := scanner.DefaultScanner
	qs, qargs := scan.StringArgs(q)

	expqs := "SELECT o1.id FROM organizations AS o1 JOIN organizations AS o2 ON (o1.root_organization_id = o2.root_organization_id AND o1.nested_set_left BETWEEN o2.nested_set_left AND o2.nested_set_right) WHERE o2.id = ?"
	expqargs := []interface{}{2}

	assert.Equal(expqs, qs)
	assert.Equal(expqargs, qargs)
}

func TestNestedSetWithAdditionalJoin(t *testing.T) {
	// ref: https://github.com/jaypipes/sqlb/issues/60
	assert := assert.New(t)

	sc := testutil.Schema()
	orgs := T(sc, "organizations")
	orgUsers := T(sc, "organization_users")

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

	nestedJoinCond := expression.And(
		expression.Equal(o1rootid, o2rootid),
		expression.Between(o1nestedleft, o2nestedleft, o2nestedright),
	)
	ouJoin := expression.And(
		expression.Equal(o2id, ouOrgId),
		expression.Equal(ouUserId, 1),
	)
	q := Select(o1id).Join(o2, nestedJoinCond).Join(ou, ouJoin)

	assert.Nil(q.e)

	scan := scanner.DefaultScanner
	qs, qargs := scan.StringArgs(q)

	expqs := "SELECT o1.id FROM organizations AS o1 JOIN organizations AS o2 ON (o1.root_organization_id = o2.root_organization_id AND o1.nested_set_left BETWEEN o2.nested_set_left AND o2.nested_set_right) JOIN organization_users AS ou ON (o2.id = ou.organization_id AND ou.user_id = ?)"
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
	// WHERE (
	//   o.visibility = 1
	//   OR (o.visibility = 0 AND private_orgs.id IS NOT NULL)

	sc := testutil.Schema()
	orgs := T(sc, "organizations")
	orgUsers := T(sc, "organization_users")

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

	nestedJoinCond := expression.And(
		expression.Equal(o1rootid, o2rootid),
		expression.Between(o1nestedleft, o2nestedleft, o2nestedright),
	)
	ouJoin := expression.And(
		expression.Equal(o2id, ouOrgId),
		expression.Equal(ouUserId, 1),
	)
	subq := Select(o1id).Join(o2, nestedJoinCond).Join(ou, ouJoin).As("derived")
	subqOrgId := subq.C("id")

	assert.Nil(subq.e)

	scan := scanner.DefaultScanner
	qs, qargs := scan.StringArgs(subq)

	expqs := "SELECT derived.id FROM (SELECT o1.id FROM organizations AS o1 JOIN organizations AS o2 ON (o1.root_organization_id = o2.root_organization_id AND o1.nested_set_left BETWEEN o2.nested_set_left AND o2.nested_set_right) JOIN organization_users AS ou ON (o2.id = ou.organization_id AND ou.user_id = ?)) AS derived"
	expqargs := []interface{}{1}

	assert.Equal(expqs, qs)
	assert.Equal(expqargs, qargs)

	q := Select(
		orgs.C("uuid"),
	).OuterJoin(
		subq,
		expression.Equal(
			orgs.C("id"),
			subqOrgId,
		),
	).Where(expression.IsNotNull(subqOrgId))

	assert.Nil(q.e)

	qs, qargs = scan.StringArgs(q)

	expqs = "SELECT organizations.uuid FROM organizations LEFT JOIN (SELECT o1.id FROM organizations AS o1 JOIN organizations AS o2 ON (o1.root_organization_id = o2.root_organization_id AND o1.nested_set_left BETWEEN o2.nested_set_left AND o2.nested_set_right) JOIN organization_users AS ou ON (o2.id = ou.organization_id AND ou.user_id = ?)) AS derived ON organizations.id = derived.id WHERE derived.id IS NOT NULL"
	expqargs = []interface{}{1}

	assert.Equal(expqs, qs)
	assert.Equal(expqargs, qargs)
}

func TestModifyingSelectQueryUpdatesBuffer(t *testing.T) {
	assert := assert.New(t)

	scan := scanner.DefaultScanner
	sc := testutil.Schema()
	users := T(sc, "users")

	q := Select(users)

	qs, qargs := scan.StringArgs(q)
	assert.Equal("SELECT users.id, users.name FROM users", qs)
	assert.Empty(qargs)

	// Modify the underlying SELECT and verify string and args changed
	q.Where(expression.Equal(users.C("id"), 1))
	qs, qargs = scan.StringArgs(q)
	assert.Equal("SELECT users.id, users.name FROM users WHERE users.id = ?", qs)
	assert.Equal([]interface{}{1}, qargs)
}

func TestSelectQueryErrors(t *testing.T) {
	assert := assert.New(t)

	q := &SelectQuery{}

	assert.False(q.IsValid()) // Doesn't have a selectClause yet...
	assert.Nil(q.Error())     // But there is no error set yet...

	sc := testutil.Schema()
	users := T(sc, "users")

	q = Select(users)

	assert.True(q.IsValid())
	assert.Nil(q.Error())

	q.e = fmt.Errorf("Cannot determine left side of JOIN expression.")
	assert.False(q.IsValid())
	assert.NotNil(q.Error())
}
