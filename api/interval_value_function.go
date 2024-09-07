//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

/*
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
*/
