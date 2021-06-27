//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"testing"

	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestLimitClause(t *testing.T) {
	assert := assert.New(t)

	lc := &LimitClause{
		limit: 20,
	}

	exp := " LIMIT ?"
	expArgCount := 1

	argc := lc.ArgCount()
	assert.Equal(expArgCount, argc)

	size := lc.Size(scanner.DefaultScanner)
	size += scanner.InterpolationLength(types.DIALECT_MYSQL, argc)
	expLen := len(exp)
	assert.Equal(expLen, size)

	args := make([]interface{}, expArgCount)
	b := make([]byte, size)
	curArg := 0
	written := lc.Scan(scanner.DefaultScanner, b, args, &curArg)

	assert.Equal(size, written)
	assert.Equal(exp, string(b))
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
	expArgCount := 2

	argc := lc.ArgCount()
	assert.Equal(expArgCount, argc)

	size := lc.Size(scanner.DefaultScanner)
	size += scanner.InterpolationLength(types.DIALECT_MYSQL, argc)
	expLen := len(exp)
	assert.Equal(expLen, size)

	args := make([]interface{}, expArgCount)
	b := make([]byte, size)
	curArg := 0
	written := lc.Scan(scanner.DefaultScanner, b, args, &curArg)

	assert.Equal(size, written)
	assert.Equal(exp, string(b))
	assert.Equal(20, args[0])
	assert.Equal(10, args[1])
}
