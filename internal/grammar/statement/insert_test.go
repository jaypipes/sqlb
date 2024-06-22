//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package statement_test

import (
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar/identifier"
	"github.com/jaypipes/sqlb/internal/grammar/statement"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestInsertStatement(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	colUserName := users.C("name")
	colUserId := users.C("id")

	tests := []struct {
		name  string
		s     *statement.Insert
		qs    string
		qargs []interface{}
	}{
		{
			name: "Simple INSERT",
			s: statement.NewInsert(
				users,
				[]*identifier.Column{colUserId, colUserName},
				[]interface{}{nil, "foo"},
			),
			qs:    "INSERT INTO users (id, name) VALUES (?, ?)",
			qargs: []interface{}{nil, "foo"},
		},
		{
			name: "Ensure no aliasing in table names",
			s: statement.NewInsert(
				users.As("u"),
				[]*identifier.Column{colUserId, colUserName},
				[]interface{}{nil, "foo"},
			),
			qs:    "INSERT INTO users (id, name) VALUES (?, ?)",
			qargs: []interface{}{nil, "foo"},
		},
		{
			name: "Ensure no aliasing in column names",
			s: statement.NewInsert(
				users,
				[]*identifier.Column{
					colUserId.As("user_id").(*identifier.Column),
					colUserName,
				},
				[]interface{}{nil, "foo"},
			),
			qs:    "INSERT INTO users (id, name) VALUES (?, ?)",
			qargs: []interface{}{nil, "foo"},
		},
	}
	for _, test := range tests {
		b := builder.New()

		expArgc := len(test.qargs)
		argc := test.s.ArgCount()
		assert.Equal(expArgc, argc)

		qs, args := b.StringArgs(test.s)
		assert.Equal(test.qs, qs)
		if len(test.qargs) > 0 {
			assert.Equal(test.qargs, args)
		}
	}
}
