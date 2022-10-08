//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package identifier_test

import (
	"strings"
	"testing"

	"github.com/jaypipes/sqlb/pkg/grammar/identifier"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := identifier.TableFromSchema(sc, "users")

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

	sc := testutil.Schema()
	u := identifier.TableFromSchema(sc, "users").As("u")

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

	sc := testutil.Schema()
	users := identifier.TableFromSchema(sc, "users").As("u")

	assert.Equal(2, len(users.Projections()))
}

func TestTableC(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := identifier.TableFromSchema(sc, "users")

	c := users.C("name")

	assert.Equal(users, c.From())
	assert.Equal("name", c.Name)

	// Check an unknown column name returns nil
	unknown := users.C("unknown")
	assert.Nil(unknown)
}
