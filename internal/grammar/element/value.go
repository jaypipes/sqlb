//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package element

import (
	"strings"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/sortcolumn"
)

// Value is a concrete struct wrapper around a constant that implements the
// Scannable interface. Typically, users won't directly construct value
// structs but instead helper functions like sqlb.Equal() will construct a
// value and bind it to the containing element.
type Value struct {
	sel   api.Selection
	alias string
	val   interface{}
}

func (v *Value) From() api.Selection {
	return v.sel
}

func (v *Value) As(alias string) api.Projection {
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

func (v *Value) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	qargs[*curarg] = v.val
	b.WriteString(builder.InterpolationMarker(opts, *curarg))
	*curarg++
	if v.alias != "" {
		b.Write(grammar.Symbols[grammar.SYM_AS])
		b.WriteString(v.alias)
	}
	return b.String()
}

func (v *Value) Desc() api.Orderable {
	return sortcolumn.NewDesc(v)
}

func (v *Value) Asc() api.Orderable {
	return sortcolumn.NewAsc(v)
}

// NewValue returns an AST node representing a Value
func NewValue(sel api.Selection, val interface{}) *Value {
	return &Value{
		sel: sel,
		val: val,
	}
}
