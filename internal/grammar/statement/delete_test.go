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
	"github.com/jaypipes/sqlb/internal/grammar/statement"
	"github.com/jaypipes/sqlb/internal/scanner"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/jaypipes/sqlb/types"
	"github.com/stretchr/testify/assert"
)

func TestDeleteStatement(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		s     *statement.Delete
		qs    string
		qargs []interface{}
	}{
		{
			name: "DELETE no WHERE",
			s:    statement.NewDelete(users, nil),
			qs:   "DELETE FROM users",
		},
		{
			name: "DELETE simple WHERE",
			s: statement.NewDelete(
				users,
				clause.NewWhere(
					expression.Equal(colUserName, "foo"),
				),
			),
			qs:    "DELETE FROM users WHERE users.name = ?",
			qargs: []interface{}{"foo"},
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
