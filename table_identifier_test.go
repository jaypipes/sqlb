//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"testing"

	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := T(sc, "users")

	exp := "users"
	expLen := len(exp)
	s := users.Size(scanner.DefaultScanner)
	assert.Equal(expLen, s)

	b := make([]byte, s)
	written := users.Scan(scanner.DefaultScanner, b, nil, nil)

	assert.Equal(written, s)
	assert.Equal(exp, string(b))
}

func TestTableAlias(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	u := T(sc, "users").As("u")

	exp := "users AS u"
	expLen := len(exp)
	s := u.Size(scanner.DefaultScanner)
	assert.Equal(expLen, s)

	b := make([]byte, s)
	written := u.Scan(scanner.DefaultScanner, b, nil, nil)

	assert.Equal(written, s)
	assert.Equal(exp, string(b))
}

func TestTableColumns(t *testing.T) {
	assert := assert.New(t)

	td := &TableIdentifier{
		name: "users",
	}

	cols := []*ColumnIdentifier{
		&ColumnIdentifier{
			name: "id",
			tbl:  td,
		},
		&ColumnIdentifier{
			name: "email",
			tbl:  td,
		},
	}
	td.columns = cols

	defs := td.columns

	assert.Equal(2, len(defs))
	for _, def := range defs {
		assert.Equal(td, def.tbl)
	}

	// Check stable order of insertion from above...
	assert.Equal(defs[0].name, "id")
	assert.Equal(defs[1].name, "email")
}

func TestTableC(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := T(sc, "users")

	c := users.C("name")

	assert.Equal(users, c.tbl)
	assert.Equal("name", c.name)

	// Check an unknown column name returns nil
	unknown := users.C("unknown")
	assert.Nil(unknown)
}
