//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/types"
)

type sqlFunc struct {
	sel      types.Selection
	alias    string
	ScanInfo grammar.ScanInfo
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
	for _, sym := range f.ScanInfo {
		switch sym {
		case grammar.SYM_ELEMENT:
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
			size += len(grammar.Symbols[sym])
		}
	}
	if f.alias != "" {
		size += len(grammar.Symbols[grammar.SYM_AS]) + len(f.alias)
	}
	return size
}

func (f *sqlFunc) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	elidx := 0
	for _, sym := range f.ScanInfo {
		if sym == grammar.SYM_ELEMENT {
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
			bw += copy(b[bw:], grammar.Symbols[sym])
		}
	}
	if f.alias != "" {
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_AS])
		bw += copy(b[bw:], f.alias)
	}
	return bw
}

func Max(p types.Projection) *sqlFunc {
	return &sqlFunc{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_MAX),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func (c *ColumnIdentifier) Max() *sqlFunc {
	return Max(c)
}

func Min(p types.Projection) *sqlFunc {
	return &sqlFunc{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_MIN),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func (c *ColumnIdentifier) Min() *sqlFunc {
	return Min(c)
}

func Sum(p types.Projection) *sqlFunc {
	return &sqlFunc{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_SUM),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func (c *ColumnIdentifier) Sum() *sqlFunc {
	return Sum(c)
}

func Avg(p types.Projection) *sqlFunc {
	return &sqlFunc{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_AVG),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func (c *ColumnIdentifier) Avg() *sqlFunc {
	return Avg(c)
}

func Count(sel types.Selection) *sqlFunc {
	return &sqlFunc{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_COUNT_STAR),
		sel:      sel,
	}
}

func CountDistinct(p types.Projection) *sqlFunc {
	return &sqlFunc{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_COUNT_DISTINCT),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func Cast(p types.Projection, stype grammar.SqlType) *sqlFunc {
	si := make([]grammar.Symbol, len(grammar.FunctionScanTable(grammar.FUNC_CAST)))
	copy(si, grammar.FunctionScanTable(grammar.FUNC_CAST))
	// Replace the placeholder with the SQL type's appropriate []byte
	// representation
	si[3] = grammar.SQLTypeToSymbol(stype)
	return &sqlFunc{
		ScanInfo: si,
		elements: []types.Element{p.(types.Element)},
	}
}

func CharLength(p types.Projection) *sqlFunc {
	return &sqlFunc{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CHAR_LENGTH),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func (c *ColumnIdentifier) CharLength() *sqlFunc {
	return CharLength(c)
}

func BitLength(p types.Projection) *sqlFunc {
	return &sqlFunc{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_BIT_LENGTH),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func (c *ColumnIdentifier) BitLength() *sqlFunc {
	return BitLength(c)
}

func Ascii(p types.Projection) *sqlFunc {
	return &sqlFunc{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_ASCII),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func (c *ColumnIdentifier) Ascii() *sqlFunc {
	return Ascii(c)
}

func Reverse(p types.Projection) *sqlFunc {
	return &sqlFunc{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_REVERSE),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func (c *ColumnIdentifier) Reverse() *sqlFunc {
	return Reverse(c)
}

func Concat(projs ...types.Projection) *sqlFunc {
	els := make([]types.Element, len(projs))
	for x, p := range projs {
		els[x] = p.(types.Element)
	}
	subjects := &List{elements: els}
	return &sqlFunc{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CONCAT),
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
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CONCAT_WS),
		elements: []types.Element{ast.NewValue(nil, sep), subjects},
		// TODO(jaypipes): Clearly we need to support >1 selection...
		sel: projs[0].From(),
	}
}

func Now() *sqlFunc {
	return &sqlFunc{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_NOW),
	}
}

func CurrentTimestamp() *sqlFunc {
	return &sqlFunc{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_TIMESTAMP),
	}
}

func CurrentTime() *sqlFunc {
	return &sqlFunc{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_TIME),
	}
}

func CurrentDate() *sqlFunc {
	return &sqlFunc{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_DATE),
	}
}

func Extract(p types.Projection, unit grammar.IntervalUnit) *sqlFunc {
	si := make([]grammar.Symbol, len(grammar.FunctionScanTable(grammar.FUNC_EXTRACT)))
	copy(si, grammar.FunctionScanTable(grammar.FUNC_EXTRACT))
	// Replace the placeholder with the interval unit's appropriate []byte
	// representation
	si[1] = grammar.IntervalUnitToSymbol(unit)
	return &sqlFunc{
		ScanInfo: si,
		elements: []types.Element{p.(types.Element)},
	}
}
