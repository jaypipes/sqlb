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

func (b *Builder) doWhereClause(
	el *grammar.WhereClause,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(b.opts.FormatSeparateClauseWith())
	b.WriteString(symbol.Where)
	b.WriteString(symbol.Space)
	b.doBooleanValueExpression(&el.Search, qargs, curarg)
}
