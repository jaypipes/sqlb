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

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	tests := []GroupByTest{
		// Single column
		GroupByTest{
			c:  clause.NewGroupBy(colUserName),
			qs: " GROUP BY users.name",
		},
		// Multiple columns
		GroupByTest{
			c:  clause.NewGroupBy(colUserName, colUserId),
			qs: " GROUP BY users.name, users.id",
		},
		// Aliased column should NOT output alias in GROUP BY
		GroupByTest{
			c:  clause.NewGroupBy(colUserName.As("user_name")),
			qs: " GROUP BY users.name",
		},
	}
	for _, test := range tests {
		b := builder.New()

		expLen := len(test.qs)
		s := test.c.Size(b)
		assert.Equal(expLen, s)

		expArgc := len(test.qargs)
		assert.Equal(expArgc, test.c.ArgCount())

		curArg := 0
		test.c.Scan(b, test.qargs, &curArg)

		assert.Equal(test.qs, b.String())
	}
}
