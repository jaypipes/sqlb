//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package element

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/grammar/sortcolumn"
	pkgscanner "github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/types"
)

// Value is a concrete struct wrapper around a constant that implements the
// Scannable interface. Typically, users won't directly construct value
// structs but instead helper functions like sqlb.Equal() will construct a
// value and bind it to the containing element.
type Value struct {
	sel   types.Selection
	alias string
	val   interface{}
}

func (v *Value) From() types.Selection {
	return v.sel
}

func (v *Value) As(alias string) types.Projection {
	return &Value{
		alias: alias,
		val:   v.val,
	}
}

func (v *Value) DisableAliasScan() func() {
	origAlias := v.alias
	v.alias = ""
	return func() { v.alias = origAlias }
}

func (v *Value) ArgCount() int {
	return 1
}

func (v *Value) Size(scanner types.Scanner) int {
	// Due to dialect handling, we do not include the length of interpolation
	// markers for query parameters. This is calculated separately by the
	// top-level scanning struct before malloc'ing the buffer to inject the SQL
	// string into.
	size := 0
	if v.alias != "" {
		size += len(grammar.Symbols[grammar.SYM_AS]) + len(v.alias)
	}
	return size
}

func (v *Value) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	args[*curArg] = v.val
	bw := pkgscanner.ScanInterpolationMarker(scanner.Dialect(), b, *curArg)
	*curArg++
	if v.alias != "" {
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_AS])
		bw += copy(b[bw:], v.alias)
	}
	return bw
}

func (v *Value) Desc() types.Sortable {
	return sortcolumn.NewDesc(v)
}

func (v *Value) Asc() types.Sortable {
	return sortcolumn.NewAsc(v)
}

// NewValue returns an AST node representing a Value
func NewValue(sel types.Selection, val interface{}) *Value {
	return &Value{
		sel: sel,
		val: val,
	}
}
