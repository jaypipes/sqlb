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

	s := v.size()
	assert.Equal(expLen, s)

	argc := v.argCount()
	assert.Equal(expArgCount, argc)

	args := make([]interface{}, 1)
	b := make([]byte, s)
	written, numArgs := v.scan(b, args)

	assert.Equal(s, written)
	assert.Equal(exp, string(b))
	assert.Equal("foo", args[0])
	assert.Equal(expArgCount, numArgs)
}
