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
	} else if el.Length != nil {
		b.doLengthExpression(el.Length, qargs, curarg)
	}
}

func (b *Builder) doPositionExpression(
	el *grammar.PositionExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.String != nil {
		b.Write(grammar.Symbols[grammar.SYM_POSITION])
		b.doStringValueExpression(&el.String.Subject, qargs, curarg)
		b.Write(grammar.Symbols[grammar.SYM_IN])
		b.doStringValueExpression(&el.String.In, qargs, curarg)
		if el.String.Using != grammar.CharacterLengthUnitsCharacters {
			b.WriteRune(' ')
			b.Write(grammar.Symbols[grammar.SYM_USING])
			b.WriteString(grammar.CharacterLengthUnitsSymbol[el.String.Using])
		}
		b.Write(grammar.Symbols[grammar.SYM_RPAREN])
	} else if el.Blob != nil {
		b.Write(grammar.Symbols[grammar.SYM_POSITION])
		b.doBlobValueExpression(&el.Blob.Subject, qargs, curarg)
		b.Write(grammar.Symbols[grammar.SYM_IN])
		b.doBlobValueExpression(&el.Blob.In, qargs, curarg)
		b.Write(grammar.Symbols[grammar.SYM_RPAREN])
	}
}

func (b *Builder) doLengthExpression(
	el *grammar.LengthExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.Character != nil {
		b.Write(grammar.Symbols[grammar.SYM_CHAR_LENGTH])
		b.doStringValueExpression(&el.Character.Subject, qargs, curarg)
		if el.Character.Using != grammar.CharacterLengthUnitsCharacters {
			b.WriteRune(' ')
			b.Write(grammar.Symbols[grammar.SYM_USING])
			b.WriteString(grammar.CharacterLengthUnitsSymbol[el.Character.Using])
		}
		b.Write(grammar.Symbols[grammar.SYM_RPAREN])
	} else if el.Octet != nil {
		b.Write(grammar.Symbols[grammar.SYM_OCTET_LENGTH])
		b.doStringValueExpression(&el.Octet.Subject, qargs, curarg)
		b.Write(grammar.Symbols[grammar.SYM_RPAREN])
	}
}
