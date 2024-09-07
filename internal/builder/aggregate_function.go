//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doAggregateFunction(
	el *grammar.AggregateFunction,
	qargs []interface{},
	curarg *int,
) {
	if el.CountStar != nil {
		b.Write(grammar.Symbols[grammar.SYM_COUNT_STAR])
	} else if el.GeneralSetFunction != nil {
		b.doGeneralSetFunction(el.GeneralSetFunction, qargs, curarg)
	}
}

func (b *Builder) doGeneralSetFunction(
	el *grammar.GeneralSetFunction,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(grammar.ComputationalOperationSymbol[el.Operation])
	b.Write(grammar.Symbols[grammar.SYM_LPAREN])
	if el.Quantifier == grammar.SetQuantifierDistinct {
		b.Write(grammar.Symbols[grammar.SYM_DISTINCT])
	}
	b.doValueExpression(&el.ValueExpression, qargs, curarg)
	b.Write(grammar.Symbols[grammar.SYM_RPAREN])
}
