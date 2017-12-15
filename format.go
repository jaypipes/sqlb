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

func (w *sqlScanner) scan(b []byte, args []interface{}, scannables ...Scannable) {
	curArg := 0
	scannables[0].scan(b, args, &curArg)
}
