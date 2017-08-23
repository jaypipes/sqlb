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
    expArgCount := 1

    s := lc.size()
    assert.Equal(expLen, s)

    argc := lc.argCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := lc.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
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
    expArgCount := 2

    s := lc.size()
    assert.Equal(expLen, s)

    argc := lc.argCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := lc.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
    assert.Equal(20, args[0])
    assert.Equal(10, args[1])
}
