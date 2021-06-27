//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	pkgscanner "github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/types"
)

// TRIM/BTRIM/LTRIM/RTRIM SQL function support
//
// For MySQL, the TRIM() SQL function takes the following forms.
//
// Remove whitespace from before and after string:
//		TRIM(string)
// Remove whitespace from before string:
//		LTRIM(string)
// Remove whitespace from after string:
//		RTRIM(string)
// Remove the string remstr from before and after string:
//		TRIM(remstr FROM string)
// Remove the string remstr from before string:
//		TRIM(LEADING remstr FROM string)
// Remove the string remstr from after string:
//		TRIM(TRAILING remstr FROM string)
//
// For PostgreSQL, the TRIM() SQL function takes the following forms.
//
// Remove whitespace from before and after string:
//		BTRIM(string)
// Remove whitespace from before string:
//		TRIM(LEADING FROM string)
// Remove whitespace from after string:
//		TRIM(TRAILING FROM string)
// Remove longest string containing any character in chars from before
// and after string:
//		BTRIM(string, chars)
// Remove longest string containing any character in chars from before
// string:
//		TRIM(LEADING chars FROM string)
// Remove longest string containing any character in chars from after
// string:
//		TRIM(TRAILING chars FROM string)

type TrimLocation int

const (
	TRIM_BOTH TrimLocation = iota
	TRIM_TRAILING
	TRIM_LEADING
)

type trimFunc struct {
	sel      types.Selection
	alias    string
	subject  types.Element
	chars    string
	location TrimLocation
}

func (f *trimFunc) From() types.Selection {
	return f.sel
}

func (f *trimFunc) DisableAliasScan() func() {
	origAlias := f.alias
	f.alias = ""
	return func() { f.alias = origAlias }
}

func (f *trimFunc) As(alias string) *trimFunc {
	aliased := &trimFunc{
		sel:      f.sel,
		alias:    alias,
		subject:  f.subject,
		location: f.location,
		chars:    f.chars,
	}
	return aliased
}

func (f *trimFunc) ArgCount() int {
	argc := f.subject.ArgCount()
	if f.chars != "" {
		argc++
	}
	return argc
}

// Helper function that returns the non-subject, non-interpolation size of the
// TRIM() function for MySQL variants
func trimFuncSizeMySQL(f *trimFunc) int {
	size := 0
	switch f.location {
	case TRIM_LEADING:
		if f.chars == "" {
			// LTRIM(string)
			size = len(grammar.Symbols[grammar.SYM_LTRIM])
		} else {
			// TRIM(LEADING remstr FROM string)
			size = (len(grammar.Symbols[grammar.SYM_TRIM]) + len(grammar.Symbols[grammar.SYM_LEADING]) +
				len(grammar.Symbols[grammar.SYM_FROM]) + 2)
		}
	case TRIM_TRAILING:
		if f.chars == "" {
			// LTRIM(string)
			size = len(grammar.Symbols[grammar.SYM_RTRIM])
		} else {
			// TRIM(TRAILING remstr FROM string)
			size = (len(grammar.Symbols[grammar.SYM_TRIM]) + len(grammar.Symbols[grammar.SYM_TRAILING]) +
				len(grammar.Symbols[grammar.SYM_FROM]) + 2)
		}
	case TRIM_BOTH:
		if f.chars == "" {
			// TRIM(string)
			size = len(grammar.Symbols[grammar.SYM_TRIM])
		} else {
			// TRIM(remstr FROM string)
			size = len(grammar.Symbols[grammar.SYM_TRIM]) + len(grammar.Symbols[grammar.SYM_FROM]) + 1
		}
	}
	return size
}

