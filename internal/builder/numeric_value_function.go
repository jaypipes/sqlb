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
	} else if el.Extract != nil {
		b.doExtractExpression(el.Extract, qargs, curarg)
	} else if el.Natural != nil {
		b.doNaturalLogarithm(el.Natural, qargs, curarg)
	} else if el.AbsoluteValue != nil {
		b.doAbsoluteValueExpression(el.AbsoluteValue, qargs, curarg)
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

func (b *Builder) doExtractExpression(
	el *grammar.ExtractExpression,
	qargs []interface{},
	curarg *int,
) {
	b.Write(grammar.Symbols[grammar.SYM_EXTRACT])
	b.doExtractField(&el.What, qargs, curarg)
	b.WriteRune(' ')
	b.Write(grammar.Symbols[grammar.SYM_FROM])
	b.doExtractSource(&el.From, qargs, curarg)
	b.Write(grammar.Symbols[grammar.SYM_RPAREN])
}

func (b *Builder) doExtractField(
	el *grammar.ExtractField,
	qargs []interface{},
	curarg *int,
) {
	if el.Datetime != nil {
		b.doPrimaryDatetimeField(el.Datetime, qargs, curarg)
	} else if el.Timezone != nil {
		b.WriteString(grammar.TimezoneFieldSymbols[*el.Timezone])
	}
}

func (b *Builder) doExtractSource(
	el *grammar.ExtractSource,
	qargs []interface{},
	curarg *int,
) {
	if el.Datetime != nil {
		b.doDatetimeValueExpression(el.Datetime, qargs, curarg)
	} else if el.Interval != nil {
		b.doIntervalValueExpression(el.Interval, qargs, curarg)
	}
}

func (b *Builder) doNaturalLogarithm(
	el *grammar.NaturalLogarithm,
	qargs []interface{},
	curarg *int,
) {
	b.Write(grammar.Symbols[grammar.SYM_LN])
	b.doNumericValueExpression(&el.Subject, qargs, curarg)
	b.Write(grammar.Symbols[grammar.SYM_RPAREN])
}

func (b *Builder) doAbsoluteValueExpression(
	el *grammar.AbsoluteValueExpression,
	qargs []interface{},
	curarg *int,
) {
	b.Write(grammar.Symbols[grammar.SYM_ABS])
	b.doNumericValueExpression(&el.Subject, qargs, curarg)
	b.Write(grammar.Symbols[grammar.SYM_RPAREN])
}
