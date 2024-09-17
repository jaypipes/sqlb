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

func TestSelect(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserName := users.C("name")
	colUserId := users.C("id")

	tests := []struct {
		name  string
		q     *expr.Selection
		qs    string
		qargs []interface{}
	}{
		{
			name: "Simple FROM",
			q:    expr.Select(users),
			qs:   "SELECT users.created_on, users.id, users.name FROM users",
		},
		{
			name:  "Simple WHERE",
			q:     expr.Select(colUserId, colUserName).Where(expr.Equal(colUserName, "foo")),
			qs:    "SELECT users.id, users.name FROM users WHERE users.name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "WHERE with an OR sqlb",
			q: expr.Select(colUserId, colUserName).Where(
				expr.Or(
					expr.Equal(colUserName, "foo"),
					expr.Equal(colUserName, "bar"),
				),
			),
			qs:    "SELECT users.id, users.name FROM users WHERE (users.name = ? OR users.name = ?)",
			qargs: []interface{}{"foo", "bar"},
		},
		{
			name: "Simple GROUP BY",
			q:    expr.Select(colUserId, colUserName).GroupBy(colUserName),
			qs:   "SELECT users.id, users.name FROM users GROUP BY users.name",
		},
		{
			name:  "Simple HAVING",
			q:     expr.Select(colUserId, colUserName).Having(expr.Equal(colUserName, "foo")),
			qs:    "SELECT users.id, users.name FROM users HAVING users.name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "Simple ORDER BY",
			q:    expr.Select(colUserId, colUserName).OrderBy(colUserName.Desc()),
			qs:   "SELECT users.id, users.name FROM users ORDER BY users.name DESC",
		},
		{
			name:  "Simple LIMIT",
			q:     expr.Select(colUserId, colUserName).Limit(10),
			qs:    "SELECT users.id, users.name FROM users LIMIT ?",
			qargs: []interface{}{10},
		},
		{
			name:  "Simple LIMIT with OFFSET",
			q:     expr.Select(colUserId, colUserName).LimitWithOffset(10, 20),
			qs:    "SELECT users.id, users.name FROM users LIMIT ? OFFSET ?",
			qargs: []interface{}{10, 20},
		},
		{
			name: "Simple named derived table",
			q:    expr.Select(expr.Select(colUserId, colUserName).As("u")),
			qs:   "SELECT u.id, u.name FROM (SELECT users.id, users.name FROM users) AS u",
		},
		{
			name: "Simple un-named derived table",
			q:    expr.Select(expr.Select(colUserId, colUserName)),
			qs:   "SELECT derived0.id, derived0.name FROM (SELECT users.id, users.name FROM users) AS derived0",
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

func TestModifyingSelectQueryUpdatesBuffer(t *testing.T) {
	assert := assert.New(t)

	b := builder.New()
	m := testutil.M()
	users := m.T("users")
	colUserName := users.C("name")
	colUserId := users.C("id")

	q := expr.Select(colUserId, colUserName)

	qs, qargs := b.StringArgs(q.Query())
	assert.Equal("SELECT users.id, users.name FROM users", qs)
	assert.Empty(qargs)

	b = builder.New()

	// Modify the underlying SELECT and verify string and args changed
	q.Where(expr.Equal(users.C("id"), 1))
	qs, qargs = b.StringArgs(q.Query())
	assert.Equal("SELECT users.id, users.name FROM users WHERE users.id = ?", qs)
	assert.Equal([]interface{}{1}, qargs)
}