// Helper function that scans into the supplied SQL []byte buffer for the
// TRIM/BTRIM() SQL function for MySQL
// TRIM() function for MySQL variants
func trimFuncScanMySQL(f *trimFunc, scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	switch f.location {
	case TRIM_LEADING:
		if f.chars == "" {
			bw += copy(b[bw:], grammar.Symbols[grammar.SYM_LTRIM])
		} else {
			bw += copy(b[bw:], grammar.Symbols[grammar.SYM_TRIM])
			bw += copy(b[bw:], grammar.Symbols[grammar.SYM_LEADING])
			bw += copy(b[bw:], " ")
			bw += pkgscanner.ScanInterpolationMarker(types.DIALECT_MYSQL, b[bw:], *curArg)
			args[*curArg] = f.chars
			*curArg++
			bw += copy(b[bw:], " ")
			bw += copy(b[bw:], grammar.Symbols[grammar.SYM_FROM])
		}
		bw += trimFuncScanSubject(f, scanner, b[bw:], args, curArg)
	case TRIM_TRAILING:
		if f.chars == "" {
			bw += copy(b[bw:], grammar.Symbols[grammar.SYM_RTRIM])
		} else {
			bw += copy(b[bw:], grammar.Symbols[grammar.SYM_TRIM])
			bw += copy(b[bw:], grammar.Symbols[grammar.SYM_TRAILING])
			bw += copy(b[bw:], " ")
			bw += pkgscanner.ScanInterpolationMarker(types.DIALECT_MYSQL, b[bw:], *curArg)
			args[*curArg] = f.chars
			*curArg++
			bw += copy(b[bw:], " ")
			bw += copy(b[bw:], grammar.Symbols[grammar.SYM_FROM])
		}
		bw += trimFuncScanSubject(f, scanner, b[bw:], args, curArg)
	case TRIM_BOTH:
		if f.chars == "" {
			bw += copy(b[bw:], grammar.Symbols[grammar.SYM_TRIM])
		} else {
			bw += copy(b[bw:], grammar.Symbols[grammar.SYM_TRIM])
			bw += pkgscanner.ScanInterpolationMarker(types.DIALECT_MYSQL, b[bw:], *curArg)
			args[*curArg] = f.chars
			*curArg++
			bw += copy(b[bw:], " ")
			bw += copy(b[bw:], grammar.Symbols[grammar.SYM_FROM])
		}
		bw += trimFuncScanSubject(f, scanner, b[bw:], args, curArg)
	}
	return bw
}

// Helper function that returns the non-subject, non-interpolation size of the
// TRIM() function for PostgreSQL variants
func trimFuncSizePostgreSQL(f *trimFunc) int {
	size := 0
	switch f.location {
	case TRIM_LEADING:
		// TRIM(LEADING FROM string)
		size = (len(grammar.Symbols[grammar.SYM_TRIM]) + len(grammar.Symbols[grammar.SYM_LEADING]) +
			len(grammar.Symbols[grammar.SYM_FROM]) + 1)
		if f.chars != "" {
			// TRIM(LEADING chars FROM string)
			size += 1
		}
	case TRIM_TRAILING:
		// TRIM(TRAILING FROM string)
		size = (len(grammar.Symbols[grammar.SYM_TRIM]) + len(grammar.Symbols[grammar.SYM_TRAILING]) +
			len(grammar.Symbols[grammar.SYM_FROM]) + 1)
		if f.chars != "" {
			// TRIM(TRAILING chars FROM string)
			size += 1
		}
	case TRIM_BOTH:
		if f.chars == "" {
			// BTRIM(string)
			size = len(grammar.Symbols[grammar.SYM_BTRIM])
		} else {
			// BTRIM(string, chars)
			size = len(grammar.Symbols[grammar.SYM_BTRIM]) + len(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}
	return size
}

// Helper function that scans into the supplied SQL []byte buffer for the
// TRIM/BTRIM() SQL function for PostgreSQL
// TRIM() function for PostgreSQL variants
func trimFuncScanPostgreSQL(f *trimFunc, scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	switch f.location {
	case TRIM_LEADING:
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_TRIM])
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_LEADING])
		if f.chars != "" {
			bw += copy(b[bw:], " ")
			bw += pkgscanner.ScanInterpolationMarker(types.DIALECT_POSTGRESQL, b[bw:], *curArg)
			args[*curArg] = f.chars
			*curArg++
		}
		bw += copy(b[bw:], " ")
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_FROM])
		bw += trimFuncScanSubject(f, scanner, b[bw:], args, curArg)
	case TRIM_TRAILING:
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_TRIM])
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_TRAILING])
		if f.chars != "" {
			bw += copy(b[bw:], " ")
			bw += pkgscanner.ScanInterpolationMarker(types.DIALECT_POSTGRESQL, b[bw:], *curArg)
			args[*curArg] = f.chars
			*curArg++
		}
		bw += copy(b[bw:], " ")
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_FROM])
		bw += trimFuncScanSubject(f, scanner, b[bw:], args, curArg)
	case TRIM_BOTH:
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_BTRIM])
		bw += trimFuncScanSubject(f, scanner, b[bw:], args, curArg)
		if f.chars != "" {
			bw += copy(b[bw:], grammar.Symbols[grammar.SYM_COMMA_WS])
			bw += pkgscanner.ScanInterpolationMarker(types.DIALECT_POSTGRESQL, b[bw:], *curArg)
			args[*curArg] = f.chars
			*curArg++
		}
	}
	return bw
}

