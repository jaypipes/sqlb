//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

/*
// Ascii returns a Projection that contains the ASCII() SQL function
func Ascii(p api.Projection) api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_ASCII),
		elements: []api.Element{p.(api.Element)},
		sel:      p.From(),
	}
}

// Reverse returns a Projection that contains the REVERSE() SQL function
func Reverse(p api.Projection) api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_REVERSE),
		elements: []api.Element{p.(api.Element)},
		sel:      p.From(),
	}
}

// Concat returns a Projection that contains the CONCAT() SQL function
func Concat(projs ...api.Projection) api.Projection {
	els := make([]api.Element, len(projs))
	for x, p := range projs {
		els[x] = p.(api.Element)
	}
	subjects := element.NewList(els...)
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CONCAT),
		elements: []api.Element{subjects},
		// TODO(jaypipes): Clearly we need to support >1 selection...
		sel: projs[0].From(),
	}
}

// ConcatWs returns a Projection that contains the CONCAT_WS() SQL function
func ConcatWs(sep string, projs ...api.Projection) api.Projection {
	els := make([]api.Element, len(projs))
	for x, p := range projs {
		els[x] = p.(api.Element)
	}
	subjects := element.NewList(els...)
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CONCAT_WS),
		elements: []api.Element{element.NewValue(nil, sep), subjects},
		// TODO(jaypipes): Clearly we need to support >1 selection...
		sel: projs[0].From(),
	}
}
*/

