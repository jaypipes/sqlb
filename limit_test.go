package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestLimitClause(t *testing.T) {
    assert := assert.New(t)

    lc := &LimitClause{
        limit: 20,
    }

    exp := " LIMIT ?"
    expLen := len(exp)
    expargCount := 1

    s := lc.size()
    assert.Equal(expLen, s)

    argc := lc.argCount()
    assert.Equal(expargCount, argc)

    args := make([]interface{}, expargCount)
    b := make([]byte, s)
    written, numArgs := lc.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expargCount, numArgs)
    assert.Equal(20, args[0])
}

func TestLimitClauseWithOffset(t *testing.T) {
    assert := assert.New(t)

    lc := &LimitClause{
        limit: 20,
    }
    offset := 10
    lc.offset = &offset

    exp := " LIMIT ? OFFSET ?"
    expLen := len(exp)
    expargCount := 2

    s := lc.size()
    assert.Equal(expLen, s)

    argc := lc.argCount()
    assert.Equal(expargCount, argc)

    args := make([]interface{}, expargCount)
    b := make([]byte, s)
    written, numArgs := lc.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expargCount, numArgs)
    assert.Equal(20, args[0])
    assert.Equal(10, args[1])
}
