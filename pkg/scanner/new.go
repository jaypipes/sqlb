//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package scanner

import (
	"github.com/jaypipes/sqlb/pkg/types"
)

// The struct that holds information about the formatting and dialect of the
// output SQL that sqlb writes to the output buffer
//
// implements pkg/types.Scanner
type sqlScanner struct {
	dialect types.Dialect
	format  *types.FormatOptions
}

func (s *sqlScanner) Scan(b []byte, args []interface{}, scannables ...types.Scannable) {
	curArg := 0
	bw := 0
	bw += copy(b[bw:], s.format.PrefixWith)
	for _, scannable := range scannables {
		bw += scannable.Scan(s, b[bw:], args, &curArg)
	}
}

func (s *sqlScanner) Size(elements ...types.Element) *types.ElementSizes {
	buflen := 0
	argc := 0

	for _, el := range elements {
		argc += el.ArgCount()
		buflen += el.Size(s)
	}
	buflen += InterpolationLength(s.dialect, argc)
	buflen += len(s.format.PrefixWith)

	return &types.ElementSizes{
		ArgCount:   argc,
		BufferSize: buflen,
	}
}

// StringArgs returns the built query string and a slice of interface{}
// representing the values of the query args used in the query string, if any.
func (s *sqlScanner) StringArgs(el types.Element) (string, []interface{}) {
	sizes := s.Size(el)
	qargs := make([]interface{}, sizes.ArgCount)
	b := make([]byte, sizes.BufferSize)
	s.Scan(b, qargs, el)
	return string(b), qargs
}

func (s *sqlScanner) Dialect() types.Dialect {
	return s.dialect
}

func (s *sqlScanner) WithDialect(dialect types.Dialect) types.Scanner {
	s.dialect = dialect
	return s
}

func (s *sqlScanner) FormatOptions() *types.FormatOptions {
	return s.format
}

func (s *sqlScanner) WithFormatOptions(opts *types.FormatOptions) types.Scanner {
	s.format = opts
	return s
}

// New returns a scanner for the supplied dialect
func New(dialect types.Dialect) types.Scanner {
	return &sqlScanner{
		dialect: dialect,
		format:  DefaultFormatOptions,
	}
}
