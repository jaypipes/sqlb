//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause_test

import (
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar/clause"
	"github.com/jaypipes/sqlb/internal/grammar/statement"
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
		{
			c: clause.NewDerivedTable(
				"u",
				statement.NewSelect(
					[]api.Projection{
						colUserName,
					},
					[]api.Selection{
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
		b := builder.New()

		expArgc := len(test.qargs)
		assert.Equal(expArgc, test.c.ArgCount())

		qs, _ := b.StringArgs(test.c)

		assert.Equal(test.qs, qs)
	}
}
