package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
    assert := assert.New(t)

    val := &Value{value: "foo"}

    exp := "?"
    expLen := len(exp)
    expargCount := 1

    s := val.size()
    assert.Equal(expLen, s)

    argc := val.argCount()
    assert.Equal(expargCount, argc)

    args := make([]interface{}, 1)
    b := make([]byte, s)
    written, numArgs  := val.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal("foo", args[0])
    assert.Equal(expargCount, numArgs)
}
