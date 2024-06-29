//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sqlb_test

import (
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar/identifier"
	"github.com/jaypipes/sqlb/internal/query"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")

	tests := []struct {
		name   string
		t      *identifier.Table
		values map[string]interface{}
		qs     string
		qargs  []interface{}
		qe     error
	}{
		{
			name:   "Table missing",
			t:      nil,
			values: map[string]interface{}{"unknown": 1},
			qe:     api.TableRequired,
		},
		{
			name:   "Values missing",
			t:      users,
			values: nil,
			qe:     api.NoValues,
		},
		{
			name:   "Unknown column",
			t:      users,
			values: map[string]interface{}{"unknown": 1},
			qe:     api.UnknownColumn,
		},
		{
			name:   "Simple INSERT",
			t:      users,
			values: map[string]interface{}{"id": 1},
			qs:     "INSERT INTO users (id) VALUES (?)",
			qargs:  []interface{}{1},
		},
		//{
		//	name:  "INSERT using Table.Insert() adapter",
		//	q:     users.Insert(map[string]interface{}{"id": 1}),
		//	qs:    "INSERT INTO users (id) VALUES (?)",
		//	qargs: []interface{}{1},
		//},
	}
	for _, tt := range tests {
		got, err := sqlb.Insert(tt.t, tt.values)
		if tt.qe != nil {
			assert.Equal(tt.qe, err)
			continue
		}
		assert.Nil(err)
		b := builder.New()
		qs, qargs := b.StringArgs(got)
		assert.Equal(len(tt.qargs), len(qargs))
		assert.Equal(tt.qs, qs)
	}
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	//colUserName := users.C("name")

	tests := []struct {
		name  string
		t     *identifier.Table
		qs    string
		qargs []interface{}
		qe    error
	}{
		{
			name: "No target table",
			qe:   api.TableRequired,
		},
		{
			name: "DELETE all rows",
			t:    users,
			qs:   "DELETE FROM users",
		},
		//{
		//	name: "Table.Delete() variant",
		//	q:    users.Delete(),
		//	qs:   "DELETE FROM users",
		//},
		//{
		//	name:  "DELETE simple WHERE",
		//	q:     query.Delete(users).Where(sqlb.Equal(colUserName, "foo")),
		//	qs:    "DELETE FROM users WHERE users.name = ?",
		//	qargs: []interface{}{"foo"},
		//},
	}
	for _, tt := range tests {
		got, err := sqlb.Delete(tt.t)
		if tt.qe != nil {
			assert.Equal(tt.qe, err)
			continue
		}
		assert.Nil(err)
		b := builder.New()
		qs, qargs := b.StringArgs(got)
		assert.Equal(len(tt.qargs), len(qargs))
		assert.Equal(tt.qs, qs)
	}
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	//colUserName := users.C("name")

	tests := []struct {
		name   string
		t      *identifier.Table
		values map[string]interface{}
		qs     string
		qargs  []interface{}
		qe     error
	}{
		{
			name:   "Values missing",
			t:      users,
			values: nil,
			qe:     api.NoValues,
		},
		{
			name:   "Target table missing",
			t:      nil,
			values: map[string]interface{}{"name": "foo"},
			qe:     api.TableRequired,
		},
		{
			name:   "Unknown column",
			t:      users,
			values: map[string]interface{}{"unknown": 1},
			qe:     api.UnknownColumn,
		},
		{
			name:   "UPDATE no WHERE",
			t:      users,
			values: map[string]interface{}{"name": "foo"},
			qs:     "UPDATE users SET name = ?",
			qargs:  []interface{}{"foo"},
		},
		//{
		//	name:  "UPDATE no WHERE using Table.Update()",
		//	q:     users.Update(map[string]interface{}{"name": "foo"}),
		//	qs:    "UPDATE users SET name = ?",
		//	qargs: []interface{}{"foo"},
		//},
		//{
		//	name: "UPDATE simple WHERE",
		//	q: query.Update(users, map[string]interface{}{"name": "bar"}).Where(
		//		sqlb.Equal(colUserName, "foo"),
		//	),
		//	qs:    "UPDATE users SET name = ? WHERE users.name = ?",
		//	qargs: []interface{}{"bar", "foo"},
		//},
	}
	for _, tt := range tests {
		got, err := sqlb.Update(tt.t, tt.values)
		if tt.qe != nil {
			assert.Equal(tt.qe, err)
			continue
		}
		assert.Nil(err)
		b := builder.New()
		qs, qargs := b.StringArgs(got)
		assert.Equal(len(tt.qargs), len(qargs))
		assert.Equal(tt.qs, qs)
	}
}

