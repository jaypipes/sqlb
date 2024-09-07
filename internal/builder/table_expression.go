//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doTableExpression(
	el *grammar.TableExpression,
	qargs []interface{},
	curarg *int,
) {
	b.doFromClause(&el.FromClause, qargs, curarg)
	if el.WhereClause != nil {
		b.doWhereClause(el.WhereClause, qargs, curarg)
	}
	if el.GroupByClause != nil {
		b.doGroupByClause(el.GroupByClause, qargs, curarg)
	}
	if el.HavingClause != nil {
		b.doHavingClause(el.HavingClause, qargs, curarg)
	}
}
