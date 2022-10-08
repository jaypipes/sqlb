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
	"github.com/jaypipes/sqlb/pkg/grammar/clause"
	"github.com/jaypipes/sqlb/pkg/grammar/function"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

type orderByTest struct {
	c     *clause.OrderBy
	qs    string
	qargs []interface{}
}

func TestOrderBy(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := sqlb.T(sc, "users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	tests := []orderByTest{
		// column asc
		orderByTest{
			c:  clause.NewOrderBy(colUserName.Asc()),
			qs: " ORDER BY users.name",
		},
		// column desc
		orderByTest{
			c:  clause.NewOrderBy(colUserName.Desc()),
			qs: " ORDER BY users.name DESC",
		},
		// Aliased column should NOT output alias in ORDER BY
		orderByTest{
			c:  clause.NewOrderBy(colUserName.As("user_name").Desc()),
			qs: " ORDER BY users.name DESC",
		},
		// multi column mixed
		orderByTest{
			c:  clause.NewOrderBy(colUserName.Asc(), colUserId.Desc()),
			qs: " ORDER BY users.name, users.id DESC",
		},
		// sort by a function
		orderByTest{
			c:  clause.NewOrderBy(function.Count(users).Desc()),
			qs: " ORDER BY COUNT(*) DESC",
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
