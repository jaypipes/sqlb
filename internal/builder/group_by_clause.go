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

func (b *Builder) doGroupByClause(
	el *grammar.GroupByClause,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(b.opts.FormatSeparateClauseWith())
	b.WriteString(symbol.Group)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.By)
	b.WriteString(symbol.Space)
	for x, ge := range el.GroupingElements {
		if x > 0 {
			b.WriteString(symbol.Comma)
			b.WriteString(symbol.Space)
		}
		b.doGroupingElement(&ge, qargs, curarg)
	}
}

func (b *Builder) doGroupingElement(
	el *grammar.GroupingElement,
	qargs []interface{},
	curarg *int,
) {
	if el.OrdinaryGroupingSet != nil {
		b.doOrdinaryGroupingSet(el.OrdinaryGroupingSet, qargs, curarg)
	}
}

func (b *Builder) doOrdinaryGroupingSet(
	el *grammar.OrdinaryGroupingSet,
	qargs []interface{},
	curarg *int,
) {
	b.doColumnReference(el.GroupingColumnReference.ColumnReference, qargs, curarg)
	if el.GroupingColumnReference.Collation != nil {
		b.WriteString(symbol.Collate)
		b.WriteString(symbol.Space)
		b.WriteString(*el.GroupingColumnReference.Collation)
	}
}
