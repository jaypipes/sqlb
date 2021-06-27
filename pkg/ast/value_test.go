//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ast_test

import (
	"testing"

	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	assert := assert.New(t)

	v := ast.NewValue(nil, "foo")

	s := v.Size(scanner.DefaultScanner)
	// Due to dialect handling, we can't include interpolation markers in the
	// size calculation, so size() always returns 0 for non-aliased values.
	assert.Equal(0, s)

	argc := v.ArgCount()
	assert.Equal(1, argc)

	args := make([]interface{}, 1)
	b := make([]byte, 1)
	curArg := 0
	written := v.Scan(scanner.DefaultScanner, b, args, &curArg)

	exp := "?"
	expLen := len(exp)

	assert.Equal(expLen, written)
	assert.Equal(exp, string(b))
	assert.Equal("foo", args[0])
}
