//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package function

import (
	"strings"

	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/element"
	"github.com/jaypipes/sqlb/internal/grammar/sortcolumn"
	"github.com/jaypipes/sqlb/internal/scanner"
)

// Function is a SQL function that accepts zero or more parameters
type Function struct {
	sel      scanner.Selection
	Alias    string
	ScanInfo grammar.ScanInfo
	elements []scanner.Element
}

func (f *Function) From() scanner.Selection {
	return f.sel
}

func (f *Function) DisableAliasScan() func() {
	origAlias := f.Alias
	f.Alias = ""
	return func() { f.Alias = origAlias }
}

func (f *Function) As(alias string) scanner.Projection {
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

func (f *Function) Size(s *scanner.Scanner) int {
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
			case scanner.Projection:
				reset := el.(scanner.Projection).DisableAliasScan()
				defer reset()
			}
			elidx++
			size += el.Size(s)
		default:
			size += len(grammar.Symbols[sym])
		}
	}
	if f.Alias != "" {
		size += len(grammar.Symbols[grammar.SYM_AS]) + len(f.Alias)
	}
	return size
}

func (f *Function) Scan(s *scanner.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	elidx := 0
	for _, sym := range f.ScanInfo {
		if sym == grammar.SYM_ELEMENT {
			el := f.elements[elidx]
			// We need to disable alias output for elements that are
			// projections. We don't want to output, for example,
			// "ON users.id AS user_id = articles.author"
			switch el.(type) {
			case scanner.Projection:
				reset := el.(scanner.Projection).DisableAliasScan()
				defer reset()
			}
			elidx++
			el.Scan(s, b, args, curArg)
		} else {
			b.Write(grammar.Symbols[sym])
		}
	}
	if f.Alias != "" {
		b.Write(grammar.Symbols[grammar.SYM_AS])
		b.WriteString(f.Alias)
	}
}

func (f *Function) Desc() scanner.Sortable {
	return sortcolumn.NewDesc(f)
}

func (f *Function) Asc() scanner.Sortable {
	return sortcolumn.NewAsc(f)
}

func Max(p scanner.Projection) scanner.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_MAX),
		elements: []scanner.Element{p.(scanner.Element)},
		sel:      p.From(),
	}
}

func Min(p scanner.Projection) scanner.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_MIN),
		elements: []scanner.Element{p.(scanner.Element)},
		sel:      p.From(),
	}
}

func Sum(p scanner.Projection) scanner.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_SUM),
		elements: []scanner.Element{p.(scanner.Element)},
		sel:      p.From(),
	}
}

func Avg(p scanner.Projection) scanner.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_AVG),
		elements: []scanner.Element{p.(scanner.Element)},
		sel:      p.From(),
	}
}

func Count(sel scanner.Selection) scanner.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_COUNT_STAR),
		sel:      sel,
	}
}

func CountDistinct(p scanner.Projection) scanner.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_COUNT_DISTINCT),
		elements: []scanner.Element{p.(scanner.Element)},
		sel:      p.From(),
	}
}

func Cast(p scanner.Projection, stype grammar.SQLType) scanner.Projection {
	si := make([]grammar.Symbol, len(grammar.FunctionScanTable(grammar.FUNC_CAST)))
	copy(si, grammar.FunctionScanTable(grammar.FUNC_CAST))
	// Replace the placeholder with the SQL type's appropriate []byte
	// representation
	si[3] = grammar.SQLTypeToSymbol(stype)
	return &Function{
		ScanInfo: si,
		elements: []scanner.Element{p.(scanner.Element)},
	}
}

func CharLength(p scanner.Projection) scanner.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CHAR_LENGTH),
		elements: []scanner.Element{p.(scanner.Element)},
		sel:      p.From(),
	}
}

func BitLength(p scanner.Projection) scanner.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_BIT_LENGTH),
		elements: []scanner.Element{p.(scanner.Element)},
		sel:      p.From(),
	}
}

func Ascii(p scanner.Projection) scanner.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_ASCII),
		elements: []scanner.Element{p.(scanner.Element)},
		sel:      p.From(),
	}
}

func Reverse(p scanner.Projection) scanner.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_REVERSE),
		elements: []scanner.Element{p.(scanner.Element)},
		sel:      p.From(),
	}
}

func Concat(projs ...scanner.Projection) scanner.Projection {
	els := make([]scanner.Element, len(projs))
	for x, p := range projs {
		els[x] = p.(scanner.Element)
	}
	subjects := element.NewList(els...)
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CONCAT),
		elements: []scanner.Element{subjects},
		// TODO(jaypipes): Clearly we need to support >1 selection...
		sel: projs[0].From(),
	}
}

func ConcatWs(sep string, projs ...scanner.Projection) scanner.Projection {
	els := make([]scanner.Element, len(projs))
	for x, p := range projs {
		els[x] = p.(scanner.Element)
	}
	subjects := element.NewList(els...)
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CONCAT_WS),
		elements: []scanner.Element{element.NewValue(nil, sep), subjects},
		// TODO(jaypipes): Clearly we need to support >1 selection...
		sel: projs[0].From(),
	}
}

func Now() scanner.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_NOW),
	}
}

func CurrentTimestamp() scanner.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_TIMESTAMP),
	}
}

func CurrentTime() scanner.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_TIME),
	}
}

func CurrentDate() scanner.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_DATE),
	}
}

func Extract(p scanner.Projection, unit grammar.IntervalUnit) scanner.Projection {
	si := make([]grammar.Symbol, len(grammar.FunctionScanTable(grammar.FUNC_EXTRACT)))
	copy(si, grammar.FunctionScanTable(grammar.FUNC_EXTRACT))
	// Replace the placeholder with the interval unit's appropriate []byte
	// representation
	si[1] = grammar.IntervalUnitToSymbol(unit)
	return &Function{
		ScanInfo: si,
		elements: []scanner.Element{p.(scanner.Element)},
	}
}
