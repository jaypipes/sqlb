//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"strconv"

	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doDatetimeValueFunction(
	el *grammar.DatetimeValueFunction,
	qargs []interface{},
	curarg *int,
) {
	if el.CurrentDate {
		b.Write(grammar.Symbols[grammar.SYM_CURRENT_DATE])
	} else if el.CurrentTime != nil {
		b.Write(grammar.Symbols[grammar.SYM_CURRENT_TIME])
		if el.CurrentTime.Precision != nil {
			b.WriteString(strconv.Itoa(int(*el.CurrentTime.Precision)))
		}
		b.Write(grammar.Symbols[grammar.SYM_RPAREN])
	} else if el.CurrentTimestamp != nil {
		b.Write(grammar.Symbols[grammar.SYM_CURRENT_TIMESTAMP])
		if el.CurrentTimestamp.Precision != nil {
			b.WriteString(strconv.Itoa(int(*el.CurrentTimestamp.Precision)))
		}
		b.Write(grammar.Symbols[grammar.SYM_RPAREN])
	} else if el.LocalTime != nil {
		b.Write(grammar.Symbols[grammar.SYM_LOCALTIME])
		if el.LocalTime.Precision != nil {
			b.WriteString(strconv.Itoa(int(*el.LocalTime.Precision)))
		}
		b.Write(grammar.Symbols[grammar.SYM_RPAREN])
	} else if el.LocalTimestamp != nil {
		b.Write(grammar.Symbols[grammar.SYM_LOCALTIMESTAMP])
		if el.LocalTimestamp.Precision != nil {
			b.WriteString(strconv.Itoa(int(*el.LocalTimestamp.Precision)))
		}
		b.Write(grammar.Symbols[grammar.SYM_RPAREN])
	}
}
