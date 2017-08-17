package sqlb

// Given a slice of interface{} variables, returns a slice of Element members.
// If any of the interface{} variables are *not* of type Element already, we
// construct a Literal{} for the variable.
func toElements(vars ...interface{}) []Element {
    els := make([]Element, len(vars))
    for x, v := range vars {
        switch v.(type) {
        case Element:
            els[x] = v.(Element)
        default:
            els[x] = &Literal{value: v}
        }
    }
    return els
}
