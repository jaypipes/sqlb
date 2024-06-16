//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sqlb

import (
	"github.com/jaypipes/sqlb/internal/grammar/identifier"
	"github.com/jaypipes/sqlb/meta"
)

// T returns a TableIdentifier of a given name from a supplied Meta
func T(m *meta.Meta, name string) *identifier.Table {
	return identifier.TableFromMeta(m, name)
}

// Reflect examines the supplied database connection and discovers Table
// definitions within that connection's associated database, returning a
// pointer to a Meta struct with the discovered information.
var Reflect = meta.Reflect
