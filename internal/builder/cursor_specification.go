//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
)

func (b *Builder) doCursorSpecification(
	el *grammar.CursorSpecification,
	qargs []interface{},
	curarg *int,
) {
	b.doQueryExpression(&el.Query, qargs, curarg)
	if el.OrderBy != nil {
		b.doOrderByClause(el.OrderBy, qargs, curarg)
	}
	if el.Limit != nil {
		b.doLimitClause(el.Limit, qargs, curarg)
	}
}
