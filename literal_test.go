package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestLiteral(t *testing.T) {
    assert := assert.New(t)

    lit := &Literal{value: "foo"}

    exp := "?"
    expLen := len(exp)
    expArgCount := 1

    s := lit.Size()
    assert.Equal(expLen, s)

    argc := lit.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, 1)
    b := make([]byte, s)
    written, numArgs  := lit.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal("foo", args[0])
    assert.Equal(expArgCount, numArgs)
}
