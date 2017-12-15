package sqlb

// INSERT INTO <table> (<columns>) VALUES (<values>)

type insertStatement struct {
	table   *Table
	columns []*Column
	values  []interface{}
}

func (s *insertStatement) argCount() int {
	return len(s.values)
}

func (s *insertStatement) size(scanner *sqlScanner) int {
	size := len(Symbols[SYM_INSERT]) + len(s.table.name) + 1 // space after table name
	ncols := len(s.columns)
	for _, c := range s.columns {
		// We don't add the table identifier or use an alias when outputting
		// the column names in the <columns> element of the INSERT statement
		size += len(c.name)
	}
	// We don't include interpolation marks in our sizing, since the length
	// differs with SQL dialects. This is accounted for by callers of scan().
	size += len(Symbols[SYM_LPAREN]) + len(Symbols[SYM_VALUES])
	// Two comma-delimited lists of same number of elements (columns and
	// values)
	size += 2 * (len(Symbols[SYM_COMMA_WS]) * (ncols - 1)) // the commas...
	size += len(Symbols[SYM_RPAREN])
	return size
}

func (s *insertStatement) scan(scanner *sqlScanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], Symbols[SYM_INSERT])
	// We don't add any table alias when outputting the table identifier
	bw += copy(b[bw:], s.table.name)
	bw += copy(b[bw:], " ")
	bw += copy(b[bw:], Symbols[SYM_LPAREN])

	ncols := len(s.columns)
	for x, c := range s.columns {
		// We don't add the table identifier or use an alias when outputting
		// the column names in the <columns> element of the INSERT statement
		bw += copy(b[bw:], c.name)
		if x != (ncols - 1) {
			bw += copy(b[bw:], Symbols[SYM_COMMA_WS])
		}
	}
	bw += copy(b[bw:], Symbols[SYM_VALUES])
	for x, v := range s.values {
		bw += scanInterpolationMarker(scanner.dialect, b[bw:], *curArg)
		args[*curArg] = v
		*curArg++
		if x != (ncols - 1) {
			bw += copy(b[bw:], Symbols[SYM_COMMA_WS])
		}
	}
	bw += copy(b[bw:], Symbols[SYM_RPAREN])
	return bw
}
