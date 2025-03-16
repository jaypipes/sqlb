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

func (b *Builder) doSubquery(
	el *grammar.Subquery,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.LeftParen)
	b.doQueryExpression(&el.QueryExpression, qargs, curarg)
	b.WriteString(symbol.RightParen)
}
