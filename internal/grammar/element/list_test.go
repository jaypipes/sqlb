//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package element_test

import (
	"strings"
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/internal/grammar/element"
	"github.com/jaypipes/sqlb/internal/scanner"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestListSingle(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	colUserName := users.C("name")

	cl := element.NewList(colUserName)

	exp := "users.name"
	expLen := len(exp)
	s := cl.Size(scanner.DefaultScanner)
	assert.Equal(expLen, s)

	var b strings.Builder
	b.Grow(s)
	cl.Scan(scanner.DefaultScanner, &b, nil, nil)

	assert.Equal(exp, b.String())
}

func TestListMulti(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	cl := element.NewList(colUserId, colUserName)

	exp := "users.id, users.name"
	expLen := len(exp)
	s := cl.Size(scanner.DefaultScanner)
	assert.Equal(expLen, s)

	var b strings.Builder
	b.Grow(s)
	cl.Scan(scanner.DefaultScanner, &b, nil, nil)

	assert.Equal(exp, b.String())
}
