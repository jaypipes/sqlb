//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doColumn(
	el *api.Column,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(el.TableName())
	b.Write(grammar.Symbols[grammar.SYM_PERIOD])
	b.WriteString(el.Name())
}
