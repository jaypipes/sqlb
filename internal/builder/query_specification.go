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

func (b *Builder) doQuerySpecification(
	el *grammar.QuerySpecification,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.Select)
	b.WriteString(symbol.Space)
	b.doSelectList(&el.SelectList, qargs, curarg)
	b.doTableExpression(&el.TableExpression, qargs, curarg)
}
