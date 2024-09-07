//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doColumnReference(
	el *grammar.ColumnReference,
	qargs []interface{},
	curarg *int,
) {
	if el.BasicIdentifierChain != nil {
		b.doIdentifierChain(el.BasicIdentifierChain, qargs, curarg)
	}
	if el.Correlation != nil {
		b.Write(grammar.Symbols[grammar.SYM_AS])
		b.WriteString(el.Correlation.Name)
	}
}
