//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause_test

import (
	"testing"

	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar/clause"
	"github.com/stretchr/testify/assert"
)

func TestLimitClause(t *testing.T) {
	assert := assert.New(t)

	lc := clause.NewLimit(20, nil)

	exp := " LIMIT ?"
	expArgCount := 1

	argc := lc.ArgCount()
	assert.Equal(expArgCount, argc)

	b := builder.New()

	qs, args := b.StringArgs(lc)
	assert.Equal(exp, qs)
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

	b := builder.New()

	qs, args := b.StringArgs(lc)
	assert.Equal(exp, qs)
	assert.Equal(20, args[0])
	assert.Equal(10, args[1])
}
