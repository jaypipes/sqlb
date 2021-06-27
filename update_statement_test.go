//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"testing"

	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/jaypipes/sqlb/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestUpdateStatement(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := T(sc, "users")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		s     *UpdateStatement
		qs    string
		qargs []interface{}
	}{
		{
			name: "UPDATE no WHERE",
			s: &UpdateStatement{
				table:   users,
				columns: []*ast.ColumnIdentifier{colUserName},
				values:  []interface{}{"foo"},
			},
			qs:    "UPDATE users SET name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "UPDATE simple WHERE",
			s: &UpdateStatement{
				table:   users,
				columns: []*ast.ColumnIdentifier{colUserName},
				values:  []interface{}{"foo"},
				where: ast.NewWhereClause(
					ast.Equal(colUserName, "bar"),
				),
			},
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
		size += scanner.InterpolationLength(types.DIALECT_MYSQL, argc)
		assert.Equal(expLen, size)

		b := make([]byte, size)
		curArg := 0
		written := test.s.Scan(scanner.DefaultScanner, b, test.qargs, &curArg)

		assert.Equal(written, size)
		assert.Equal(test.qs, string(b))
	}
}
