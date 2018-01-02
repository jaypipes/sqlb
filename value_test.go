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

func Testvalue(t *testing.T) {
	assert := assert.New(t)

	v := &value{val: "foo"}

	exp := "?"
	expLen := len(exp)
	expArgCount := 1

	s := v.size(defaultScanner)
	assert.Equal(expLen, s)

	argc := v.argCount()
	assert.Equal(expArgCount, argc)

	args := make([]interface{}, 1)
	b := make([]byte, s)
	curArg := 0
	written := v.scan(defaultScanner, b, args, &curArg)

	assert.Equal(s, written)
	assert.Equal(exp, string(b))
	assert.Equal("foo", args[0])
}
