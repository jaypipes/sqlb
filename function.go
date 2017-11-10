package sqlb

type funcId int

const (
	FUNC_MAX funcId = iota
	FUNC_MIN
	FUNC_SUM
	FUNC_AVG
	FUNC_COUNT_STAR
	FUNC_COUNT_DISTINCT
	FUNC_CAST
	FUNC_CHAR_LENGTH
	FUNC_BIT_LENGTH
	FUNC_ASCII
	FUNC_REVERSE
	FUNC_CONCAT
	FUNC_CONCAT_WS
	FUNC_NOW
	FUNC_CURRENT_TIMESTAMP
	FUNC_CURRENT_TIME
	FUNC_CURRENT_DATE
	FUNC_EXTRACT
)

var (
	// A static table containing information used in constructing the
	// expression's SQL string during scan() calls
	funcScanTable = map[funcId]scanInfo{
		FUNC_MAX: scanInfo{
			SYM_MAX, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_MIN: scanInfo{
			SYM_MIN, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_SUM: scanInfo{
			SYM_SUM, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_AVG: scanInfo{
			SYM_AVG, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_COUNT_STAR: scanInfo{
			SYM_COUNT_STAR,
		},
		FUNC_COUNT_DISTINCT: scanInfo{
			SYM_COUNT_DISTINCT, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_CAST: scanInfo{
			SYM_CAST, SYM_ELEMENT, SYM_AS, SYM_PLACEHOLDER, SYM_RPAREN,
		},
		FUNC_CHAR_LENGTH: scanInfo{
			SYM_CHAR_LENGTH, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_BIT_LENGTH: scanInfo{
			SYM_BIT_LENGTH, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_ASCII: scanInfo{
			SYM_ASCII, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_REVERSE: scanInfo{
			SYM_REVERSE, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_CONCAT: scanInfo{
			SYM_CONCAT, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_CONCAT_WS: scanInfo{
			SYM_CONCAT_WS, SYM_ELEMENT, SYM_COMMA_WS, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_NOW: scanInfo{
			SYM_NOW,
		},
		FUNC_CURRENT_TIMESTAMP: scanInfo{
			SYM_CURRENT_TIMESTAMP,
		},
		FUNC_CURRENT_TIME: scanInfo{
			SYM_CURRENT_TIME,
		},
		FUNC_CURRENT_DATE: scanInfo{
			SYM_CURRENT_DATE,
		},
		// This is the MySQL variant of EXTRACT, which follows the form
		// EXTRACT(field FROM source). PostgreSQL has a different format for
		// EXTRACT() which follows the following format:
		// EXTRACT(field FROM [interval|timestamp] source)
		FUNC_EXTRACT: scanInfo{
			SYM_EXTRACT, SYM_PLACEHOLDER, SYM_FROM, SYM_ELEMENT, SYM_RPAREN,
		},
	}
)

type sqlFunc struct {
	sel      selection
	alias    string
	scanInfo scanInfo
	elements []element
	dialect  Dialect
}

// Sets the sqlFunc's dialect and pushes the dialect down into any of the
// sqlFunc's elements
func (f *sqlFunc) setDialect(dialect Dialect) {
	f.dialect = dialect
	for _, el := range f.elements {
		switch el.(type) {
		case *value:
			v := el.(*value)
			v.dialect = dialect
		}
	}
}

func (f *sqlFunc) from() selection {
	return f.sel
}

func (f *sqlFunc) disableAliasScan() func() {
	origAlias := f.alias
	f.alias = ""
	return func() { f.alias = origAlias }
}

func (f *sqlFunc) Alias(alias string) {
	f.alias = alias
}

func (f *sqlFunc) As(alias string) *sqlFunc {
	f.Alias(alias)
	return f
}

func (e *sqlFunc) argCount() int {
	ac := 0
	for _, el := range e.elements {
		ac += el.argCount()
	}
	return ac
}

func (f *sqlFunc) size() int {
	size := 0
	elidx := 0
	for _, sym := range f.scanInfo {
		switch sym {
		case SYM_ELEMENT:
			el := f.elements[elidx]
			// We need to disable alias output for elements that are
			// projections. We don't want to output, for example,
			// "ON users.id AS user_id = articles.author"
			switch el.(type) {
			case projection:
				reset := el.(projection).disableAliasScan()
				defer reset()
			}
			elidx++
			size += el.size()
		default:
			size += len(Symbols[sym])
		}
	}
	if f.alias != "" {
		size += len(Symbols[SYM_AS]) + len(f.alias)
	}
	return size
}

func (f *sqlFunc) scan(b []byte, args []interface{}, curArg *int) int {
	bw := 0
	elidx := 0
	for _, sym := range f.scanInfo {
		if sym == SYM_ELEMENT {
			el := f.elements[elidx]
			// We need to disable alias output for elements that are
			// projections. We don't want to output, for example,
			// "ON users.id AS user_id = articles.author"
			switch el.(type) {
			case projection:
				reset := el.(projection).disableAliasScan()
				defer reset()
			}
			elidx++
			bw += el.scan(b[bw:], args, curArg)
		} else {
			bw += copy(b[bw:], Symbols[sym])
		}
	}
	if f.alias != "" {
		bw += copy(b[bw:], Symbols[SYM_AS])
		bw += copy(b[bw:], f.alias)
	}
	return bw
}

func Max(p projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_MAX],
		elements: []element{p.(element)},
		sel:      p.from(),
	}
}

func (c *Column) Max() *sqlFunc {
	return Max(c)
}

func Min(p projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_MIN],
		elements: []element{p.(element)},
		sel:      p.from(),
	}
}

func (c *Column) Min() *sqlFunc {
	return Min(c)
}

func Sum(p projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_SUM],
		elements: []element{p.(element)},
		sel:      p.from(),
	}
}

func (c *Column) Sum() *sqlFunc {
	return Sum(c)
}

func Avg(p projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_AVG],
		elements: []element{p.(element)},
		sel:      p.from(),
	}
}

func (c *Column) Avg() *sqlFunc {
	return Avg(c)
}

func Count(sel selection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_COUNT_STAR],
		sel:      sel,
	}
}

func CountDistinct(p projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_COUNT_DISTINCT],
		elements: []element{p.(element)},
		sel:      p.from(),
	}
}

