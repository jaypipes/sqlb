//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package statement_test

import (
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/grammar/statement"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/jaypipes/sqlb/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestInsertStatement(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := sqlb.T(sc, "users")
	colUserName := users.C("name")
	colUserId := users.C("id")

	tests := []struct {
		name  string
		s     *statement.Insert
		qs    string
		qargs []interface{}
	}{
		{
			name: "Simple INSERT",
			s: statement.NewInsert(
				users,
				[]*ast.ColumnIdentifier{colUserId, colUserName},
				[]interface{}{nil, "foo"},
			),
			qs:    "INSERT INTO users (id, name) VALUES (?, ?)",
			qargs: []interface{}{nil, "foo"},
		},
		{
			name: "Ensure no aliasing in table names",
			s: statement.NewInsert(
				users.As("u"),
				[]*ast.ColumnIdentifier{colUserId, colUserName},
				[]interface{}{nil, "foo"},
			),
			qs:    "INSERT INTO users (id, name) VALUES (?, ?)",
			qargs: []interface{}{nil, "foo"},
		},
		{
			name: "Ensure no aliasing in column names",
			s: statement.NewInsert(
				users,
				[]*ast.ColumnIdentifier{
					colUserId.As("user_id").(*ast.ColumnIdentifier),
					colUserName,
				},
				[]interface{}{nil, "foo"},
			),
			qs:    "INSERT INTO users (id, name) VALUES (?, ?)",
			qargs: []interface{}{nil, "foo"},
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
