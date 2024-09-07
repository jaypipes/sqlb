//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doTableReference(
	el *grammar.TableReference,
	qargs []interface{},
	curarg *int,
) {
	if el.Primary != nil {
		b.doTablePrimary(el.Primary, qargs, curarg)
	} else if el.Joined != nil {
		b.doJoinedTable(el.Joined, qargs, curarg)
	}
}