func Cast(p projection, stype SqlType) *sqlFunc {
	si := make([]Symbol, len(funcScanTable[FUNC_CAST]))
	copy(si, funcScanTable[FUNC_CAST])
	// Replace the placeholder with the SQL type's appropriate []byte
	// representation
	si[3] = sqlTypeToSymbol[stype]
	return &sqlFunc{
		scanInfo: si,
		elements: []element{p.(element)},
	}
}

func CharLength(p projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_CHAR_LENGTH],
		elements: []element{p.(element)},
		sel:      p.from(),
	}
}

func (c *Column) CharLength() *sqlFunc {
	return CharLength(c)
}

func BitLength(p projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_BIT_LENGTH],
		elements: []element{p.(element)},
		sel:      p.from(),
	}
}

func (c *Column) BitLength() *sqlFunc {
	return BitLength(c)
}

func Ascii(p projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_ASCII],
		elements: []element{p.(element)},
		sel:      p.from(),
	}
}

func (c *Column) Ascii() *sqlFunc {
	return Ascii(c)
}

func Reverse(p projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_REVERSE],
		elements: []element{p.(element)},
		sel:      p.from(),
	}
}

func (c *Column) Reverse() *sqlFunc {
	return Reverse(c)
}

func Concat(projs ...projection) *sqlFunc {
	els := make([]element, len(projs))
	for x, p := range projs {
		els[x] = p.(element)
	}
	subjects := &List{elements: els}
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_CONCAT],
		elements: []element{subjects},
		// TODO(jaypipes): Clearly we need to support >1 selection...
		sel: projs[0].from(),
	}
}

func ConcatWs(sep string, projs ...projection) *sqlFunc {
	els := make([]element, len(projs))
	for x, p := range projs {
		els[x] = p.(element)
	}
	subjects := &List{elements: els}
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_CONCAT_WS],
		elements: []element{&value{val: sep}, subjects},
		// TODO(jaypipes): Clearly we need to support >1 selection...
		sel: projs[0].from(),
	}
}

func Now() *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_NOW],
	}
}

func CurrentTimestamp() *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_CURRENT_TIMESTAMP],
	}
}

func CurrentTime() *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_CURRENT_TIME],
	}
}

func CurrentDate() *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_CURRENT_DATE],
	}
}

func Extract(p projection, unit IntervalUnit) *sqlFunc {
	si := make([]Symbol, len(funcScanTable[FUNC_EXTRACT]))
	copy(si, funcScanTable[FUNC_EXTRACT])
	// Replace the placeholder with the interval unit's appropriate []byte
	// representation
	si[1] = intervalUnitToSymbol[unit]
	return &sqlFunc{
		scanInfo: si,
		elements: []element{p.(element)},
	}
}
