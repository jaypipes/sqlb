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
	"github.com/jaypipes/sqlb/internal/grammar/function"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

type orderByTest struct {
	c     *clause.OrderBy
	qs    string
	qargs []interface{}
}

func TestOrderBy(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	tests := []orderByTest{
		// column asc
		{
			c:  clause.NewOrderBy(colUserName.Asc()),
			qs: " ORDER BY users.name",
		},
		// column desc
		{
			c:  clause.NewOrderBy(colUserName.Desc()),
			qs: " ORDER BY users.name DESC",
		},
		// Aliased column should NOT output alias in ORDER BY
		{
			c:  clause.NewOrderBy(colUserName.As("user_name").Desc()),
			qs: " ORDER BY users.name DESC",
		},
		// multi column mixed
		{
			c:  clause.NewOrderBy(colUserName.Asc(), colUserId.Desc()),
			qs: " ORDER BY users.name, users.id DESC",
		},
		// sort by a function
		{
			c:  clause.NewOrderBy(function.Count(users).Desc()),
			qs: " ORDER BY COUNT(*) DESC",
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
