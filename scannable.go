package sqlb

type Scannable interface {
    Size() int
    // Scan takes a slice of bytes that the implementation should copy its
    // string representation to and return the number of bytes that it copied
    // into the byte slice
    Scan([]byte) int
}
