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

func TestValue(t *testing.T) {
	assert := assert.New(t)

	v := &value{val: "foo"}

	s := v.size(defaultScanner)
	// Due to dialect handling, we can't include interpolation markers in the
	// size calculation, so size() always returns 0 for non-aliased values.
	assert.Equal(0, s)

	argc := v.argCount()
	assert.Equal(1, argc)

	args := make([]interface{}, 1)
	b := make([]byte, 1)
	curArg := 0
	written := v.scan(defaultScanner, b, args, &curArg)

	exp := "?"
	expLen := len(exp)

	assert.Equal(expLen, written)
	assert.Equal(exp, string(b))
	assert.Equal("foo", args[0])
}
