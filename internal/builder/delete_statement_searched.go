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

func (b *Builder) doDeleteStatementSearched(
	el *grammar.DeleteStatementSearched,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.Delete)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.From)
	b.WriteString(symbol.Space)
	// We don't add any table alias when outputting the table identifier
	b.WriteString(el.TableName)

	if el.Where != nil {
		b.doWhereClause(el.Where, qargs, curarg)
	}
}
