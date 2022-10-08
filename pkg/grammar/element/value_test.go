//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package element_test

import (
	"strings"
	"testing"

	"github.com/jaypipes/sqlb/pkg/grammar/element"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	assert := assert.New(t)

	v := element.NewValue(nil, "foo")

	s := v.Size(scanner.DefaultScanner)
	// Due to dialect handling, we can't include interpolation markers in the
	// size calculation, so size() always returns 0 for non-aliased values.
	assert.Equal(0, s)

	argc := v.ArgCount()
	assert.Equal(1, argc)

	args := make([]interface{}, 1)
	var b strings.Builder
	b.Grow(s)
	curArg := 0
	v.Scan(scanner.DefaultScanner, &b, args, &curArg)

	exp := "?"

	assert.Equal(exp, b.String())
	assert.Equal("foo", args[0])
}
