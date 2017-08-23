package sqlb

// A Value is a concrete struct wrapper around a constant that implements the
// scannable interface. Typically, users won't directly construct Value
// structs but instead helper functions like sqlb.Equal() will construct a
// Value and bind it to the containing element.
type Value struct {
    value interface{}
}

func (val *Value) argCount() int {
    return 1
}

func (val  *Value) size() int {
    // The value is always injected as a question mark in the produced SQL
    // string
    return 1
}

func (val *Value) scan(b []byte, args []interface{}) (int, int) {
    args[0] = val.value
    copy(b, Symbols[SYM_QUEST_MARK])
    return 1, 1
}
