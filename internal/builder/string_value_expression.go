//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doStringValueExpression(
	el *grammar.StringValueExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.CharacterValueExpression != nil {
		b.doCharacterValueExpression(el.CharacterValueExpression, qargs, curarg)
	} else if el.BlobValueExpression != nil {
		b.doBlobValueExpression(el.BlobValueExpression, qargs, curarg)
	}
}

func (b *Builder) doCharacterValueExpression(
	el *grammar.CharacterValueExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.CharacterFactor != nil {
		b.doCharacterFactor(el.CharacterFactor, qargs, curarg)
	}
}

func (b *Builder) doBlobValueExpression(
	el *grammar.BlobValueExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.BlobFactor != nil {
		b.doBlobFactor(el.BlobFactor, qargs, curarg)
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
	if el.ValueExpressionPrimary != nil {
		b.doValueExpressionPrimary(el.ValueExpressionPrimary, qargs, curarg)
	} else if el.StringValueFunction != nil {
		b.doStringValueFunction(el.StringValueFunction, qargs, curarg)
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
	if el.ValueExpressionPrimary != nil {
		b.doValueExpressionPrimary(el.ValueExpressionPrimary, qargs, curarg)
	} else if el.StringValueFunction != nil {
		b.doStringValueFunction(el.StringValueFunction, qargs, curarg)
	}
}
