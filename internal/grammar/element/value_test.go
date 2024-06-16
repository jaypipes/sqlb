//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package element_test

import (
	"testing"

	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar/element"
	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	assert := assert.New(t)

	v := element.NewValue(nil, "foo")

	b := builder.New()
	s := v.Size(b)
	// Due to dialect handling, we can't include interpolation markers in the
	// size calculation, so size() always returns 0 for non-aliased values.
	assert.Equal(0, s)

	argc := v.ArgCount()
	assert.Equal(1, argc)

	args := make([]interface{}, 1)
	curArg := 0
	v.Scan(b, args, &curArg)

	exp := "?"

	assert.Equal(exp, b.String())
	assert.Equal("foo", args[0])
}
