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
//		TRIM(chars FROM string)
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
	chars    []byte
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
	return f.subject.argCount()
}

// Helper function that returns the non-subject, non-interpolation size of the
// TRIM() function for MySQL variants
func trimFuncSizeMySQL(f *trimFunc) int {
	size := 0
	switch f.location {
	case TRIM_LEADING:
		if f.chars == nil {
			// LTRIM(string)
			size = len(Symbols[SYM_LTRIM])
		} else {
			// TRIM(LEADING remstr FROM string)
			size = (len(Symbols[SYM_TRIM]) + len(Symbols[SYM_LEADING]) +
				len(Symbols[SYM_FROM]) + 1)
		}
	case TRIM_TRAILING:
		if f.chars == nil {
			// LTRIM(string)
			size = len(Symbols[SYM_RTRIM])
		} else {
			// TRIM(TRAILING remstr FROM string)
			size = (len(Symbols[SYM_TRIM]) + len(Symbols[SYM_TRAILING]) +
				len(Symbols[SYM_FROM]) + 1)
		}
	case TRIM_BOTH:
		if f.chars == nil {
			// TRIM(string)
			size = len(Symbols[SYM_TRIM])
		} else {
			// TRIM(remstr FROM string)
			size = len(Symbols[SYM_TRIM]) + len(Symbols[SYM_FROM]) + 1
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
		if f.chars == nil {
			bw += copy(b[bw:], Symbols[SYM_LTRIM])
		} else {
			bw += copy(b[bw:], Symbols[SYM_TRIM])
			bw += copy(b[bw:], Symbols[SYM_LEADING])
			args[*curArg] = string(f.chars)
			*curArg++
			bw += copy(b[bw:], []byte{' '})
			bw += copy(b[bw:], Symbols[SYM_FROM])
		}
	case TRIM_TRAILING:
		if f.chars == nil {
			bw += copy(b[bw:], Symbols[SYM_RTRIM])
		} else {
			bw += copy(b[bw:], Symbols[SYM_TRIM])
			bw += copy(b[bw:], Symbols[SYM_TRAILING])
			args[*curArg] = string(f.chars)
			*curArg++
			bw += copy(b[bw:], []byte{' '})
			bw += copy(b[bw:], Symbols[SYM_FROM])
		}
	case TRIM_BOTH:
		if f.chars == nil {
			bw += copy(b[bw:], Symbols[SYM_TRIM])
		} else {
			bw += copy(b[bw:], Symbols[SYM_TRIM])
			bw += copy(b[bw:], Symbols[SYM_BOTH])
			args[*curArg] = string(f.chars)
			*curArg++
			bw += copy(b[bw:], []byte{' '})
			bw += copy(b[bw:], Symbols[SYM_FROM])
		}
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
		if f.chars != nil {
			// TRIM(LEADING chars FROM string)
			size += 1
		}
	case TRIM_TRAILING:
		// TRIM(TRAILING FROM string)
		size = (len(Symbols[SYM_TRIM]) + len(Symbols[SYM_TRAILING]) +
			len(Symbols[SYM_FROM]))
		if f.chars != nil {
			// TRIM(TRAILING chars FROM string)
			size += 1
		}
	case TRIM_BOTH:
		if f.chars == nil {
			// BTRIM(string)
			size = len(Symbols[SYM_BTRIM])
		} else {
			// TRIM(BOTH chars FROM string)
			size = (len(Symbols[SYM_TRIM]) + len(Symbols[SYM_BOTH]) +
				len(Symbols[SYM_FROM]) + 1)
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
		if f.chars != nil {
			args[*curArg] = string(f.chars)
			*curArg++
			bw += copy(b[bw:], []byte{' '})
		}
		bw += copy(b[bw:], Symbols[SYM_FROM])
	case TRIM_TRAILING:
		bw += copy(b[bw:], Symbols[SYM_TRIM])
		bw += copy(b[bw:], Symbols[SYM_TRAILING])
		if f.chars != nil {
			args[*curArg] = string(f.chars)
			*curArg++
			bw += copy(b[bw:], []byte{' '})
		}
		bw += copy(b[bw:], Symbols[SYM_FROM])
	case TRIM_BOTH:
		if f.chars == nil {
			bw += copy(b[bw:], Symbols[SYM_BTRIM])
		} else {
			bw += copy(b[bw:], Symbols[SYM_TRIM])
			bw += copy(b[bw:], Symbols[SYM_BOTH])
			args[*curArg] = string(f.chars)
			*curArg++
			bw += copy(b[bw:], []byte{' '})
			bw += copy(b[bw:], Symbols[SYM_FROM])
		}
	}
	return bw
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
	// "ON users.id AS user_id = articles.author"
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
	// We need to disable alias output for elements that are
	// projections. We don't want to output, for example,
	// "ON users.id AS user_id = articles.author"
	switch f.subject.(type) {
	case projection:
		reset := f.subject.(projection).disableAliasScan()
		defer reset()
	}
	bw += f.subject.scan(b[bw:], args, curArg)
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
