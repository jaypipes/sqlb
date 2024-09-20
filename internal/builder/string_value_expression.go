//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
)

func (b *Builder) doStringValueExpression(
	el *grammar.StringValueExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.Character != nil {
		b.doCharacterValueExpression(el.Character, qargs, curarg)
	} else if el.Blob != nil {
		b.doBlobValueExpression(el.Blob, qargs, curarg)
	}
}

func (b *Builder) doCharacterValueExpression(
	el *grammar.CharacterValueExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.Factor != nil {
		b.doCharacterFactor(el.Factor, qargs, curarg)
	}
}

func (b *Builder) doBlobValueExpression(
	el *grammar.BlobValueExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.Factor != nil {
		b.doBlobFactor(el.Factor, qargs, curarg)
	}
}

func (b *Builder) doCharacterFactor(
	el *grammar.CharacterFactor,
	qargs []interface{},
	curarg *int,
) {
	b.doCharacterPrimary(&el.Primary, qargs, curarg)
	if el.Collation != nil {
		b.Write(grammar.Symbols[grammar.SYM_COLLATE])
		b.WriteString(*el.Collation)
	}
}

func (b *Builder) doCharacterPrimary(
	el *grammar.CharacterPrimary,
	qargs []interface{},
	curarg *int,
) {
	if el.Primary != nil {
		b.doValueExpressionPrimary(el.Primary, qargs, curarg)
	} else if el.Function != nil {
		b.doStringValueFunction(el.Function, qargs, curarg)
	}
}

func (b *Builder) doBlobFactor(
	el *grammar.BlobFactor,
	qargs []interface{},
	curarg *int,
) {
	b.doBlobPrimary(&el.Primary, qargs, curarg)
}

func (b *Builder) doBlobPrimary(
	el *grammar.BlobPrimary,
	qargs []interface{},
	curarg *int,
) {
	if el.Primary != nil {
		b.doValueExpressionPrimary(el.Primary, qargs, curarg)
	} else if el.Function != nil {
		b.doStringValueFunction(el.Function, qargs, curarg)
	}
}