// Scan in the subject of the TRIM() function
func trimFuncScanSubject(f *trimFunc, scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	// We need to disable alias output for elements that are
	// projections. We don't want to output, for example,
	// "ON users.id AS user_id = TRIM(articles.author)"
	switch f.subject.(type) {
	case types.Projection:
		reset := f.subject.(types.Projection).DisableAliasScan()
		defer reset()
	}
	return f.subject.Scan(scanner, b, args, curArg)
}

func (f *trimFunc) Size(scanner types.Scanner) int {
	size := 0
	switch scanner.Dialect() {
	case types.DIALECT_POSTGRESQL:
		size = trimFuncSizePostgreSQL(f)
	default:
		size = trimFuncSizeMySQL(f)
	}
	size += len(grammar.Symbols[grammar.SYM_RPAREN])
	// We need to disable alias output for elements that are
	// projections. We don't want to output, for example,
	// "ON users.id AS user_id = TRIM(articles.author)"
	switch f.subject.(type) {
	case types.Projection:
		reset := f.subject.(types.Projection).DisableAliasScan()
		defer reset()
	}
	size += f.subject.Size(scanner)
	if f.alias != "" {
		size += len(grammar.Symbols[grammar.SYM_AS]) + len(f.alias)
	}
	return size
}

func (f *trimFunc) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	switch scanner.Dialect() {
	case types.DIALECT_POSTGRESQL:
		bw += trimFuncScanPostgreSQL(f, scanner, b[bw:], args, curArg)
	default:
		bw += trimFuncScanMySQL(f, scanner, b[bw:], args, curArg)
	}
	bw += copy(b[bw:], grammar.Symbols[grammar.SYM_RPAREN])
	if f.alias != "" {
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_AS])
		bw += copy(b[bw:], f.alias)
	}
	return bw
}

// Returns a struct that will output the TRIM() SQL function, trimming leading
// and trailing whitespace from the supplied projection
func Trim(p types.Projection) *trimFunc {
	return &trimFunc{
		subject:  p.(types.Element),
		sel:      p.From(),
		location: TRIM_BOTH,
	}
}

func (c *Column) Trim() *trimFunc {
	f := Trim(c)
	return f
}

// Returns a struct that will output the LTRIM() SQL function for MySQL and the
// TRIM(LEADING FROM column) SQL function for PostgreSQL. The SQL function in
// either case will remove whitespace from the start of the supplied projection
func LTrim(p types.Projection) *trimFunc {
	return &trimFunc{
		subject:  p.(types.Element),
		sel:      p.From(),
		location: TRIM_LEADING,
	}
}

func (c *Column) LTrim() *trimFunc {
	f := LTrim(c)
	return f
}

// Returns a struct that will output the RTRIM() SQL function for MySQL and the
// TRIM(TRAILING FROM column) SQL function for PostgreSQL. The SQL function in
// either case will remove whitespace from the start of the supplied projection
func RTrim(p types.Projection) *trimFunc {
	return &trimFunc{
		subject:  p.(types.Element),
		sel:      p.From(),
		location: TRIM_TRAILING,
	}
}

func (c *Column) RTrim() *trimFunc {
	f := LTrim(c)
	return f
}

// Returns a struct that will output the TRIM() SQL function, trimming leading
// and trailing specified characters from the supplied projection
func TrimChars(p types.Projection, chars string) *trimFunc {
	return &trimFunc{
		subject:  p.(types.Element),
		sel:      p.From(),
		location: TRIM_BOTH,
		chars:    chars,
	}
}

func (c *Column) TrimChars(chars string) *trimFunc {
	f := TrimChars(c, chars)
	return f
}

// Returns a struct that will output the TRIM(LEADING chars FROM column) SQL
// function, trimming leading specified characters from the supplied projection
func LTrimChars(p types.Projection, chars string) *trimFunc {
	return &trimFunc{
		subject:  p.(types.Element),
		sel:      p.From(),
		location: TRIM_LEADING,
		chars:    chars,
	}
}

func (c *Column) LTrimChars(chars string) *trimFunc {
	f := LTrimChars(c, chars)
	return f
}

// Returns a struct that will output the TRIM(TRAILING chars FROM column) SQL
// function, trimming trailing specified characters from the supplied
// projection
func RTrimChars(p types.Projection, chars string) *trimFunc {
	return &trimFunc{
		subject:  p.(types.Element),
		sel:      p.From(),
		location: TRIM_TRAILING,
		chars:    chars,
	}
}

func (c *Column) RTrimChars(chars string) *trimFunc {
	f := RTrimChars(c, chars)
	return f
}
