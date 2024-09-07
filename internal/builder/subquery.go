//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doSubquery(
	el *grammar.Subquery,
	qargs []interface{},
	curarg *int,
) {
	b.Write(grammar.Symbols[grammar.SYM_LPAREN])
	b.doQueryExpression(&el.QueryExpression, qargs, curarg)
	b.Write(grammar.Symbols[grammar.SYM_RPAREN])
}
