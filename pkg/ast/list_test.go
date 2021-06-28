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

func TestListSingle(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := sqlb.T(sc, "users")
	colUserName := users.C("name")

	cl := ast.NewList(colUserName)

	exp := "users.name"
	expLen := len(exp)
	s := cl.Size(scanner.DefaultScanner)
	assert.Equal(expLen, s)

	b := make([]byte, s)
	written := cl.Scan(scanner.DefaultScanner, b, nil, nil)

	assert.Equal(written, s)
	assert.Equal(exp, string(b))
}

func TestListMulti(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := sqlb.T(sc, "users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	cl := ast.NewList(colUserId, colUserName)

	exp := "users.id, users.name"
	expLen := len(exp)
	s := cl.Size(scanner.DefaultScanner)
	assert.Equal(expLen, s)

	b := make([]byte, s)
	written := cl.Scan(scanner.DefaultScanner, b, nil, nil)

	assert.Equal(written, s)
	assert.Equal(exp, string(b))
}
