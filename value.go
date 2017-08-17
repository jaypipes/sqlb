package sqlb

// A Value is a concrete struct wrapper around a constant that implements the
// Scannable interface. Typically, users won't directly construct Value
// structs but instead helper functions like sqlb.Equal() will construct a
// Value and bind it to the containing Element.
type Value struct {
    value interface{}
}

func (val *Value) ArgCount() int {
    return 1
}

func (val  *Value) Size() int {
    // The value is always injected as a question mark in the produced SQL
    // string
    return 1
}

func (val *Value) Scan(b []byte, args []interface{}) (int, int) {
    args[0] = val.value
    copy(b, SYM_QM)
    return 1, 1
}
