//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doQuerySpecification(
	el *grammar.QuerySpecification,
	qargs []interface{},
	curarg *int,
) {
	b.Write(grammar.Symbols[grammar.SYM_SELECT])
	b.doSelectList(&el.SelectList, qargs, curarg)
	b.doTableExpression(&el.TableExpression, qargs, curarg)
}
