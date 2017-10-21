package sqlb

import (
	"fmt"
)

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

func (v *value) projectionId() uint64 {
	// Each construction of a value is unique, so here we cheat and just
	// return the hash of the struct's address in memory
	return toId(fmt.Sprintf("%p", v))
}

func (v *value) argCount() int {
	return 1
}

func (v *value) size() int {
	size := 1 // the question mark for the value
	if v.alias != "" {
		size += len(Symbols[SYM_AS]) + len(v.alias)
	}
	return size
}

func (v *value) scan(b []byte, args []interface{}) (int, int) {
	args[0] = v.val
	bw := copy(b, Symbols[SYM_QUEST_MARK])
	if v.alias != "" {
		bw += copy(b[bw:], Symbols[SYM_AS])
		bw += copy(b[bw:], v.alias)
	}
	return bw, 1
}
