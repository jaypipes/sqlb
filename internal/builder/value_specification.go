//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doUnsignedValueSpecification(
	el *grammar.UnsignedValueSpecification,
	qargs []interface{},
	curarg *int,
) {
	if el.UnsignedLiteral != nil {
		b.doUnsignedLiteral(el.UnsignedLiteral, qargs, curarg)
	}
}
