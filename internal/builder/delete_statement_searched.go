//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
)

func (b *Builder) doDeleteStatementSearched(
	el *grammar.DeleteStatementSearched,
	qargs []interface{},
	curarg *int,
) {
	b.Write(grammar.Symbols[grammar.SYM_DELETE])
	// We don't add any table alias when outputting the table identifier
	b.WriteString(el.TableName)

	if el.WhereClause != nil {
		b.doWhereClause(el.WhereClause, qargs, curarg)
	}
}
