//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause_test

import (
	"strings"
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/internal/grammar/clause"
	"github.com/jaypipes/sqlb/internal/grammar/statement"
	"github.com/jaypipes/sqlb/internal/scanner"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

type derivedTest struct {
	c     *clause.DerivedTable
	qs    string
	qargs []interface{}
}

func TestDerived(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	colUserName := users.C("name")

	tests := []derivedTest{
		// Simple one-column sub-SELECT
		derivedTest{
			c: clause.NewDerivedTable(
				"u",
				statement.NewSelect(
					[]scanner.Projection{
						colUserName,
					},
					[]scanner.Selection{
						users,
					},
					nil,
					nil,
					nil,
					nil,
					nil,
					nil,
				),
			),
			qs: "(SELECT users.name FROM users) AS u",
		},
	}
	for _, test := range tests {
		expLen := len(test.qs)
		s := test.c.Size(scanner.DefaultScanner)
		assert.Equal(expLen, s)

		expArgc := len(test.qargs)
		assert.Equal(expArgc, test.c.ArgCount())

		var b strings.Builder
		b.Grow(s)
		curArg := 0
		test.c.Scan(scanner.DefaultScanner, &b, test.qargs, &curArg)

		assert.Equal(test.qs, b.String())
	}
}
