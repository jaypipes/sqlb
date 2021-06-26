//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sqlb

import (
	"github.com/jaypipes/sqlb/pkg/types"
)

var defaultFormatOptions = &types.FormatOptions{
	SeparateClauseWith: " ",
	PrefixWith:         "",
}

var defaultScanner = &sqlScanner{
	dialect: types.DIALECT_MYSQL,
	format:  defaultFormatOptions,
}

// The struct that holds information about the formatting and dialect of the
// output SQL that sqlb writes to the output buffer
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
	buflen += interpolationLength(s.dialect, argc)
	buflen += len(s.format.PrefixWith)

	return &types.ElementSizes{
		ArgCount:   argc,
		BufferSize: buflen,
	}
}

func (s *sqlScanner) Dialect() types.Dialect {
	return s.dialect
}

func (s *sqlScanner) FormatOptions() *types.FormatOptions {
	return s.format
}
