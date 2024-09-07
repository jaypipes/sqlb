//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doUnsignedLiteral(
	el *grammar.UnsignedLiteral,
	qargs []interface{},
	curarg *int,
) {
	if el.UnsignedNumericLiteral != nil {
		b.doScalar(el.UnsignedNumericLiteral.Value, qargs, curarg)
	} else {
		b.doScalar(el.GeneralLiteral.Value, qargs, curarg)
	}
}
