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
    expArgCount := 1

    s := val.Size()
    assert.Equal(expLen, s)

    argc := val.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, 1)
    b := make([]byte, s)
    written, numArgs  := val.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal("foo", args[0])
    assert.Equal(expArgCount, numArgs)
}
