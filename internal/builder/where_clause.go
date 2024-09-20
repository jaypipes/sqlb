//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
)

func (b *Builder) doWhereClause(
	el *grammar.WhereClause,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(b.opts.FormatSeparateClauseWith())
	b.Write(grammar.Symbols[grammar.SYM_WHERE])
	b.doBooleanValueExpression(&el.SearchCondition, qargs, curarg)
}
