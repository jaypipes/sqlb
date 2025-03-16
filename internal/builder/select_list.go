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

func (b *Builder) doSelectList(
	el *grammar.SelectList,
	qargs []interface{},
	curarg *int,
) {
	if el.Asterisk {
		b.WriteString(symbol.Asterisk)
	} else {
		for x, s := range el.Sublists {
			if x > 0 {
				b.WriteString(symbol.Comma)
				b.WriteString(symbol.Space)
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
