//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package identifier_test

import (
	"testing"

	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar/identifier"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	ut := m.Table("users")
	users := identifier.TableFromMeta(ut, "users")

	exp := "users"

	b := builder.New()

	qs, _ := b.StringArgs(users)

	assert.Equal(exp, qs)
}

func TestTableAlias(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	ut := m.Table("users")
	u := identifier.TableFromMeta(ut, "users").As("u")

	exp := "users AS u"

	b := builder.New()

	qs, _ := b.StringArgs(u)

	assert.Equal(exp, qs)
}

func TestTableColumns(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	ut := m.Table("users")
	users := identifier.TableFromMeta(ut, "users").As("u")

	assert.Equal(2, len(users.Projections()))
}

func TestTableC(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	ut := m.Table("users")
	users := identifier.TableFromMeta(ut, "users")

	c := users.C("name")

	assert.Equal(users, c.From())
	assert.Equal("name", c.Name)

	// Check an unknown column name returns nil
	unknown := users.C("unknown")
	assert.Nil(unknown)
}