func TestSelect(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	articles := sqlb.T(m, "articles")
	articleStates := sqlb.T(m, "article_states")
	userProfiles := sqlb.T(m, "user_profiles")
	colUserId := users.C("id")
	colUserName := users.C("name")
	colArticleId := articles.C("id")
	colArticleAuthor := articles.C("author")
	colArticleState := articles.C("state")
	colArticleStateId := articleStates.C("id")
	colArticleStateName := articleStates.C("name")
	colUserProfileContent := userProfiles.C("content")
	colUserProfileUser := userProfiles.C("user")

	subq := sqlb.Select(colUserId).As("users_derived")

	var sq *query.SelectQuery

	tests := []struct {
		name     string
		q        func()
		qs       string
		qargs    []interface{}
		expPanic bool
	}{
		{
			name: "Simple FROM",
			q:    func() { sq = sqlb.Select(users) },
			qs:   "SELECT users.id, users.name FROM users",
		},
		{
			name: "Simple SELECT COUNT(*) FROM",
			q:    func() { sq = sqlb.Select(sqlb.Count(users)) },
			qs:   "SELECT COUNT(*) FROM users",
		},
		{
			name:  "Simple WHERE",
			q:     func() { sq = sqlb.Select(users).Where(sqlb.Equal(colUserName, "foo")) },
			qs:    "SELECT users.id, users.name FROM users WHERE users.name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "WHERE with an OR sqlb",
			q: func() {
				sq = sqlb.Select(users).Where(
					sqlb.Or(
						sqlb.Equal(colUserName, "foo"),
						sqlb.Equal(colUserName, "bar"),
					),
				)
			},
			qs:    "SELECT users.id, users.name FROM users WHERE (users.name = ? OR users.name = ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "Simple GROUP BY",
			q: func() {
				sq = sqlb.Select(users).GroupBy(colUserName)
			},
			qs: "SELECT users.id, users.name FROM users GROUP BY users.name",
		},
		{
			name: "Simple HAVING",
			q: func() {
				sq = sqlb.Select(users).Having(sqlb.Equal(colUserName, "foo"))
			},
			qs:    "SELECT users.id, users.name FROM users HAVING users.name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "Simple ORDER BY",
			q: func() {
				sq = sqlb.Select(users).OrderBy(colUserName.Desc())
			},
			qs: "SELECT users.id, users.name FROM users ORDER BY users.name DESC",
		},
		{
			name: "Simple LIMIT",
			q: func() {
				sq = sqlb.Select(users).Limit(10)
			},
			qs:    "SELECT users.id, users.name FROM users LIMIT ?",
			qargs: []interface{}{10},
		},
		{
			name: "Simple LIMIT with OFFSET",
			q: func() {
				sq = sqlb.Select(users).LimitWithOffset(10, 20)
			},
			qs:    "SELECT users.id, users.name FROM users LIMIT ? OFFSET ?",
			qargs: []interface{}{10, 20},
		},
		{
			name: "Simple named derived table",
			q: func() {
				sq = sqlb.Select(sqlb.Select(users).As("u"))
			},
			qs: "SELECT u.id, u.name FROM (SELECT users.id, users.name FROM users) AS u",
		},
		{
			name: "Simple un-named derived table",
			q: func() {
				sq = sqlb.Select(sqlb.Select(users))
			},
			qs: "SELECT derived0.id, derived0.name FROM (SELECT users.id, users.name FROM users) AS derived0",
		},
		{
			name: "Bad JOIN. Can't Join() against no selection",
			q: func() {
				sq = sqlb.Select().Join(users, sqlb.Equal(colArticleAuthor, colUserId))
			},
			expPanic: true,
		},
		{
			name: "Bad JOIN. Can't Join() against a selection that isn't in the containing SELECT",
			q: func() {
				sq = sqlb.Select(articleStates).Join(users, sqlb.Equal(colArticleAuthor, colUserId))
			},
			expPanic: true,
		},
		{
			name: "Simple INNER JOIN",
			q: func() {
				sq = sqlb.Select(colArticleId, colUserName.As("author")).Join(users, sqlb.Equal(colArticleAuthor, colUserId))
			},
			qs: "SELECT articles.id, users.name AS author FROM articles JOIN users ON articles.author = users.id",
		},
		{
			name: "Simple LEFT JOIN",
			q: func() {
				sq = sqlb.Select(colArticleId, colUserName.As("author")).OuterJoin(users, sqlb.Equal(colArticleAuthor, colUserId))
			},
			qs: "SELECT articles.id, users.name AS author FROM articles LEFT JOIN users ON articles.author = users.id",
		},
		{
			name: "JOIN A to B and A to C",
			q: func() {
				sq = sqlb.Select(
					colArticleId,
					colUserName.As("author"),
					colArticleStateName.As("state"),
				).Join(users, sqlb.Equal(colArticleAuthor, colUserId)).Join(articleStates, sqlb.Equal(colArticleState, colArticleStateId))
			},
			qs: "SELECT articles.id, users.name AS author, article_states.name AS state FROM articles JOIN users ON articles.author = users.id JOIN article_states ON articles.state = article_states.id",
		},
		{
			name: "LEFT JOIN with WHERE",
			q: func() {
				sq = sqlb.Select(
					colArticleId, colUserName.As("author"),
				).OuterJoin(
					users,
					sqlb.Equal(colArticleAuthor, colUserId),
				).Where(
					sqlb.IsNull(colArticleAuthor),
				)
			},
			qs: "SELECT articles.id, users.name AS author FROM articles LEFT JOIN users ON articles.author = users.id WHERE articles.author IS NULL",
		},
		{
			name: "LEFT JOIN to derived table (subquery in FROM clause)",
			q: func() {
				sq = sqlb.Select(
					colUserId,
					colUserName,
				).OuterJoin(subq, sqlb.Equal(colUserId, subq.C("id")))
			},
			qs: "SELECT users.id, users.name FROM users LEFT JOIN (SELECT users.id FROM users) AS users_derived ON users.id = users_derived.id",
		},
		{
			name: "JOIN to derived table (subquery in FROM clause)",
			q: func() {
				sq = sqlb.Select(
					colUserId,
					colUserName,
				).Join(subq, sqlb.Equal(colUserId, subq.C("id")))
			},
			qs: "SELECT users.id, users.name FROM users JOIN (SELECT users.id FROM users) AS users_derived ON users.id = users_derived.id",
		},
		{
			name: "JOIN A to B and B to C",
			q: func() {
				sq = sqlb.Select(
					colArticleId,
					colUserName.As("author"),
					colUserProfileContent.As("author_profile"),
				).Join(users, sqlb.Equal(colArticleAuthor, colUserId)).Join(userProfiles, sqlb.Equal(colUserId, colUserProfileUser))
			},
			qs: "SELECT articles.id, users.name AS author, user_profiles.content AS author_profile FROM articles JOIN users ON articles.author = users.id JOIN user_profiles ON users.id = user_profiles.user",
		},
		{
			name: "LEFT JOIN to derived table with WHERE referencing derived table",
			q: func() {
				sq = sqlb.Select(
					colUserId,
					colUserName,
				).OuterJoin(subq, sqlb.Equal(colUserId, subq.C("id"))).Where(sqlb.Equal(subq.C("id"), 1))
			},
			qs:    "SELECT users.id, users.name FROM users LEFT JOIN (SELECT users.id FROM users) AS users_derived ON users.id = users_derived.id WHERE users_derived.id = ?",
			qargs: []interface{}{1},
		},
	}
	for _, tt := range tests {
		if tt.expPanic {
			assert.Panics(func() { tt.q() })
			continue
		}
		b := builder.New()

		tt.q()

		qs, qargs := b.StringArgs(sq)
		assert.Equal(len(tt.qargs), len(qargs))
		assert.Equal(tt.qs, qs)
	}
}

