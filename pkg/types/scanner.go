//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package types

type ElementSizes struct {
	// The number of interface{} arguments that the element will add to the
	// slice of interface{} arguments that will eventually be constructed for
	// the Query
	ArgCount int
	// The number of bytes in the output buffer to represent this element
	BufferSize int
}

type Scanner interface {
	// Scan takes two slices and a pointer to an int. The first slice is a slice of bytes that the
	// implementation should copy its string representation to and the other slice is a slice of interface{} values that the element should add its
	// arguments to. The pointer to an int is the index of the current argument to be processed. The method returns a single int, the number of bytes written to the buffer.
	Scan([]byte, []interface{}, ...Scannable)
	Size(...Element) *ElementSizes
	Dialect() Dialect
	WithDialect(Dialect) Scanner
	FormatOptions() *FormatOptions
	WithFormatOptions(*FormatOptions) Scanner
}
