//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doSortSpecification(
	el *grammar.SortSpecification,
	qargs []interface{},
	curarg *int,
) {
	b.doValueExpression(&el.Key, qargs, curarg)
	if el.Order == grammar.OrderSpecificationDesc {
		b.Write(grammar.Symbols[grammar.SYM_DESC])
	}
	if el.NullOrder != grammar.NullOrderSpecificationNone {
		b.WriteRune(' ')
		b.WriteString(grammar.NullOrderSpecificationSymbol[el.NullOrder])
	}
}
