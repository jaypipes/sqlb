package sqlb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitClause(t *testing.T) {
	assert := assert.New(t)

	lc := &limitClause{
		limit: 20,
	}

	exp := " LIMIT ?"
	expArgCount := 1

	argc := lc.argCount()
	assert.Equal(expArgCount, argc)

	size := lc.size(defaultScanner)
	size += interpolationLength(DIALECT_MYSQL, argc)
	expLen := len(exp)
	assert.Equal(expLen, size)

	args := make([]interface{}, expArgCount)
	b := make([]byte, size)
	curArg := 0
	written := lc.scan(defaultScanner, b, args, &curArg)

	assert.Equal(size, written)
	assert.Equal(exp, string(b))
	assert.Equal(20, args[0])
}

func TestLimitClauseWithOffset(t *testing.T) {
	assert := assert.New(t)

	lc := &limitClause{
		limit: 20,
	}
	offset := 10
	lc.offset = &offset

	exp := " LIMIT ? OFFSET ?"
	expArgCount := 2

	argc := lc.argCount()
	assert.Equal(expArgCount, argc)

	size := lc.size(defaultScanner)
	size += interpolationLength(DIALECT_MYSQL, argc)
	expLen := len(exp)
	assert.Equal(expLen, size)

	args := make([]interface{}, expArgCount)
	b := make([]byte, size)
	curArg := 0
	written := lc.scan(defaultScanner, b, args, &curArg)

	assert.Equal(size, written)
	assert.Equal(exp, string(b))
	assert.Equal(20, args[0])
	assert.Equal(10, args[1])
}
