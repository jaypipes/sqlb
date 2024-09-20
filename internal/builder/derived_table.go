//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
)

func (b *Builder) doDerivedTable(
	el *grammar.DerivedTable,
	qargs []interface{},
	curarg *int,
) {
	b.doSubquery(&el.Subquery, qargs, curarg)
}
