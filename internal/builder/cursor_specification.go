//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doCursorSpecification(
	el *grammar.CursorSpecification,
	qargs []interface{},
	curarg *int,
) {
	b.doQueryExpression(&el.QueryExpression, qargs, curarg)
	if el.OrderByClause != nil {
		b.doOrderByClause(el.OrderByClause, qargs, curarg)
	}
	if el.LimitClause != nil {
		b.doLimitClause(el.LimitClause, qargs, curarg)
	}
}
