//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
)

func (b *Builder) doSelectList(
	el *grammar.SelectList,
	qargs []interface{},
	curarg *int,
) {
	if el.Asterisk {
		b.Write(grammar.Symbols[grammar.SYM_ASTERISK])
	} else {
		for x, s := range el.Sublists {
			if x > 0 {
				b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
			}
			b.doSelectSublist(&s, qargs, curarg)
		}
	}
}

func (b *Builder) doSelectSublist(
	el *grammar.SelectSublist,
	qargs []interface{},
	curarg *int,
) {
	if el.Asterisk {
		// TODO(jaypipes): Handle qualified asterisk expression
	} else {
		b.doDerivedColumn(el.DerivedColumn, qargs, curarg)
	}
}
