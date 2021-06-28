//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ast

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/types"
)

// Function is a SQL function that accepts zero or more parameters
type Function struct {
	sel      types.Selection
	Alias    string
	ScanInfo grammar.ScanInfo
	elements []types.Element
}

func (f *Function) From() types.Selection {
	return f.sel
}

func (f *Function) DisableAliasScan() func() {
	origAlias := f.Alias
	f.Alias = ""
	return func() { f.Alias = origAlias }
}

func (f *Function) As(alias string) *Function {
	return &Function{
		sel:      f.sel,
		Alias:    alias,
		ScanInfo: f.ScanInfo,
		elements: f.elements,
	}
}

func (e *Function) ArgCount() int {
	ac := 0
	for _, el := range e.elements {
		ac += el.ArgCount()
	}
	return ac
}

func (f *Function) Size(scanner types.Scanner) int {
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
	if f.Alias != "" {
		size += len(grammar.Symbols[grammar.SYM_AS]) + len(f.Alias)
	}
	return size
}

func (f *Function) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
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
	if f.Alias != "" {
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_AS])
		bw += copy(b[bw:], f.Alias)
	}
	return bw
}

func (f *Function) Desc() *SortColumn {
	return &SortColumn{p: f, desc: true}
}

func (f *Function) Asc() *SortColumn {
	return &SortColumn{p: f}
}

func Max(p types.Projection) *Function {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_MAX),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func Min(p types.Projection) *Function {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_MIN),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func Sum(p types.Projection) *Function {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_SUM),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func Avg(p types.Projection) *Function {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_AVG),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func Count(sel types.Selection) *Function {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_COUNT_STAR),
		sel:      sel,
	}
}

func CountDistinct(p types.Projection) *Function {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_COUNT_DISTINCT),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func Cast(p types.Projection, stype grammar.SqlType) *Function {
	si := make([]grammar.Symbol, len(grammar.FunctionScanTable(grammar.FUNC_CAST)))
	copy(si, grammar.FunctionScanTable(grammar.FUNC_CAST))
	// Replace the placeholder with the SQL type's appropriate []byte
	// representation
	si[3] = grammar.SQLTypeToSymbol(stype)
	return &Function{
		ScanInfo: si,
		elements: []types.Element{p.(types.Element)},
	}
}

func CharLength(p types.Projection) *Function {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CHAR_LENGTH),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func BitLength(p types.Projection) *Function {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_BIT_LENGTH),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func Ascii(p types.Projection) *Function {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_ASCII),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func Reverse(p types.Projection) *Function {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_REVERSE),
		elements: []types.Element{p.(types.Element)},
		sel:      p.From(),
	}
}

func Concat(projs ...types.Projection) *Function {
	els := make([]types.Element, len(projs))
	for x, p := range projs {
		els[x] = p.(types.Element)
	}
	subjects := NewList(els...)
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CONCAT),
		elements: []types.Element{subjects},
		// TODO(jaypipes): Clearly we need to support >1 selection...
		sel: projs[0].From(),
	}
}

func ConcatWs(sep string, projs ...types.Projection) *Function {
	els := make([]types.Element, len(projs))
	for x, p := range projs {
		els[x] = p.(types.Element)
	}
	subjects := NewList(els...)
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CONCAT_WS),
		elements: []types.Element{NewValue(nil, sep), subjects},
		// TODO(jaypipes): Clearly we need to support >1 selection...
		sel: projs[0].From(),
	}
}

func Now() *Function {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_NOW),
	}
}

func CurrentTimestamp() *Function {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_TIMESTAMP),
	}
}

func CurrentTime() *Function {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_TIME),
	}
}

func CurrentDate() *Function {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_DATE),
	}
}

func Extract(p types.Projection, unit grammar.IntervalUnit) *Function {
	si := make([]grammar.Symbol, len(grammar.FunctionScanTable(grammar.FUNC_EXTRACT)))
	copy(si, grammar.FunctionScanTable(grammar.FUNC_EXTRACT))
	// Replace the placeholder with the interval unit's appropriate []byte
	// representation
	si[1] = grammar.IntervalUnitToSymbol(unit)
	return &Function{
		ScanInfo: si,
		elements: []types.Element{p.(types.Element)},
	}
}
