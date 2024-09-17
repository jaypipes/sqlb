//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import "github.com/jaypipes/sqlb/grammar"

func (b *Builder) doDatetimeValueFunction(
	el *grammar.DatetimeValueFunction,
	qargs []interface{},
	curarg *int,
) {
	if el.CurrentDate {
		b.Write(grammar.Symbols[grammar.SYM_CURRENT_DATE])
	}
}
