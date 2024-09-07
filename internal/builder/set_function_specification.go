//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doSetFunctionSpecification(
	el *grammar.SetFunctionSpecification,
	qargs []interface{},
	curarg *int,
) {
	if el.AggregateFunction != nil {
		b.doAggregateFunction(el.AggregateFunction, qargs, curarg)
	}
}
