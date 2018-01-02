//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

// A value is a concrete struct wrapper around a constant that implements the
// scannable interface. Typically, users won't directly construct value
// structs but instead helper functions like sqlb.Equal() will construct a
// value and bind it to the containing element.
type value struct {
	sel   selection
	alias string
	val   interface{}
}

func (v *value) from() selection {
	return v.sel
}

func (v *value) setAlias(alias string) *value {
	v.alias = alias
	return v
}

func (v *value) As(alias string) *value {
	return &value{
		alias: alias,
		val:   v.val,
	}
}

func (v *value) disableAliasScan() func() {
	origAlias := v.alias
	v.alias = ""
	return func() { v.alias = origAlias }
}

func (v *value) argCount() int {
	return 1
}

func (v *value) size(scanner *sqlScanner) int {
	// Due to dialect handling, we do not include the length of interpolation
	// markers for query parameters. This is calculated separately by the
	// top-level scanning struct before malloc'ing the buffer to inject the SQL
	// string into.
	size := 0
	if v.alias != "" {
		size += len(Symbols[SYM_AS]) + len(v.alias)
	}
	return size
}

func (v *value) scan(scanner *sqlScanner, b []byte, args []interface{}, curArg *int) int {
	args[*curArg] = v.val
	bw := scanInterpolationMarker(scanner.dialect, b, *curArg)
	*curArg++
	if v.alias != "" {
		bw += copy(b[bw:], Symbols[SYM_AS])
		bw += copy(b[bw:], v.alias)
	}
	return bw
}
