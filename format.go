package sqlb

type FormatOptions struct {
	SeparateClauseWith string
}

var defaultFormatOptions = &FormatOptions{
	SeparateClauseWith: " ",
}

// The struct that holds information about the formatting and dialect of the
// output SQL that sqlb writes to the output buffer
type sqlScanner struct {
	dialect Dialect
	format  *FormatOptions
}

func (s *sqlScanner) scan(b []byte, args []interface{}, scannables ...Scannable) {
	curArg := 0
	scannables[0].scan(b, args, &curArg)
}

// Returns the length (in bytes) of the interpolation markers, which depends on
// the dialect in use when constructing the SQL buffer
func (s *sqlScanner) interpolationLength(argc int) int {
	return interpolationLength(s.dialect, argc)
}
