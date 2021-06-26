//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import "github.com/jaypipes/sqlb/pkg/types"

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
			SYM_EXTRACT, SYM_PLACEHOLDER, SYM_SPACE, SYM_FROM, SYM_ELEMENT, SYM_RPAREN,
		},
	}
)

type sqlFunc struct {
	sel      types.Selection
	alias    string
	scanInfo scanInfo
	elements []types.Element
}

func (f *sqlFunc) From() types.Selection {
	return f.sel
}

func (f *sqlFunc) DisableAliasScan() func() {
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

func (e *sqlFunc) ArgCount() int {
	ac := 0
	for _, el := range e.elements {
		ac += el.ArgCount()
	}
	return ac
}

func (f *sqlFunc) Size(scanner types.Scanner) int {
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
			case types.Projection:
				reset := el.(types.Projection).DisableAliasScan()
				defer reset()
			}
			elidx++
			size += el.Size(scanner)
		default:
			size += len(Symbols[sym])
		}
	}
	if f.alias != "" {
		size += len(Symbols[SYM_AS]) + len(f.alias)
	}
	return size
}

func (f *sqlFunc) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	elidx := 0
	for _, sym := range f.scanInfo {
		if sym == SYM_ELEMENT {
			el := f.elements[elidx]
			// We need to disable alias output for elements that are
			// projections. We don't want to output, for example,
			// "ON users.id AS user_id = articles.author"
			switch el.(type) {
			case types.Projection:
				reset := el.(types.Projection).DisableAliasScan()
				defer reset()
			}
			elidx++
			bw += el.Scan(scanner, b[bw:], args, curArg)
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

func Max(p types.Projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_MAX],
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func (c *Column) Max() *sqlFunc {
	return Max(c)
}

func Min(p types.Projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_MIN],
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func (c *Column) Min() *sqlFunc {
	return Min(c)
}

func Sum(p types.Projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_SUM],
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func (c *Column) Sum() *sqlFunc {
	return Sum(c)
}

func Avg(p types.Projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_AVG],
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func (c *Column) Avg() *sqlFunc {
	return Avg(c)
}

func Count(sel types.Selection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_COUNT_STAR],
		sel:      sel,
	}
}

func CountDistinct(p types.Projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_COUNT_DISTINCT],
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func Cast(p types.Projection, stype SqlType) *sqlFunc {
	si := make([]Symbol, len(funcScanTable[FUNC_CAST]))
	copy(si, funcScanTable[FUNC_CAST])
	// Replace the placeholder with the SQL type's appropriate []byte
	// representation
	si[3] = sqlTypeToSymbol[stype]
	return &sqlFunc{
		scanInfo: si,
		elements: []types.Element{p.(types.Element)},
	}
}

func CharLength(p types.Projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_CHAR_LENGTH],
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func (c *Column) CharLength() *sqlFunc {
	return CharLength(c)
}

func BitLength(p types.Projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_BIT_LENGTH],
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func (c *Column) BitLength() *sqlFunc {
	return BitLength(c)
}

func Ascii(p types.Projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_ASCII],
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func (c *Column) Ascii() *sqlFunc {
	return Ascii(c)
}

func Reverse(p types.Projection) *sqlFunc {
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_REVERSE],
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func (c *Column) Reverse() *sqlFunc {
	return Reverse(c)
}

func Concat(projs ...types.Projection) *sqlFunc {
	els := make([]types.Element, len(projs))
	for x, p := range projs {
		els[x] = p.(types.Element)
	}
	subjects := &List{elements: els}
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_CONCAT],
		elements: []types.Element{subjects},
		// TODO(jaypipes): Clearly we need to support >1 selection...
		sel: projs[0].From(),
	}
}

func ConcatWs(sep string, projs ...types.Projection) *sqlFunc {
	els := make([]types.Element, len(projs))
	for x, p := range projs {
		els[x] = p.(types.Element)
	}
	subjects := &List{elements: els}
	return &sqlFunc{
		scanInfo: funcScanTable[FUNC_CONCAT_WS],
		elements: []types.Element{&value{val: sep}, subjects},
		// TODO(jaypipes): Clearly we need to support >1 selection...
		sel: projs[0].From(),
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

func Extract(p types.Projection, unit IntervalUnit) *sqlFunc {
	si := make([]Symbol, len(funcScanTable[FUNC_EXTRACT]))
	copy(si, funcScanTable[FUNC_EXTRACT])
	// Replace the placeholder with the interval unit's appropriate []byte
	// representation
	si[1] = intervalUnitToSymbol[unit]
	return &sqlFunc{
		scanInfo: si,
		elements: []types.Element{p.(types.Element)},
	}
}