func TestNestedSetQueries(t *testing.T) {
	// ref: https://github.com/jaypipes/sqlb/issues/49
	assert := assert.New(t)

	m := testutil.Meta()
	orgs := sqlb.T(m, "organizations")

	o1 := orgs.As("o1")
	o2 := orgs.As("o2")

	o1id := o1.C("id")
	o2id := o2.C("id")
	o1rootid := o1.C("root_organization_id")
	o2rootid := o2.C("root_organization_id")
	o1nestedleft := o1.C("nested_set_left")
	o2nestedleft := o2.C("nested_set_left")
	o2nestedright := o2.C("nested_set_right")

	joinCond := sqlb.And(
		sqlb.Equal(o1rootid, o2rootid),
		sqlb.Between(o1nestedleft, o2nestedleft, o2nestedright),
	)
	q := sqlb.Select(o1id).Join(o2, joinCond)
	q.Where(sqlb.Equal(o2id, 2))

	b := builder.New()
	qs, qargs := b.StringArgs(q)

	expqs := "SELECT o1.id FROM organizations AS o1 JOIN organizations AS o2 ON (o1.root_organization_id = o2.root_organization_id AND o1.nested_set_left BETWEEN o2.nested_set_left AND o2.nested_set_right) WHERE o2.id = ?"
	expqargs := []interface{}{2}

	assert.Equal(expqs, qs)
	assert.Equal(expqargs, qargs)
}

