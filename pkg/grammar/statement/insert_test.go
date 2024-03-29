//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package statement_test

import (
	"strings"
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/pkg/grammar/identifier"
	"github.com/jaypipes/sqlb/pkg/grammar/statement"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/jaypipes/sqlb/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestInsertStatement(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := sqlb.T(sc, "users")
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
		expArgc := len(test.qargs)
		argc := test.s.ArgCount()
		assert.Equal(expArgc, argc)

		expLen := len(test.qs)
		size := test.s.Size(scanner.DefaultScanner)
		size += scanner.InterpolationLength(types.DIALECT_MYSQL, argc)
		assert.Equal(expLen, size)

		var b strings.Builder
		b.Grow(size)
		curArg := 0
		test.s.Scan(scanner.DefaultScanner, &b, test.qargs, &curArg)

		assert.Equal(test.qs, b.String())
	}
}
