//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause_test

import (
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar/clause"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

type GroupByTest struct {
	c     *clause.GroupBy
	qs    string
	qargs []interface{}
}

func TestGroupBy(t *testing.T) {
	assert := assert.New(t)

	m := testutil.M()
	users := sqlb.T(m, "users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	tests := []GroupByTest{
		// Single column
		{
			c:  clause.NewGroupBy(colUserName),
			qs: " GROUP BY users.name",
		},
		// Multiple columns
		{
			c:  clause.NewGroupBy(colUserName, colUserId),
			qs: " GROUP BY users.name, users.id",
		},
		// Aliased column should NOT output alias in GROUP BY
		{
			c:  clause.NewGroupBy(colUserName.As("user_name")),
			qs: " GROUP BY users.name",
		},
	}
	for _, test := range tests {
		b := builder.New()

		expArgc := len(test.qargs)
		argc := test.c.ArgCount()
		assert.Equal(expArgc, argc)

		qs, _ := b.StringArgs(test.c)

		assert.Equal(test.qs, qs)
	}
}
