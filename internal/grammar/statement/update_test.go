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
	"github.com/jaypipes/sqlb/internal/grammar/clause"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
	"github.com/jaypipes/sqlb/internal/grammar/identifier"
	"github.com/jaypipes/sqlb/internal/grammar/statement"
	"github.com/jaypipes/sqlb/internal/scanner"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/jaypipes/sqlb/types"
	"github.com/stretchr/testify/assert"
)

func TestUpdateStatement(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		s     *statement.Update
		qs    string
		qargs []interface{}
	}{
		{
			name: "UPDATE no WHERE",
			s: statement.NewUpdate(
				users,
				[]*identifier.Column{colUserName},
				[]interface{}{"foo"},
				nil,
			),
			qs:    "UPDATE users SET name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "UPDATE simple WHERE",
			s: statement.NewUpdate(
				users,
				[]*identifier.Column{colUserName},
				[]interface{}{"foo"},
				clause.NewWhere(
					expression.Equal(colUserName, "bar"),
				),
			),
			qs:    "UPDATE users SET name = ? WHERE users.name = ?",
			qargs: []interface{}{"foo", "bar"},
		},
	}
	for _, test := range tests {
		expArgc := len(test.qargs)
		argc := test.s.ArgCount()
		assert.Equal(expArgc, argc)

		expLen := len(test.qs)
		size := test.s.Size(scanner.DefaultScanner)
		size += scanner.InterpolationLength(types.DialectMySQL, argc)
		assert.Equal(expLen, size)

		var b strings.Builder
		b.Grow(size)
		curArg := 0
		test.s.Scan(scanner.DefaultScanner, &b, test.qargs, &curArg)

		assert.Equal(test.qs, b.String())
	}
}
