//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/api"
)

func (b *Builder) doTable(
	el *api.Table,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(el.Name())
}
