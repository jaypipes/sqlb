//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"testing"

	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestListSingle(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	colUserName := users.C("name")

	cl := &List{elements: []types.Element{colUserName}}

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

	m := testFixtureMeta()
	users := m.Table("users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	cl := &List{elements: []types.Element{colUserId, colUserName}}

	exp := "users.id, users.name"
	expLen := len(exp)
	s := cl.Size(scanner.DefaultScanner)
	assert.Equal(expLen, s)

	b := make([]byte, s)
	written := cl.Scan(scanner.DefaultScanner, b, nil, nil)

	assert.Equal(written, s)
	assert.Equal(exp, string(b))
}
