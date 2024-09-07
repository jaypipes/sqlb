//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doGroupByClause(
	el *grammar.GroupByClause,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(b.opts.FormatSeparateClauseWith())
	b.Write(grammar.Symbols[grammar.SYM_GROUP_BY])
	for x, ge := range el.GroupingElements {
		if x > 0 {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
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
		b.Write(grammar.Symbols[grammar.SYM_COLLATE])
		b.WriteString(*el.GroupingColumnReference.Collation)
	}
}
