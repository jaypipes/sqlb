//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package identifier_test

import (
	"strings"
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestC(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := sqlb.T(sc, "users")
	c := users.C("name")

	exp := "users.name"
	expLen := len(exp)
	s := c.Size(scanner.DefaultScanner)
	assert.Equal(expLen, s)

	var b strings.Builder
	b.Grow(s)
	c.Scan(scanner.DefaultScanner, &b, nil, nil)

	assert.Equal(exp, b.String())
}

func TestColumnWithTableAlias(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := sqlb.T(sc, "users").As("u")
	c := users.C("name")

	exp := "u.name"
	expLen := len(exp)
	s := c.Size(scanner.DefaultScanner)
	assert.Equal(expLen, s)

	var b strings.Builder
	b.Grow(s)
	c.Scan(scanner.DefaultScanner, &b, nil, nil)

	assert.Equal(exp, b.String())
}

func TestColumnAlias(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := sqlb.T(sc, "users")
	c := users.C("name").As("user_name")

	exp := "users.name AS user_name"
	expLen := len(exp)
	s := c.Size(scanner.DefaultScanner)
	assert.Equal(expLen, s)

	var b strings.Builder
	b.Grow(s)
	c.Scan(scanner.DefaultScanner, &b, nil, nil)

	assert.Equal(exp, b.String())
}
