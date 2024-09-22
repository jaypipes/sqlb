//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
)

func (b *Builder) doTableExpression(
	el *grammar.TableExpression,
	qargs []interface{},
	curarg *int,
) {
	b.doFromClause(&el.From, qargs, curarg)
	if el.Where != nil {
		b.doWhereClause(el.Where, qargs, curarg)
	}
	if el.GroupBy != nil {
		b.doGroupByClause(el.GroupBy, qargs, curarg)
	}
	if el.Having != nil {
		b.doHavingClause(el.Having, qargs, curarg)
	}
}
