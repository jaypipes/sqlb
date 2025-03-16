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
		b.doCharacterSubstringFunction(el.Substring, qargs, curarg)
	} else if el.RegexSubstring != nil {
		b.doRegexSubstringFunction(el.RegexSubstring, qargs, curarg)
	} else if el.Fold != nil {
		b.doFoldFunction(el.Fold, qargs, curarg)
	} else if el.Transcoding != nil {
		b.doTranscodingFunction(el.Transcoding, qargs, curarg)
	} else if el.Transliteration != nil {
		b.doCharacterTransliterationFunction(el.Transliteration, qargs, curarg)
	} else if el.Trim != nil {
		b.doTrimFunction(el.Trim, qargs, curarg)
	} else if el.Overlay != nil {
		b.doCharacterOverlayFunction(el.Overlay, qargs, curarg)
	} else if el.Normalize != nil {
		b.doNormalizeFunction(el.Normalize, qargs, curarg)
	}
}

func (b *Builder) doBlobValueFunction(
	el *grammar.BlobValueFunction,
	qargs []interface{},
	curarg *int,
) {
}

func (b *Builder) doCharacterSubstringFunction(
	el *grammar.CharacterSubstringFunction,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.Substring)
	b.WriteString(symbol.LeftParen)
	b.doCharacterValueExpression(&el.Subject, qargs, curarg)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.From)
	b.WriteString(symbol.Space)
	b.doNumericValueExpression(&el.From, qargs, curarg)
	if el.For != nil {
		b.WriteString(symbol.Space)
		b.WriteString(symbol.For)
		b.WriteString(symbol.Space)
		b.doNumericValueExpression(el.For, qargs, curarg)
	}
	if el.Using != grammar.CharacterLengthUnitsCharacters {
		b.WriteString(symbol.Space)
		b.WriteString(symbol.Using)
		b.WriteString(symbol.Space)
		b.WriteString(grammar.CharacterLengthUnitsSymbol[el.Using])
	}
	b.WriteString(symbol.RightParen)
}

func (b *Builder) doRegexSubstringFunction(
	el *grammar.RegexSubstringFunction,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.Substring)
	b.WriteString(symbol.LeftParen)
	b.doCharacterValueExpression(&el.Subject, qargs, curarg)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.Similar)
	b.WriteString(symbol.Space)
	b.doCharacterValueExpression(&el.Similar, qargs, curarg)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.Escape)
	b.WriteString(symbol.Space)
	b.doCharacterValueExpression(&el.Escape, qargs, curarg)
	b.WriteString(symbol.RightParen)
}

func (b *Builder) doFoldFunction(
	el *grammar.FoldFunction,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(grammar.FoldCaseSymbols[el.Case])
	b.WriteString(symbol.LeftParen)
	b.doCharacterValueExpression(&el.Subject, qargs, curarg)
	b.WriteString(symbol.RightParen)
}

func (b *Builder) doTranscodingFunction(
	el *grammar.TranscodingFunction,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.Convert)
	b.WriteString(symbol.LeftParen)
	b.doCharacterValueExpression(&el.Subject, qargs, curarg)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.Using)
	b.WriteString(symbol.Space)
	b.doSchemaQualifiedName(&el.Using, qargs, curarg)
	b.WriteString(symbol.RightParen)
}

func (b *Builder) doCharacterTransliterationFunction(
	el *grammar.CharacterTransliterationFunction,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.Translate)
	b.WriteString(symbol.LeftParen)
	b.doCharacterValueExpression(&el.Subject, qargs, curarg)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.Using)
	b.WriteString(symbol.Space)
	b.doSchemaQualifiedName(&el.Using, qargs, curarg)
	b.WriteString(symbol.RightParen)
}

func (b *Builder) doTrimFunction(
	el *grammar.TrimFunction,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.Trim)
	b.WriteString(symbol.LeftParen)
	if el.Specification != grammar.TrimSpecificationBoth {
		b.WriteString(grammar.TrimSpecificationSymbols[el.Specification])
		b.WriteString(symbol.Space)
	}
	if el.Character != nil {
		b.doCharacterValueExpression(el.Character, qargs, curarg)
		b.WriteString(symbol.Space)
		b.WriteString(symbol.From)
		b.WriteString(symbol.Space)
	}
	b.doCharacterValueExpression(&el.Subject, qargs, curarg)
	b.WriteString(symbol.RightParen)
}

func (b *Builder) doCharacterOverlayFunction(
	el *grammar.CharacterOverlayFunction,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.Overlay)
	b.WriteString(symbol.LeftParen)
	b.doCharacterValueExpression(&el.Subject, qargs, curarg)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.Placing)
	b.WriteString(symbol.Space)
	b.doCharacterValueExpression(&el.Placing, qargs, curarg)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.From)
	b.WriteString(symbol.Space)
	b.doNumericValueExpression(&el.From, qargs, curarg)
	if el.For != nil {
		b.WriteString(symbol.Space)
		b.WriteString(symbol.For)
		b.WriteString(symbol.Space)
		b.doNumericValueExpression(el.For, qargs, curarg)
	}
	if el.Using != grammar.CharacterLengthUnitsCharacters {
		b.WriteString(symbol.Space)
		b.WriteString(symbol.Using)
		b.WriteString(symbol.Space)
		b.WriteString(grammar.CharacterLengthUnitsSymbol[el.Using])
	}
	b.WriteString(symbol.RightParen)
}

func (b *Builder) doNormalizeFunction(
	el *grammar.NormalizeFunction,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(symbol.Normalize)
	b.WriteString(symbol.LeftParen)
	b.doCharacterValueExpression(&el.Subject, qargs, curarg)
	b.WriteString(symbol.RightParen)
}
