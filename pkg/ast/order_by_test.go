//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ast_test

import (
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

type orderByTest struct {
	c     *ast.OrderByClause
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
			c:  ast.NewOrderByClause(colUserName.Asc()),
			qs: " ORDER BY users.name",
		},
		// column desc
		orderByTest{
			c:  ast.NewOrderByClause(colUserName.Desc()),
			qs: " ORDER BY users.name DESC",
		},
		// Aliased column should NOT output alias in ORDER BY
		orderByTest{
			c:  ast.NewOrderByClause(colUserName.As("user_name").Desc()),
			qs: " ORDER BY users.name DESC",
		},
		// multi column mixed
		orderByTest{
			c:  ast.NewOrderByClause(colUserName.Asc(), colUserId.Desc()),
			qs: " ORDER BY users.name, users.id DESC",
		},
		// sort by a function
		orderByTest{
			c:  ast.NewOrderByClause(ast.Count(users).Desc()),
			qs: " ORDER BY COUNT(*) DESC",
		},
	}
	for _, test := range tests {
		expLen := len(test.qs)
		s := test.c.Size(scanner.DefaultScanner)
		assert.Equal(expLen, s)

		expArgc := len(test.qargs)
		assert.Equal(expArgc, test.c.ArgCount())

		b := make([]byte, s)
		curArg := 0
		written := test.c.Scan(scanner.DefaultScanner, b, test.qargs, &curArg)

		assert.Equal(written, s)
		assert.Equal(test.qs, string(b))
	}
}
