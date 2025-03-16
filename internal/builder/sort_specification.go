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

func (b *Builder) doSortSpecification(
	el *grammar.SortSpecification,
	qargs []interface{},
	curarg *int,
) {
	b.doValueExpression(&el.Key, qargs, curarg)
	if el.Order == grammar.OrderSpecificationDesc {
		b.WriteString(symbol.Space)
		b.WriteString(symbol.Desc)
	}
	if el.NullOrder != grammar.NullOrderSpecificationNone {
		b.WriteString(symbol.Space)
		b.WriteString(grammar.NullOrderSpecificationSymbol[el.NullOrder])
	}
}
