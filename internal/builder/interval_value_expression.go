//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doIntervalValueExpression(
	el *grammar.IntervalValueExpression,
	qargs []interface{},
	curarg *int,
) {
}

func (b *Builder) doPrimaryDatetimeField(
	el *grammar.PrimaryDatetimeField,
	qargs []interface{},
	curarg *int,
) {
	if el.Second {
		b.WriteString("SECOND")
	} else {
		b.WriteString(grammar.NonsecondPrimaryDatetimeFieldSymbols[el.Nonsecond])
	}
}
