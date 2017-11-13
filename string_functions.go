package sqlb

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
	sel      selection
	alias    string
	subject  element
	dialect  Dialect
	chars    string
	location TrimLocation
}

// Sets the element's dialect and pushes the dialect down into any of the
// sub-elements
func (f *trimFunc) setDialect(dialect Dialect) {
	f.dialect = dialect
	switch f.subject.(type) {
	case *value:
		v := f.subject.(*value)
		v.dialect = dialect
	}
}

func (f *trimFunc) from() selection {
	return f.sel
}

func (f *trimFunc) disableAliasScan() func() {
	origAlias := f.alias
	f.alias = ""
	return func() { f.alias = origAlias }
}

func (f *trimFunc) As(alias string) *trimFunc {
	aliased := &trimFunc{
		sel:      f.sel,
		alias:    alias,
		subject:  f.subject,
		dialect:  f.dialect,
		location: f.location,
		chars:    f.chars,
	}
	return aliased
}

func (f *trimFunc) argCount() int {
	argc := f.subject.argCount()
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
			size = len(Symbols[SYM_LTRIM])
		} else {
			// TRIM(LEADING remstr FROM string)
			size = (len(Symbols[SYM_TRIM]) + len(Symbols[SYM_LEADING]) +
				len(Symbols[SYM_FROM]) + 1)
		}
	case TRIM_TRAILING:
		if f.chars == "" {
			// LTRIM(string)
			size = len(Symbols[SYM_RTRIM])
		} else {
			// TRIM(TRAILING remstr FROM string)
			size = (len(Symbols[SYM_TRIM]) + len(Symbols[SYM_TRAILING]) +
				len(Symbols[SYM_FROM]) + 1)
		}
	case TRIM_BOTH:
		if f.chars == "" {
			// TRIM(string)
			size = len(Symbols[SYM_TRIM])
		} else {
			// TRIM(remstr FROM string)
			size = len(Symbols[SYM_TRIM]) + len(Symbols[SYM_FROM])
		}
	}
	return size
}

// Helper function that scans into the supplied SQL []byte buffer for the
// TRIM/BTRIM() SQL function for MySQL
// TRIM() function for MySQL variants
func trimFuncScanMySQL(f *trimFunc, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	switch f.location {
	case TRIM_LEADING:
		if f.chars == "" {
			bw += copy(b[bw:], Symbols[SYM_LTRIM])
		} else {
			bw += copy(b[bw:], Symbols[SYM_TRIM])
			bw += copy(b[bw:], Symbols[SYM_LEADING])
			bw += copy(b[bw:], []byte{' '})
			bw += scanInterpolationMarker(f.dialect, b[bw:], *curArg)
			args[*curArg] = f.chars
			*curArg++
			bw += copy(b[bw:], Symbols[SYM_FROM])
		}
		bw += trimFuncScanSubject(f, b[bw:], args, curArg)
	case TRIM_TRAILING:
		if f.chars == "" {
			bw += copy(b[bw:], Symbols[SYM_RTRIM])
		} else {
			bw += copy(b[bw:], Symbols[SYM_TRIM])
			bw += copy(b[bw:], Symbols[SYM_TRAILING])
			bw += copy(b[bw:], []byte{' '})
			bw += scanInterpolationMarker(f.dialect, b[bw:], *curArg)
			args[*curArg] = f.chars
			*curArg++
			bw += copy(b[bw:], Symbols[SYM_FROM])
		}
		bw += trimFuncScanSubject(f, b[bw:], args, curArg)
	case TRIM_BOTH:
		if f.chars == "" {
			bw += copy(b[bw:], Symbols[SYM_TRIM])
		} else {
			bw += copy(b[bw:], Symbols[SYM_TRIM])
			bw += scanInterpolationMarker(f.dialect, b[bw:], *curArg)
			args[*curArg] = f.chars
			*curArg++
			bw += copy(b[bw:], Symbols[SYM_FROM])
		}
		bw += trimFuncScanSubject(f, b[bw:], args, curArg)
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
		size = (len(Symbols[SYM_TRIM]) + len(Symbols[SYM_LEADING]) +
			len(Symbols[SYM_FROM]))
		if f.chars != "" {
			// TRIM(LEADING chars FROM string)
			size += 1
		}
	case TRIM_TRAILING:
		// TRIM(TRAILING FROM string)
		size = (len(Symbols[SYM_TRIM]) + len(Symbols[SYM_TRAILING]) +
			len(Symbols[SYM_FROM]))
		if f.chars != "" {
			// TRIM(TRAILING chars FROM string)
			size += 1
		}
	case TRIM_BOTH:
		if f.chars == "" {
			// BTRIM(string)
			size = len(Symbols[SYM_BTRIM])
		} else {
			// BTRIM(string, chars)
			size = len(Symbols[SYM_BTRIM]) + len(Symbols[SYM_COMMA_WS])
		}
	}
	return size
}

