//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause_test

import (
	"strings"
	"testing"

	"github.com/jaypipes/sqlb/internal/grammar/clause"
	"github.com/jaypipes/sqlb/internal/scanner"
	"github.com/jaypipes/sqlb/types"
	"github.com/stretchr/testify/assert"
)

func TestLimitClause(t *testing.T) {
	assert := assert.New(t)

	lc := clause.NewLimit(20, nil)

	exp := " LIMIT ?"
	expArgCount := 1

	argc := lc.ArgCount()
	assert.Equal(expArgCount, argc)

	size := lc.Size(scanner.DefaultScanner)
	size += scanner.InterpolationLength(types.DialectMySQL, argc)
	expLen := len(exp)
	assert.Equal(expLen, size)

	args := make([]interface{}, expArgCount)
	var b strings.Builder
	b.Grow(size)
	curArg := 0
	lc.Scan(scanner.DefaultScanner, &b, args, &curArg)

	assert.Equal(exp, b.String())
	assert.Equal(20, args[0])
}

func TestLimitClauseWithOffset(t *testing.T) {
	assert := assert.New(t)

	offset := 10
	lc := clause.NewLimit(20, &offset)

	exp := " LIMIT ? OFFSET ?"
	expArgCount := 2

	argc := lc.ArgCount()
	assert.Equal(expArgCount, argc)

	size := lc.Size(scanner.DefaultScanner)
	size += scanner.InterpolationLength(types.DialectMySQL, argc)
	expLen := len(exp)
	assert.Equal(expLen, size)

	args := make([]interface{}, expArgCount)
	var b strings.Builder
	b.Grow(size)
	curArg := 0
	lc.Scan(scanner.DefaultScanner, &b, args, &curArg)

	assert.Equal(exp, b.String())
	assert.Equal(20, args[0])
	assert.Equal(10, args[1])
}
