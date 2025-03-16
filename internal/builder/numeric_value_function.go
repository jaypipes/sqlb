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
	} else if el.Exponential != nil {
		b.doExponentialFunction(el.Exponential, qargs, curarg)
	} else if el.SquareRoot != nil {
		b.doSquareRoot(el.SquareRoot, qargs, curarg)
	} else if el.Ceiling != nil {
		b.doCeilingFunction(el.Ceiling, qargs, curarg)
	} else if el.Floor != nil {
		b.doFloorFunction(el.Floor, qargs, curarg)
	}
}

func (b *Builder) doPositionExpression(
	el *grammar.PositionExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.String != nil {
		b.WriteString(symbol.Position)
		b.WriteString(symbol.LeftParen)
		b.doStringValueExpression(&el.String.Subject, qargs, curarg)
		b.WriteString(symbol.Space)
		b.WriteString(symbol.In)
		b.WriteString(symbol.Space)
		b.doStringValueExpression(&el.String.In, qargs, curarg)
		if el.String.Using != grammar.CharacterLengthUnitsCharacters {
			b.WriteString(symbol.Space)
			b.WriteString(symbol.Using)
			b.WriteString(symbol.Space)
			b.WriteString(grammar.CharacterLengthUnitsSymbol[el.String.Using])
		}
		b.WriteString(symbol.RightParen)
	} else if el.Blob != nil {
		b.WriteString(symbol.Position)
		b.WriteString(symbol.LeftParen)
		b.doBlobValueExpression(&el.Blob.Subject, qargs, curarg)
		b.WriteString(symbol.Space)
		b.WriteString(symbol.In)
		b.WriteString(symbol.Space)
		b.doBlobValueExpression(&el.Blob.In, qargs, curarg)
		b.WriteString(symbol.RightParen)
	}
}

func (b *Builder) doLengthExpression(
	el *grammar.LengthExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.Character != nil {
		b.WriteString(symbol.CharLength)
		b.WriteString(symbol.LeftParen)
		b.doStringValueExpression(&el.Character.Subject, qargs, curarg)
		if el.Character.Using != grammar.CharacterLengthUnitsCharacters {
			b.WriteString(symbol.Space)
			b.WriteString(symbol.Using)
			b.WriteString(symbol.Space)
			b.WriteString(grammar.CharacterLengthUnitsSymbol[el.Character.Using])
		}
		b.WriteString(symbol.RightParen)
	} else if el.Octet != nil {
		b.WriteString(symbol.OctetLength)
		b.WriteString(symbol.LeftParen)
		b.doStringValueExpression(&el.Octet.Subject, qargs, curarg)
		b.WriteString(symbol.RightParen)
	}
}

func (b *Builder) doExtractExpression(
	el *grammar.ExtractExpression,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.Extract)
	b.WriteString(symbol.LeftParen)
	b.doExtractField(&el.What, qargs, curarg)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.From)
	b.WriteString(symbol.Space)
	b.doExtractSource(&el.From, qargs, curarg)
	b.WriteString(symbol.RightParen)
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
	b.WriteString(symbol.Ln)
	b.WriteString(symbol.LeftParen)
	b.doNumericValueExpression(&el.Subject, qargs, curarg)
	b.WriteString(symbol.RightParen)
}

func (b *Builder) doAbsoluteValueExpression(
	el *grammar.AbsoluteValueExpression,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.Abs)
	b.WriteString(symbol.LeftParen)
	b.doNumericValueExpression(&el.Subject, qargs, curarg)
	b.WriteString(symbol.RightParen)
}

func (b *Builder) doExponentialFunction(
	el *grammar.ExponentialFunction,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.Exp)
	b.WriteString(symbol.LeftParen)
	b.doNumericValueExpression(&el.Subject, qargs, curarg)
	b.WriteString(symbol.RightParen)
}

func (b *Builder) doSquareRoot(
	el *grammar.SquareRoot,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.Sqrt)
	b.WriteString(symbol.LeftParen)
	b.doNumericValueExpression(&el.Subject, qargs, curarg)
	b.WriteString(symbol.RightParen)
}

func (b *Builder) doCeilingFunction(
	el *grammar.CeilingFunction,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.Ceil)
	b.WriteString(symbol.LeftParen)
	b.doNumericValueExpression(&el.Subject, qargs, curarg)
	b.WriteString(symbol.RightParen)
}

func (b *Builder) doFloorFunction(
	el *grammar.FloorFunction,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.Floor)
	b.WriteString(symbol.LeftParen)
	b.doNumericValueExpression(&el.Subject, qargs, curarg)
	b.WriteString(symbol.RightParen)
}
