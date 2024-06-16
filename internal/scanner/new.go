//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package scanner

import (
	"strings"

	"github.com/jaypipes/sqlb/types"
)

type ElementSizes struct {
	// The number of interface{} arguments that the element will add to the
	// slice of interface{} arguments that will eventually be constructed for
	// the Query
	ArgCount int
	// The number of bytes in the output buffer to represent this element
	BufferSize int
}

// Scanner holds information about the formatting and dialect of the output SQL
// that sqlb writes to the output buffer
type Scanner struct {
	Dialect types.Dialect
	Format  types.FormatOptions
}

// Scan takes two slices and a pointer to an int. The first slice is a slice of
// bytes that the implementation should copy its string representation to and
// the other slice is a slice of interface{} values that the element should add
// its arguments to. The pointer to an int is the index of the current argument
// to be processed. The method returns a single int, the number of bytes
// written to the buffer.
func (s *Scanner) Scan(b *strings.Builder, args []interface{}, scannables ...Scannable) {
	curArg := 0
	b.WriteString(s.Format.PrefixWith)
	for _, scannable := range scannables {
		scannable.Scan(s, b, args, &curArg)
	}
}

func (s *Scanner) Size(elements ...Element) *ElementSizes {
	buflen := 0
	argc := 0

	for _, el := range elements {
		argc += el.ArgCount()
		buflen += el.Size(s)
	}
	buflen += InterpolationLength(s.Dialect, argc)
	buflen += len(s.Format.PrefixWith)

	return &ElementSizes{
		ArgCount:   argc,
		BufferSize: buflen,
	}
}

// StringArgs returns the built query string and a slice of interface{}
// representing the values of the query args used in the query string, if any.
func (s *Scanner) StringArgs(el Element) (string, []interface{}) {
	sizes := s.Size(el)
	qargs := make([]interface{}, sizes.ArgCount)
	var b strings.Builder
	b.Grow(sizes.BufferSize)
	s.Scan(&b, qargs, el)
	return b.String(), qargs
}

// New returns a scanner for the supplied dialect
func New(
	mods ...ScannerOptionModifier,
) *Scanner {
	opts := mergeOpts(mods)
	return &Scanner{
		Dialect: opts.Dialect,
		Format:  opts.Format,
	}
}
