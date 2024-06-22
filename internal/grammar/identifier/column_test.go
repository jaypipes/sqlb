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

	qs, _ := b.StringArgs(c)

	assert.Equal(exp, qs)
}

func TestColumnWithTableAlias(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users").As("u")
	c := users.C("name")

	b := builder.New()

	exp := "u.name"

	qs, _ := b.StringArgs(c)

	assert.Equal(exp, qs)
}

func TestColumnAlias(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	c := users.C("name").As("user_name")

	b := builder.New()

	exp := "users.name AS user_name"

	qs, _ := b.StringArgs(c)

	assert.Equal(exp, qs)
}
