package sqlb

type FormatOptions struct {
	SeparateClauseWith string
}

var defaultFormatOptions = &FormatOptions{
	SeparateClauseWith: " ",
}

type ElementSizes struct {
	// The number of interface{} arguments that the element will add to the
	// slice of interface{} arguments that will eventually be constructed for
	// the Query
	ArgCount int
	// The number of bytes in the output buffer to represent this element
	BufferSize int
}

// The struct that holds information about the formatting and dialect of the
// output SQL that sqlb writes to the output buffer
type sqlScanner struct {
	dialect Dialect
	format  *FormatOptions
}

func (s *sqlScanner) scan(b []byte, args []interface{}, scannables ...Scannable) {
	curArg := 0

	for _, scannable := range scannables {
		scannable.scan(b, args, &curArg)
	}
}

func (s *sqlScanner) size(elements ...element) *ElementSizes {
	buflen := 0
	argc := 0

	for _, el := range elements {
		argc += el.argCount()
		buflen += el.size()
	}
	buflen += interpolationLength(s.dialect, argc)

	return &ElementSizes{
		ArgCount:   argc,
		BufferSize: buflen,
	}
}
