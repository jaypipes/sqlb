//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doNumericValueFunction(
	el *grammar.NumericValueFunction,
	qargs []interface{},
	curarg *int,
) {
	if el.Position != nil {
		b.doPositionExpression(el.Position, qargs, curarg)
	}
}

func (b *Builder) doPositionExpression(
	el *grammar.PositionExpression,
	qargs []interface{},
	curarg *int,
) {
}