// Helper function that scans into the supplied SQL []byte buffer for the
// TRIM/BTRIM() SQL function for PostgreSQL
// TRIM() function for PostgreSQL variants
func trimFuncScanPostgreSQL(f *trimFunc, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	switch f.location {
	case TRIM_LEADING:
		bw += copy(b[bw:], Symbols[SYM_TRIM])
		bw += copy(b[bw:], Symbols[SYM_LEADING])
		if f.chars != "" {
			bw += copy(b[bw:], []byte{' '})
			bw += scanInterpolationMarker(f.dialect, b[bw:], *curArg)
			args[*curArg] = f.chars
			*curArg++
		}
		bw += copy(b[bw:], Symbols[SYM_FROM])
		bw += trimFuncScanSubject(f, b[bw:], args, curArg)
	case TRIM_TRAILING:
		bw += copy(b[bw:], Symbols[SYM_TRIM])
		bw += copy(b[bw:], Symbols[SYM_TRAILING])
		if f.chars != "" {
			bw += copy(b[bw:], []byte{' '})
			bw += scanInterpolationMarker(f.dialect, b[bw:], *curArg)
			args[*curArg] = f.chars
			*curArg++
		}
		bw += copy(b[bw:], Symbols[SYM_FROM])
		bw += trimFuncScanSubject(f, b[bw:], args, curArg)
	case TRIM_BOTH:
		bw += copy(b[bw:], Symbols[SYM_BTRIM])
		bw += trimFuncScanSubject(f, b[bw:], args, curArg)
		if f.chars != "" {
			bw += copy(b[bw:], Symbols[SYM_COMMA_WS])
			bw += scanInterpolationMarker(f.dialect, b[bw:], *curArg)
			args[*curArg] = f.chars
			*curArg++
		}
	}
	return bw
}

// Scan in the subject of the TRIM() function
func trimFuncScanSubject(f *trimFunc, b []byte, args []interface{}, curArg *int) int {
	// We need to disable alias output for elements that are
	// projections. We don't want to output, for example,
	// "ON users.id AS user_id = TRIM(articles.author)"
	switch f.subject.(type) {
	case projection:
		reset := f.subject.(projection).disableAliasScan()
		defer reset()
	}
	return f.subject.scan(b, args, curArg)
}

func (f *trimFunc) size() int {
	size := 0
	switch f.dialect {
	case DIALECT_POSTGRESQL:
		size = trimFuncSizePostgreSQL(f)
	default:
		size = trimFuncSizeMySQL(f)
	}
	size += len(Symbols[SYM_RPAREN])
	// We need to disable alias output for elements that are
	// projections. We don't want to output, for example,
	// "ON users.id AS user_id = TRIM(articles.author)"
	switch f.subject.(type) {
	case projection:
		reset := f.subject.(projection).disableAliasScan()
		defer reset()
	}
	size += f.subject.size()
	if f.alias != "" {
		size += len(Symbols[SYM_AS]) + len(f.alias)
	}
	return size
}

func (f *trimFunc) scan(b []byte, args []interface{}, curArg *int) int {
	bw := 0
	switch f.dialect {
	case DIALECT_POSTGRESQL:
		bw += trimFuncScanPostgreSQL(f, b[bw:], args, curArg)
	default:
		bw += trimFuncScanMySQL(f, b[bw:], args, curArg)
	}
	bw += copy(b[bw:], Symbols[SYM_RPAREN])
	if f.alias != "" {
		bw += copy(b[bw:], Symbols[SYM_AS])
		bw += copy(b[bw:], f.alias)
	}
	return bw
}

// Returns a struct that will output the TRIM() SQL function, trimming leading
// and trailing whitespace from the supplied projection
func Trim(p projection) *trimFunc {
	return &trimFunc{
		subject:  p.(element),
		sel:      p.from(),
		location: TRIM_BOTH,
	}
}

func (c *Column) Trim() *trimFunc {
	f := Trim(c)
	f.setDialect(c.tbl.meta.dialect)
	return f
}

// Returns a struct that will output the LTRIM() SQL function for MySQL and the
// TRIM(LEADING FROM column) SQL function for PostgreSQL. The SQL function in
// either case will remove whitespace from the start of the supplied projection
func LTrim(p projection) *trimFunc {
	return &trimFunc{
		subject:  p.(element),
		sel:      p.from(),
		location: TRIM_LEADING,
	}
}

func (c *Column) LTrim() *trimFunc {
	f := LTrim(c)
	f.setDialect(c.tbl.meta.dialect)
	return f
}

// Returns a struct that will output the RTRIM() SQL function for MySQL and the
// TRIM(TRAILING FROM column) SQL function for PostgreSQL. The SQL function in
// either case will remove whitespace from the start of the supplied projection
func RTrim(p projection) *trimFunc {
	return &trimFunc{
		subject:  p.(element),
		sel:      p.from(),
		location: TRIM_TRAILING,
	}
}

func (c *Column) RTrim() *trimFunc {
	f := LTrim(c)
	f.setDialect(c.tbl.meta.dialect)
	return f
}

// Returns a struct that will output the TRIM() SQL function, trimming leading
// and trailing specified characters from the supplied projection
func TrimChars(p projection, chars string) *trimFunc {
	return &trimFunc{
		subject:  p.(element),
		sel:      p.from(),
		location: TRIM_BOTH,
		chars:    chars,
	}
}

func (c *Column) TrimChars(chars string) *trimFunc {
	f := TrimChars(c, chars)
	f.setDialect(c.tbl.meta.dialect)
	return f
}

// Returns a struct that will output the TRIM(LEADING chars FROM column) SQL
// function, trimming leading specified characters from the supplied projection
func LTrimChars(p projection, chars string) *trimFunc {
	return &trimFunc{
		subject:  p.(element),
		sel:      p.from(),
		location: TRIM_LEADING,
		chars:    chars,
	}
}

func (c *Column) LTrimChars(chars string) *trimFunc {
	f := LTrimChars(c, chars)
	f.setDialect(c.tbl.meta.dialect)
	return f
}

// Returns a struct that will output the TRIM(TRAILING chars FROM column) SQL
// function, trimming trailing specified characters from the supplied
// projection
func RTrimChars(p projection, chars string) *trimFunc {
	return &trimFunc{
		subject:  p.(element),
		sel:      p.from(),
		location: TRIM_TRAILING,
		chars:    chars,
	}
}

func (c *Column) RTrimChars(chars string) *trimFunc {
	f := RTrimChars(c, chars)
	f.setDialect(c.tbl.meta.dialect)
	return f
}
