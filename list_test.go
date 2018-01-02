//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListSingle(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	colUserName := users.C("name")

	cl := &List{elements: []element{colUserName}}

	exp := "users.name"
	expLen := len(exp)
	s := cl.size(defaultScanner)
	assert.Equal(expLen, s)

	b := make([]byte, s)
	written := cl.scan(defaultScanner, b, nil, nil)

	assert.Equal(written, s)
	assert.Equal(exp, string(b))
}

func TestListMulti(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	cl := &List{elements: []element{colUserId, colUserName}}

	exp := "users.id, users.name"
	expLen := len(exp)
	s := cl.size(defaultScanner)
	assert.Equal(expLen, s)

	b := make([]byte, s)
	written := cl.scan(defaultScanner, b, nil, nil)

	assert.Equal(written, s)
	assert.Equal(exp, string(b))
}