/*
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
	sel      api.Selection
	alias    string
	subject  api.Element
	chars    string
	location TrimLocation
}

func (f *TrimFunction) From() api.Selection {
	return f.sel
}

func (f *TrimFunction) DisableAliasScan() func() {
	origAlias := f.alias
	f.alias = ""
	return func() { f.alias = origAlias }
}

func (f *TrimFunction) As(alias string) api.Projection {
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

func (f *TrimFunction) Desc() api.Orderable {
	return sortcolumn.NewDesc(f)
}

func (f *TrimFunction) Asc() api.Orderable {
	return sortcolumn.NewAsc(f)
}

// Helper function that scans into the supplied SQL []byte buffer for the
// TRIM/BTRIM() SQL function for MySQL
// TRIM() function for MySQL variants
func TrimFunctionScanMySQL(
	f *TrimFunction,
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	switch f.location {
	case TRIM_LEADING:
		if f.chars == "" {
			b.Write(grammar.Symbols[grammar.SYM_LTRIM])
		} else {
			b.Write(grammar.Symbols[grammar.SYM_TRIM])
			b.Write(grammar.Symbols[grammar.SYM_LEADING])
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			b.WriteString(builder.InterpolationMarker(opts, *curarg))
			qargs[*curarg] = f.chars
			*curarg++
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			b.Write(grammar.Symbols[grammar.SYM_FROM])
		}
		b.WriteString(TrimFunctionScanSubject(f, opts, qargs, curarg))
	case TRIM_TRAILING:
		if f.chars == "" {
			b.Write(grammar.Symbols[grammar.SYM_RTRIM])
		} else {
			b.Write(grammar.Symbols[grammar.SYM_TRIM])
			b.Write(grammar.Symbols[grammar.SYM_TRAILING])
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			b.WriteString(builder.InterpolationMarker(opts, *curarg))
			qargs[*curarg] = f.chars
			*curarg++
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			b.Write(grammar.Symbols[grammar.SYM_FROM])
		}
		b.WriteString(TrimFunctionScanSubject(f, opts, qargs, curarg))
	case TRIM_BOTH:
		if f.chars == "" {
			b.Write(grammar.Symbols[grammar.SYM_TRIM])
		} else {
			b.Write(grammar.Symbols[grammar.SYM_TRIM])
			b.WriteString(builder.InterpolationMarker(opts, *curarg))
			qargs[*curarg] = f.chars
			*curarg++
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			b.Write(grammar.Symbols[grammar.SYM_FROM])
		}
		b.WriteString(TrimFunctionScanSubject(f, opts, qargs, curarg))
	}
	return b.String()
}

// Helper function that scans into the supplied SQL []byte buffer for the
// TRIM/BTRIM() SQL function for PostgreSQL
// TRIM() function for PostgreSQL variants
func TrimFunctionScanPostgreSQL(
	f *TrimFunction,
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	switch f.location {
	case TRIM_LEADING:
		b.Write(grammar.Symbols[grammar.SYM_TRIM])
		b.Write(grammar.Symbols[grammar.SYM_LEADING])
		if f.chars != "" {
			b.WriteRune(' ')
			b.WriteString(builder.InterpolationMarker(opts, *curarg))
			qargs[*curarg] = f.chars
			*curarg++
		}
		b.Write(grammar.Symbols[grammar.SYM_SPACE])
		b.Write(grammar.Symbols[grammar.SYM_FROM])
		b.WriteString(TrimFunctionScanSubject(f, opts, qargs, curarg))
	case TRIM_TRAILING:
		b.Write(grammar.Symbols[grammar.SYM_TRIM])
		b.Write(grammar.Symbols[grammar.SYM_TRAILING])
		if f.chars != "" {
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			b.WriteString(builder.InterpolationMarker(opts, *curarg))
			qargs[*curarg] = f.chars
			*curarg++
		}
		b.Write(grammar.Symbols[grammar.SYM_SPACE])
		b.Write(grammar.Symbols[grammar.SYM_FROM])
		b.WriteString(TrimFunctionScanSubject(f, opts, qargs, curarg))
	case TRIM_BOTH:
		b.Write(grammar.Symbols[grammar.SYM_BTRIM])
		b.WriteString(TrimFunctionScanSubject(f, opts, qargs, curarg))
		if f.chars != "" {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
			b.WriteString(builder.InterpolationMarker(opts, *curarg))
			qargs[*curarg] = f.chars
			*curarg++
		}
	}
	return b.String()
}

// Scan in the subject of the TRIM() function
func TrimFunctionScanSubject(
	f *TrimFunction,
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	// We need to disable alias output for elements that are
	// projections. We don't want to output, for example,
	// "ON users.id AS user_id = TRIM(articles.author)"
	switch f.subject.(type) {
	case api.Projection:
		reset := f.subject.(api.Projection).DisableAliasScan()
		defer reset()
	}
	b.WriteString(f.subject.String(opts, qargs, curarg))
	return b.String()
}

func (f *TrimFunction) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	switch opts.Dialect() {
	case api.DialectPostgreSQL:
		b.WriteString(TrimFunctionScanPostgreSQL(f, opts, qargs, curarg))
	default:
		b.WriteString(TrimFunctionScanMySQL(f, opts, qargs, curarg))
	}
	b.Write(grammar.Symbols[grammar.SYM_RPAREN])
	if f.alias != "" {
		b.Write(grammar.Symbols[grammar.SYM_AS])
		b.WriteString(f.alias)
	}
	return b.String()
}

// Returns a struct that will output the TRIM() SQL function, trimming leading
// and trailing whitespace from the supplied projection
func Trim(p api.Projection) api.Projection {
	return &TrimFunction{
		subject:  p.(api.Element),
		sel:      p.From(),
		location: TRIM_BOTH,
	}
}

// Returns a struct that will output the LTRIM() SQL function for MySQL and the
// TRIM(LEADING FROM column) SQL function for PostgreSQL. The SQL function in
// either case will remove whitespace from the start of the supplied projection
func LTrim(p api.Projection) api.Projection {
	return &TrimFunction{
		subject:  p.(api.Element),
		sel:      p.From(),
		location: TRIM_LEADING,
	}
}

// Returns a struct that will output the RTRIM() SQL function for MySQL and the
// TRIM(TRAILING FROM column) SQL function for PostgreSQL. The SQL function in
// either case will remove whitespace from the start of the supplied projection
func RTrim(p api.Projection) api.Projection {
	return &TrimFunction{
		subject:  p.(api.Element),
		sel:      p.From(),
		location: TRIM_TRAILING,
	}
}

// Returns a struct that will output the TRIM() SQL function, trimming leading
// and trailing specified characters from the supplied projection
func TrimChars(p api.Projection, chars string) api.Projection {
	return &TrimFunction{
		subject:  p.(api.Element),
		sel:      p.From(),
		location: TRIM_BOTH,
		chars:    chars,
	}
}

// Returns a struct that will output the TRIM(LEADING chars FROM column) SQL
// function, trimming leading specified characters from the supplied projection
func LTrimChars(p api.Projection, chars string) api.Projection {
	return &TrimFunction{
		subject:  p.(api.Element),
		sel:      p.From(),
		location: TRIM_LEADING,
		chars:    chars,
	}
}

// Returns a struct that will output the TRIM(TRAILING chars FROM column) SQL
// function, trimming trailing specified characters from the supplied
// projection
func RTrimChars(p api.Projection, chars string) api.Projection {
	return &TrimFunction{
		subject:  p.(api.Element),
		sel:      p.From(),
		location: TRIM_TRAILING,
		chars:    chars,
	}
}
*/
