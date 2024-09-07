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
	if el.CharacterValueFunction != nil {
		b.doCharacterValueFunction(el.CharacterValueFunction, qargs, curarg)
	} else if el.BlobValueFunction != nil {
		b.doBlobValueFunction(el.BlobValueFunction, qargs, curarg)
	}
}

func (b *Builder) doCharacterValueFunction(
	el *grammar.CharacterValueFunction,
	qargs []interface{},
	curarg *int,
) {
}

func (b *Builder) doBlobValueFunction(
	el *grammar.BlobValueFunction,
	qargs []interface{},
	curarg *int,
) {
}
