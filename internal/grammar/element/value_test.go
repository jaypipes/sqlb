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

	argc := v.ArgCount()
	assert.Equal(1, argc)

	exp := "?"

	b := builder.New()

	qs, args := b.StringArgs(v)

	assert.Equal(exp, qs)
	assert.Equal("foo", args[0])
}
