//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package function

import (
	"strings"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/element"
	"github.com/jaypipes/sqlb/internal/grammar/sortcolumn"
)

// Function is a SQL function that accepts zero or more parameters
type Function struct {
	sel      api.Selection
	Alias    string
	ScanInfo grammar.ScanInfo
	elements []api.Element
}

func (f *Function) From() api.Selection {
	return f.sel
}

func (f *Function) DisableAliasScan() func() {
	origAlias := f.Alias
	f.Alias = ""
	return func() { f.Alias = origAlias }
}

func (f *Function) As(alias string) api.Projection {
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

func (f *Function) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	elidx := 0
	for _, sym := range f.ScanInfo {
		if sym == grammar.SYM_ELEMENT {
			el := f.elements[elidx]
			// We need to disable alias output for elements that are
			// projections. We don't want to output, for example,
			// "ON users.id AS user_id = articles.author"
			switch el := el.(type) {
			case api.Projection:
				reset := el.DisableAliasScan()
				defer reset()
			}
			elidx++
			b.WriteString(el.String(opts, qargs, curarg))
		} else {
			b.Write(grammar.Symbols[sym])
		}
	}
	if f.Alias != "" {
		b.Write(grammar.Symbols[grammar.SYM_AS])
		b.WriteString(f.Alias)
	}
	return b.String()
}

func (f *Function) Desc() api.Orderable {
	return sortcolumn.NewDesc(f)
}

func (f *Function) Asc() api.Orderable {
	return sortcolumn.NewAsc(f)
}

// Max returns a Projection that contains the MAX() SQL function
func Max(p api.Projection) api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_MAX),
		elements: []api.Element{p.(api.Element)},
		sel:      p.From(),
	}
}

// Min returns a Projection that contains the MIN() SQL function
func Min(p api.Projection) api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_MIN),
		elements: []api.Element{p.(api.Element)},
		sel:      p.From(),
	}
}

// Sum returns a Projection that contains the SUM() SQL function
func Sum(p api.Projection) api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_SUM),
		elements: []api.Element{p.(api.Element)},
		sel:      p.From(),
	}
}

// Avg returns a Projection that contains the AVG() SQL function
func Avg(p api.Projection) api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_AVG),
		elements: []api.Element{p.(api.Element)},
		sel:      p.From(),
	}
}

// Count returns a Projection that contains the COUNT() SQL function
func Count(sel api.Selection) api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_COUNT_STAR),
		sel:      sel,
	}
}

// CountDistint returns a Projection that contains the COUNT(x DISTINCT) SQL
// function
func CountDistinct(p api.Projection) api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_COUNT_DISTINCT),
		elements: []api.Element{p.(api.Element)},
		sel:      p.From(),
	}
}

// Cast returns a Projection that contains the CAST() SQL function
func Cast(p api.Projection, stype grammar.SQLType) api.Projection {
	si := make([]grammar.Symbol, len(grammar.FunctionScanTable(grammar.FUNC_CAST)))
	copy(si, grammar.FunctionScanTable(grammar.FUNC_CAST))
	// Replace the placeholder with the SQL type's appropriate []byte
	// representation
	si[3] = grammar.SQLTypeToSymbol(stype)
	return &Function{
		ScanInfo: si,
		elements: []api.Element{p.(api.Element)},
	}
}

// CharLength returns a Projection that contains the CHAR_LENGTH() SQL function
func CharLength(p api.Projection) api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CHAR_LENGTH),
		elements: []api.Element{p.(api.Element)},
		sel:      p.From(),
	}
}

// BitLength returns a Projection that contains the BIT_LENGTH() SQL function
func BitLength(p api.Projection) api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_BIT_LENGTH),
		elements: []api.Element{p.(api.Element)},
		sel:      p.From(),
	}
}

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

// Now returns a Projection that contains the NOW() SQL function
func Now() api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_NOW),
	}
}

// CurrentTimestamp returns a Projection that contains the CURRENT_TIMESTAMP() SQL function
func CurrentTimestamp() api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_TIMESTAMP),
	}
}

// CurrentTime returns a Projection that contains the CURRENT_TIME() SQL function
func CurrentTime() api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_TIME),
	}
}

// CurrentDate returns a Projection that contains the CURRENT_DATE() SQL function
func CurrentDate() api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_DATE),
	}
}

// Extract returns a Projection that contains the EXTRACT() SQL function
func Extract(p api.Projection, unit grammar.IntervalUnit) api.Projection {
	si := make([]grammar.Symbol, len(grammar.FunctionScanTable(grammar.FUNC_EXTRACT)))
	copy(si, grammar.FunctionScanTable(grammar.FUNC_EXTRACT))
	// Replace the placeholder with the interval unit's appropriate []byte
	// representation
	si[1] = grammar.IntervalUnitToSymbol(unit)
	return &Function{
		ScanInfo: si,
		elements: []api.Element{p.(api.Element)},
	}
}
