//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package identifier_test

import (
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestC(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	c := users.C("name")

	b := builder.New()

	exp := "users.name"
	expLen := len(exp)
	s := c.Size(b)
	assert.Equal(expLen, s)

	c.Scan(b, nil, nil)

	assert.Equal(exp, b.String())
}

func TestColumnWithTableAlias(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users").As("u")
	c := users.C("name")

	b := builder.New()

	exp := "u.name"
	expLen := len(exp)
	s := c.Size(b)
	assert.Equal(expLen, s)

	c.Scan(b, nil, nil)

	assert.Equal(exp, b.String())
}

func TestColumnAlias(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	c := users.C("name").As("user_name")

	b := builder.New()

	exp := "users.name AS user_name"
	expLen := len(exp)
	s := c.Size(b)
	assert.Equal(expLen, s)

	c.Scan(b, nil, nil)

	assert.Equal(exp, b.String())
}
