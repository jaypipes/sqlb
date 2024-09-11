//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import "github.com/jaypipes/sqlb/grammar"

func (b *Builder) doStringValueFunction(
	el *grammar.StringValueFunction,
	qargs []interface{},
	curarg *int,
) {
	if el.Character != nil {
		b.doCharacterValueFunction(el.Character, qargs, curarg)
	} else if el.Blob != nil {
		b.doBlobValueFunction(el.Blob, qargs, curarg)
	}
}

func (b *Builder) doCharacterValueFunction(
	el *grammar.CharacterValueFunction,
	qargs []interface{},
	curarg *int,
) {
	if el.Substring != nil {
		b.doSubstringFunction(el.Substring, qargs, curarg)
	}
}

func (b *Builder) doBlobValueFunction(
	el *grammar.BlobValueFunction,
	qargs []interface{},
	curarg *int,
) {
}

func (b *Builder) doSubstringFunction(
	el *grammar.SubstringFunction,
	qargs []interface{},
	curarg *int,
) {
	b.Write(grammar.Symbols[grammar.SYM_SUBSTRING])
	b.doCharacterValueExpression(&el.Subject, qargs, curarg)
	b.WriteRune(' ')
	b.Write(grammar.Symbols[grammar.SYM_FROM])
	b.doNumericValueExpression(&el.From, qargs, curarg)
	if el.For != nil {
		b.WriteRune(' ')
		b.Write(grammar.Symbols[grammar.SYM_FOR])
		b.doNumericValueExpression(el.For, qargs, curarg)
	}
	if el.Using != grammar.CharacterLengthUnitsCharacters {
		b.WriteRune(' ')
		b.Write(grammar.Symbols[grammar.SYM_USING])
		b.WriteString(grammar.CharacterLengthUnitsSymbol[el.Using])
	}
	b.Write(grammar.Symbols[grammar.SYM_RPAREN])
}