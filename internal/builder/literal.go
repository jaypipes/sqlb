//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
)

func (b *Builder) doUnsignedLiteral(
	el *grammar.UnsignedLiteral,
	qargs []interface{},
	curarg *int,
) {
	if el.UnsignedNumeric != nil {
		b.doScalar(el.UnsignedNumeric.Value, qargs, curarg)
	} else {
		b.doScalar(el.General.Value, qargs, curarg)
	}
}
