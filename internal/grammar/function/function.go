//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package function

import (
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/element"
	"github.com/jaypipes/sqlb/internal/grammar/sortcolumn"
)

// Function is a SQL function that accepts zero or more parameters
type Function struct {
	sel      builder.Selection
	Alias    string
	ScanInfo grammar.ScanInfo
	elements []builder.Element
}

func (f *Function) From() builder.Selection {
	return f.sel
}

func (f *Function) DisableAliasScan() func() {
	origAlias := f.Alias
	f.Alias = ""
	return func() { f.Alias = origAlias }
}

func (f *Function) As(alias string) builder.Projection {
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

func (f *Function) Size(b *builder.Builder) int {
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
			case builder.Projection:
				reset := el.(builder.Projection).DisableAliasScan()
				defer reset()
			}
			elidx++
			size += el.Size(b)
		default:
			size += len(grammar.Symbols[sym])
		}
	}
	if f.Alias != "" {
		size += len(grammar.Symbols[grammar.SYM_AS]) + len(f.Alias)
	}
	return size
}

func (f *Function) Scan(b *builder.Builder, args []interface{}, curArg *int) {
	elidx := 0
	for _, sym := range f.ScanInfo {
		if sym == grammar.SYM_ELEMENT {
			el := f.elements[elidx]
			// We need to disable alias output for elements that are
			// projections. We don't want to output, for example,
			// "ON users.id AS user_id = articles.author"
			switch el.(type) {
			case builder.Projection:
				reset := el.(builder.Projection).DisableAliasScan()
				defer reset()
			}
			elidx++
			el.Scan(b, args, curArg)
		} else {
			b.Write(grammar.Symbols[sym])
		}
	}
	if f.Alias != "" {
		b.Write(grammar.Symbols[grammar.SYM_AS])
		b.WriteString(f.Alias)
	}
}

func (f *Function) Desc() builder.Sortable {
	return sortcolumn.NewDesc(f)
}

func (f *Function) Asc() builder.Sortable {
	return sortcolumn.NewAsc(f)
}

func Max(p builder.Projection) builder.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_MAX),
		elements: []builder.Element{p.(builder.Element)},
		sel:      p.From(),
	}
}

func Min(p builder.Projection) builder.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_MIN),
		elements: []builder.Element{p.(builder.Element)},
		sel:      p.From(),
	}
}

func Sum(p builder.Projection) builder.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_SUM),
		elements: []builder.Element{p.(builder.Element)},
		sel:      p.From(),
	}
}

func Avg(p builder.Projection) builder.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_AVG),
		elements: []builder.Element{p.(builder.Element)},
		sel:      p.From(),
	}
}

func Count(sel builder.Selection) builder.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_COUNT_STAR),
		sel:      sel,
	}
}

func CountDistinct(p builder.Projection) builder.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_COUNT_DISTINCT),
		elements: []builder.Element{p.(builder.Element)},
		sel:      p.From(),
	}
}

func Cast(p builder.Projection, stype grammar.SQLType) builder.Projection {
	si := make([]grammar.Symbol, len(grammar.FunctionScanTable(grammar.FUNC_CAST)))
	copy(si, grammar.FunctionScanTable(grammar.FUNC_CAST))
	// Replace the placeholder with the SQL type's appropriate []byte
	// representation
	si[3] = grammar.SQLTypeToSymbol(stype)
	return &Function{
		ScanInfo: si,
		elements: []builder.Element{p.(builder.Element)},
	}
}

func CharLength(p builder.Projection) builder.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CHAR_LENGTH),
		elements: []builder.Element{p.(builder.Element)},
		sel:      p.From(),
	}
}

func BitLength(p builder.Projection) builder.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_BIT_LENGTH),
		elements: []builder.Element{p.(builder.Element)},
		sel:      p.From(),
	}
}

func Ascii(p builder.Projection) builder.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_ASCII),
		elements: []builder.Element{p.(builder.Element)},
		sel:      p.From(),
	}
}

func Reverse(p builder.Projection) builder.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_REVERSE),
		elements: []builder.Element{p.(builder.Element)},
		sel:      p.From(),
	}
}

func Concat(projs ...builder.Projection) builder.Projection {
	els := make([]builder.Element, len(projs))
	for x, p := range projs {
		els[x] = p.(builder.Element)
	}
	subjects := element.NewList(els...)
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CONCAT),
		elements: []builder.Element{subjects},
		// TODO(jaypipes): Clearly we need to support >1 selection...
		sel: projs[0].From(),
	}
}

func ConcatWs(sep string, projs ...builder.Projection) builder.Projection {
	els := make([]builder.Element, len(projs))
	for x, p := range projs {
		els[x] = p.(builder.Element)
	}
	subjects := element.NewList(els...)
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CONCAT_WS),
		elements: []builder.Element{element.NewValue(nil, sep), subjects},
		// TODO(jaypipes): Clearly we need to support >1 selection...
		sel: projs[0].From(),
	}
}

func Now() builder.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_NOW),
	}
}

func CurrentTimestamp() builder.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_TIMESTAMP),
	}
}

func CurrentTime() builder.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_TIME),
	}
}

func CurrentDate() builder.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_DATE),
	}
}

func Extract(p builder.Projection, unit grammar.IntervalUnit) builder.Projection {
	si := make([]grammar.Symbol, len(grammar.FunctionScanTable(grammar.FUNC_EXTRACT)))
	copy(si, grammar.FunctionScanTable(grammar.FUNC_EXTRACT))
	// Replace the placeholder with the interval unit's appropriate []byte
	// representation
	si[1] = grammar.IntervalUnitToSymbol(unit)
	return &Function{
		ScanInfo: si,
		elements: []builder.Element{p.(builder.Element)},
	}
}