func TestNestedSetWithAdditionalJoin(t *testing.T) {
	// ref: https://github.com/jaypipes/sqlb/issues/60
	assert := assert.New(t)

	m := testutil.Meta()
	orgs := sqlb.T(m, "organizations")
	orgUsers := sqlb.T(m, "organization_users")

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

	nestedJoinCond := sqlb.And(
		sqlb.Equal(o1rootid, o2rootid),
		sqlb.Between(o1nestedleft, o2nestedleft, o2nestedright),
	)
	ouJoin := sqlb.And(
		sqlb.Equal(o2id, ouOrgId),
		sqlb.Equal(ouUserId, 1),
	)
	q := sqlb.Select(o1id).Join(o2, nestedJoinCond).Join(ou, ouJoin)

	b := builder.New()
	qs, qargs := b.StringArgs(q)

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

	m := testutil.Meta()
	orgs := sqlb.T(m, "organizations")
	orgUsers := sqlb.T(m, "organization_users")

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

	nestedJoinCond := sqlb.And(
		sqlb.Equal(o1rootid, o2rootid),
		sqlb.Between(o1nestedleft, o2nestedleft, o2nestedright),
	)
	ouJoin := sqlb.And(
		sqlb.Equal(o2id, ouOrgId),
		sqlb.Equal(ouUserId, 1),
	)
	subq := sqlb.Select(o1id).Join(o2, nestedJoinCond).Join(ou, ouJoin).As("derived")
	subqOrgId := subq.C("id")

	b := builder.New()
	qs, qargs := b.StringArgs(subq)

	expqs := "SELECT derived.id FROM (SELECT o1.id FROM organizations AS o1 JOIN organizations AS o2 ON (o1.root_organization_id = o2.root_organization_id AND o1.nested_set_left BETWEEN o2.nested_set_left AND o2.nested_set_right) JOIN organization_users AS ou ON (o2.id = ou.organization_id AND ou.user_id = ?)) AS derived"
	expqargs := []interface{}{1}

	assert.Equal(expqs, qs)
	assert.Equal(expqargs, qargs)

	b = builder.New()
	q := sqlb.Select(
		orgs.C("uuid"),
	).OuterJoin(
		subq,
		sqlb.Equal(
			orgs.C("id"),
			subqOrgId,
		),
	).Where(sqlb.IsNotNull(subqOrgId))

	qs, qargs = b.StringArgs(q)

	expqs = "SELECT organizations.uuid FROM organizations LEFT JOIN (SELECT o1.id FROM organizations AS o1 JOIN organizations AS o2 ON (o1.root_organization_id = o2.root_organization_id AND o1.nested_set_left BETWEEN o2.nested_set_left AND o2.nested_set_right) JOIN organization_users AS ou ON (o2.id = ou.organization_id AND ou.user_id = ?)) AS derived ON organizations.id = derived.id WHERE derived.id IS NOT NULL"
	expqargs = []interface{}{1}

	assert.Equal(expqs, qs)
	assert.Equal(expqargs, qargs)
}

func TestModifyingSelectQueryUpdatesBuffer(t *testing.T) {
	assert := assert.New(t)

	b := builder.New()
	m := testutil.Meta()
	users := sqlb.T(m, "users")

	q := sqlb.Select(users)

	qs, qargs := b.StringArgs(q)
	assert.Equal("SELECT users.id, users.name FROM users", qs)
	assert.Empty(qargs)

	b = builder.New()

	// Modify the underlying SELECT and verify string and args changed
	q.Where(sqlb.Equal(users.C("id"), 1))
	qs, qargs = b.StringArgs(q)
	assert.Equal("SELECT users.id, users.name FROM users WHERE users.id = ?", qs)
	assert.Equal([]interface{}{1}, qargs)
}

func TestFormat(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	articles := sqlb.T(m, "articles")
	colUserName := users.C("name")
	colUserId := users.C("id")
	colArticleId := articles.C("id")
	colArticleAuthor := articles.C("author")

	q := sqlb.Select(colArticleId, colUserName.As("author"))
	q.Join(articles, sqlb.Equal(colUserId, colArticleAuthor))
	q.Where(sqlb.Equal(colUserName, "foo"))
	q.GroupBy(colUserName)
	q.OrderBy(colUserName.Desc())
	q.Limit(10)

	tests := []struct {
		name  string
		b     *builder.Builder
		query api.Element
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
			name: "newline clause separator ",
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
		qs, qargs := tt.b.StringArgs(tt.query)
		assert.Equal(tt.qs, qs)
		assert.Equal(tt.qargs, qargs)
	}
}
