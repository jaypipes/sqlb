//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"strings"

	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doIdentifierChain(
	el *grammar.IdentifierChain,
	qargs []interface{},
	curarg *int,
) {
	b.WriteString(strings.Join(el.Identifiers, "."))
}

func (b *Builder) doSchemaQualifiedName(
	el *grammar.SchemaQualifiedName,
	qargs []interface{},
	curarg *int,
) {
	if el.SchemaName != nil {
		b.WriteString(*el.SchemaName)
		b.WriteRune('.')
	}
	b.doIdentifierChain(&el.Identifiers, qargs, curarg)
}
