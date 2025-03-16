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

func (b *Builder) doAggregateFunction(
	el *grammar.AggregateFunction,
	qargs []interface{},
	curarg *int,
) {
	if el.CountStar != nil {
		b.WriteString(symbol.Count)
		b.WriteString(symbol.LeftParen)
		b.WriteString(symbol.Asterisk)
		b.WriteString(symbol.RightParen)
	} else if el.GeneralSet != nil {
		b.doGeneralSetFunction(el.GeneralSet, qargs, curarg)
	}
}

func (b *Builder) doGeneralSetFunction(
	el *grammar.GeneralSetFunction,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(grammar.ComputationalOperationSymbol[el.Operation])
	b.WriteString(symbol.LeftParen)
	if el.Quantifier == grammar.SetQuantifierDistinct {
		b.WriteString(symbol.Distinct)
		b.WriteString(symbol.Space)
	}
	b.doValueExpression(&el.Value, qargs, curarg)
	b.WriteString(symbol.RightParen)
}
