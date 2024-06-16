//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package identifier_test

import (
	"strings"
	"testing"

	"github.com/jaypipes/sqlb/internal/grammar/identifier"
	"github.com/jaypipes/sqlb/internal/scanner"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := identifier.TableFromMeta(m, "users")

	exp := "users"
	expLen := len(exp)
	s := users.Size(scanner.DefaultScanner)
	assert.Equal(expLen, s)

	var b strings.Builder
	b.Grow(s)
	users.Scan(scanner.DefaultScanner, &b, nil, nil)

	assert.Equal(exp, b.String())
}

func TestTableAlias(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	u := identifier.TableFromMeta(m, "users").As("u")

	exp := "users AS u"
	expLen := len(exp)
	s := u.Size(scanner.DefaultScanner)
	assert.Equal(expLen, s)

	var b strings.Builder
	b.Grow(s)
	u.Scan(scanner.DefaultScanner, &b, nil, nil)

	assert.Equal(exp, b.String())
}

func TestTableColumns(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := identifier.TableFromMeta(m, "users").As("u")

	assert.Equal(2, len(users.Projections()))
}

func TestTableC(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := identifier.TableFromMeta(m, "users")

	c := users.C("name")

	assert.Equal(users, c.From())
	assert.Equal("name", c.Name)

	// Check an unknown column name returns nil
	unknown := users.C("unknown")
	assert.Nil(unknown)
}
