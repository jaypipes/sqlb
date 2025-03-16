//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"strconv"

	"github.com/jaypipes/sqlb/core/grammar"
	"github.com/jaypipes/sqlb/core/grammar/symbol"
)

func (b *Builder) doDatetimeValueFunction(
	el *grammar.DatetimeValueFunction,
	qargs []interface{},
	curarg *int,
) {
	if el.CurrentDate {
		b.WriteString(symbol.CurrentDate)
		b.WriteString(symbol.LeftParen)
		b.WriteString(symbol.RightParen)
	} else if el.CurrentTime != nil {
		b.WriteString(symbol.CurrentTime)
		b.WriteString(symbol.LeftParen)
		if el.CurrentTime.Precision != nil {
			b.WriteString(strconv.Itoa(int(*el.CurrentTime.Precision)))
		}
		b.WriteString(symbol.RightParen)
	} else if el.CurrentTimestamp != nil {
		b.WriteString(symbol.CurrentTimestamp)
		b.WriteString(symbol.LeftParen)
		if el.CurrentTimestamp.Precision != nil {
			b.WriteString(strconv.Itoa(int(*el.CurrentTimestamp.Precision)))
		}
		b.WriteString(symbol.RightParen)
	} else if el.LocalTime != nil {
		b.WriteString(symbol.LocalTime)
		b.WriteString(symbol.LeftParen)
		if el.LocalTime.Precision != nil {
			b.WriteString(strconv.Itoa(int(*el.LocalTime.Precision)))
		}
		b.WriteString(symbol.RightParen)
	} else if el.LocalTimestamp != nil {
		b.WriteString(symbol.LocalTimestamp)
		b.WriteString(symbol.LeftParen)
		if el.LocalTimestamp.Precision != nil {
			b.WriteString(strconv.Itoa(int(*el.LocalTimestamp.Precision)))
		}
		b.WriteString(symbol.RightParen)
	}
}
