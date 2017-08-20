package sqlb

import (
    "hash/fnv"
)

// Given one or more names, returns a hash of the names
func toId(names ...string) uint64 {
    hasher := fnv.New64a()
    for _, name := range names {
        hasher.Write([]byte(name))
    }
    return hasher.Sum64()
}

// Given a slice of interface{} variables, returns a slice of Element members.
// If any of the interface{} variables are *not* of type Element already, we
// construct a Value{} for the variable.
func toElements(vars ...interface{}) []Element {
    els := make([]Element, len(vars))
    for x, v := range vars {
        switch v.(type) {
        case Element:
            els[x] = v.(Element)
        default:
            els[x] = &Value{value: v}
        }
    }
    return els
}

// Given a variable number of interface{} variables, returns a List containing
// Value structs for the variables
// If any of the interface{} variables are *not* of type Element already, we
// construct a Value{} for the variable.
func toValueList(vars ...interface{}) *List {
    els := make([]Element, len(vars))
    for x, v := range vars {
        els[x] = &Value{value: v}
    }
    return &List{elements: els}
}
