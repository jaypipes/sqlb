//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package element_test

import (
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar/element"
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

	b := builder.New()

	qs, _ := b.StringArgs(cl)

	assert.Equal(exp, qs)
}

func TestListMulti(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	cl := element.NewList(colUserId, colUserName)

	exp := "users.id, users.name"

	b := builder.New()

	qs, _ := b.StringArgs(cl)

	assert.Equal(exp, qs)
}
