//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
	"github.com/jaypipes/sqlb/core/grammar/symbol"
)

func (b *Builder) doTablePrimary(
	el *grammar.TablePrimary,
	qargs []interface{},
	curarg *int,
) {
	if el.TableName != nil {
		b.WriteString(*el.TableName)
	} else if el.QueryName != nil {
		b.WriteString(*el.QueryName)
	} else if el.DerivedTable != nil {
		b.doDerivedTable(el.DerivedTable, qargs, curarg)
	}
	if el.Correlation != nil {
		b.WriteString(symbol.Space)
		b.WriteString(symbol.As)
		b.WriteString(symbol.Space)
		b.WriteString(el.Correlation.Name)
	}
}
