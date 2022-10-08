//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package function

import (
	"strings"

	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/grammar/sortcolumn"
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

type TrimFunction struct {
	sel      types.Selection
	alias    string
	subject  types.Element
	chars    string
	location TrimLocation
}

func (f *TrimFunction) From() types.Selection {
	return f.sel
}

func (f *TrimFunction) DisableAliasScan() func() {
	origAlias := f.alias
	f.alias = ""
	return func() { f.alias = origAlias }
}

func (f *TrimFunction) As(alias string) types.Projection {
	aliased := &TrimFunction{
		sel:      f.sel,
		alias:    alias,
		subject:  f.subject,
		location: f.location,
		chars:    f.chars,
	}
	return aliased
}

func (f *TrimFunction) ArgCount() int {
	argc := f.subject.ArgCount()
	if f.chars != "" {
		argc++
	}
	return argc
}

func (f *TrimFunction) Desc() types.Sortable {
	return sortcolumn.NewDesc(f)
}

func (f *TrimFunction) Asc() types.Sortable {
	return sortcolumn.NewAsc(f)
}

// Helper function that returns the non-subject, non-interpolation size of the
// TRIM() function for MySQL variants
func TrimFunctionSizeMySQL(f *TrimFunction) int {
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
func TrimFunctionScanMySQL(f *TrimFunction, scanner types.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	switch f.location {
	case TRIM_LEADING:
		if f.chars == "" {
			b.Write(grammar.Symbols[grammar.SYM_LTRIM])
		} else {
			b.Write(grammar.Symbols[grammar.SYM_TRIM])
			b.Write(grammar.Symbols[grammar.SYM_LEADING])
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			pkgscanner.ScanInterpolationMarker(types.DIALECT_MYSQL, b, *curArg)
			args[*curArg] = f.chars
			*curArg++
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			b.Write(grammar.Symbols[grammar.SYM_FROM])
		}
		TrimFunctionScanSubject(f, scanner, b, args, curArg)
	case TRIM_TRAILING:
		if f.chars == "" {
			b.Write(grammar.Symbols[grammar.SYM_RTRIM])
		} else {
			b.Write(grammar.Symbols[grammar.SYM_TRIM])
			b.Write(grammar.Symbols[grammar.SYM_TRAILING])
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			pkgscanner.ScanInterpolationMarker(types.DIALECT_MYSQL, b, *curArg)
			args[*curArg] = f.chars
			*curArg++
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			b.Write(grammar.Symbols[grammar.SYM_FROM])
		}
		TrimFunctionScanSubject(f, scanner, b, args, curArg)
	case TRIM_BOTH:
		if f.chars == "" {
			b.Write(grammar.Symbols[grammar.SYM_TRIM])
		} else {
			b.Write(grammar.Symbols[grammar.SYM_TRIM])
			pkgscanner.ScanInterpolationMarker(types.DIALECT_MYSQL, b, *curArg)
			args[*curArg] = f.chars
			*curArg++
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			b.Write(grammar.Symbols[grammar.SYM_FROM])
		}
		TrimFunctionScanSubject(f, scanner, b, args, curArg)
	}
}

// Helper function that returns the non-subject, non-interpolation size of the
// TRIM() function for PostgreSQL variants
func TrimFunctionSizePostgreSQL(f *TrimFunction) int {
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
func TrimFunctionScanPostgreSQL(f *TrimFunction, scanner types.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	switch f.location {
	case TRIM_LEADING:
		b.Write(grammar.Symbols[grammar.SYM_TRIM])
		b.Write(grammar.Symbols[grammar.SYM_LEADING])
		if f.chars != "" {
			b.WriteRune(' ')
			pkgscanner.ScanInterpolationMarker(types.DIALECT_POSTGRESQL, b, *curArg)
			args[*curArg] = f.chars
			*curArg++
		}
		b.Write(grammar.Symbols[grammar.SYM_SPACE])
		b.Write(grammar.Symbols[grammar.SYM_FROM])
		TrimFunctionScanSubject(f, scanner, b, args, curArg)
	case TRIM_TRAILING:
		b.Write(grammar.Symbols[grammar.SYM_TRIM])
		b.Write(grammar.Symbols[grammar.SYM_TRAILING])
		if f.chars != "" {
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			pkgscanner.ScanInterpolationMarker(types.DIALECT_POSTGRESQL, b, *curArg)
			args[*curArg] = f.chars
			*curArg++
		}
		b.Write(grammar.Symbols[grammar.SYM_SPACE])
		b.Write(grammar.Symbols[grammar.SYM_FROM])
		TrimFunctionScanSubject(f, scanner, b, args, curArg)
	case TRIM_BOTH:
		b.Write(grammar.Symbols[grammar.SYM_BTRIM])
		TrimFunctionScanSubject(f, scanner, b, args, curArg)
		if f.chars != "" {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
			pkgscanner.ScanInterpolationMarker(types.DIALECT_POSTGRESQL, b, *curArg)
			args[*curArg] = f.chars
			*curArg++
		}
	}
}

// Scan in the subject of the TRIM() function
func TrimFunctionScanSubject(f *TrimFunction, scanner types.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	// We need to disable alias output for elements that are
	// projections. We don't want to output, for example,
	// "ON users.id AS user_id = TRIM(articles.author)"
	switch f.subject.(type) {
	case types.Projection:
		reset := f.subject.(types.Projection).DisableAliasScan()
		defer reset()
	}
	f.subject.Scan(scanner, b, args, curArg)
}

func (f *TrimFunction) Size(scanner types.Scanner) int {
	size := 0
	switch scanner.Dialect() {
	case types.DIALECT_POSTGRESQL:
		size = TrimFunctionSizePostgreSQL(f)
	default:
		size = TrimFunctionSizeMySQL(f)
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

func (f *TrimFunction) Scan(scanner types.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	switch scanner.Dialect() {
	case types.DIALECT_POSTGRESQL:
		TrimFunctionScanPostgreSQL(f, scanner, b, args, curArg)
	default:
		TrimFunctionScanMySQL(f, scanner, b, args, curArg)
	}
	b.Write(grammar.Symbols[grammar.SYM_RPAREN])
	if f.alias != "" {
		b.Write(grammar.Symbols[grammar.SYM_AS])
		b.WriteString(f.alias)
	}
}

// Returns a struct that will output the TRIM() SQL function, trimming leading
// and trailing whitespace from the supplied projection
func Trim(p types.Projection) types.Projection {
	return &TrimFunction{
		subject:  p.(types.Element),
		sel:      p.From(),
		location: TRIM_BOTH,
	}
}

// Returns a struct that will output the LTRIM() SQL function for MySQL and the
// TRIM(LEADING FROM column) SQL function for PostgreSQL. The SQL function in
// either case will remove whitespace from the start of the supplied projection
func LTrim(p types.Projection) types.Projection {
	return &TrimFunction{
		subject:  p.(types.Element),
		sel:      p.From(),
		location: TRIM_LEADING,
	}
}

// Returns a struct that will output the RTRIM() SQL function for MySQL and the
// TRIM(TRAILING FROM column) SQL function for PostgreSQL. The SQL function in
// either case will remove whitespace from the start of the supplied projection
func RTrim(p types.Projection) types.Projection {
	return &TrimFunction{
		subject:  p.(types.Element),
		sel:      p.From(),
		location: TRIM_TRAILING,
	}
}

// Returns a struct that will output the TRIM() SQL function, trimming leading
// and trailing specified characters from the supplied projection
func TrimChars(p types.Projection, chars string) types.Projection {
	return &TrimFunction{
		subject:  p.(types.Element),
		sel:      p.From(),
		location: TRIM_BOTH,
		chars:    chars,
	}
}

// Returns a struct that will output the TRIM(LEADING chars FROM column) SQL
// function, trimming leading specified characters from the supplied projection
func LTrimChars(p types.Projection, chars string) types.Projection {
	return &TrimFunction{
		subject:  p.(types.Element),
		sel:      p.From(),
		location: TRIM_LEADING,
		chars:    chars,
	}
}

// Returns a struct that will output the TRIM(TRAILING chars FROM column) SQL
// function, trimming trailing specified characters from the supplied
// projection
func RTrimChars(p types.Projection, chars string) types.Projection {
	return &TrimFunction{
		subject:  p.(types.Element),
		sel:      p.From(),
		location: TRIM_TRAILING,
		chars:    chars,
	}
}
