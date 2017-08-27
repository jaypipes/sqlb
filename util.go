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

// Given a slice of interface{} variables, returns a slice of element members.
// If any of the interface{} variables are *not* of type element already, we
// construct a Value{} for the variable.
func toElements(vars ...interface{}) []element {
    els := make([]element, len(vars))
    for x, v := range vars {
        switch v.(type) {
        case element:
            els[x] = v.(element)
        default:
            els[x] = &value{val: v}
        }
    }
    return els
}

// Given a variable number of interface{} variables, returns a List containing
// Value structs for the variables
// If any of the interface{} variables are *not* of type element already, we
// construct a Value{} for the variable.
func toValueList(vars ...interface{}) *List {
    els := make([]element, len(vars))
    for x, v := range vars {
        els[x] = &value{val: v}
    }
    return &List{elements: els}
}
