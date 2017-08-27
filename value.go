package sqlb

// A value is a concrete struct wrapper around a constant that implements the
// scannable interface. Typically, users won't directly construct value
// structs but instead helper functions like sqlb.Equal() will construct a
// value and bind it to the containing element.
type value struct {
    val interface{}
}

func (v *value) argCount() int {
    return 1
}

func (v  *value) size() int {
    // The value is always injected as a question mark in the produced SQL
    // string
    return 1
}

func (v *value) scan(b []byte, args []interface{}) (int, int) {
    args[0] = v.val
    copy(b, Symbols[SYM_QUEST_MARK])
    return 1, 1
}
