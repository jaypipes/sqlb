//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sqlb

import (
	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/schema"
)

// T returns a TableIdentifier of a given name from a supplied Schema
func T(s *schema.Schema, name string) *ast.TableIdentifier {
	return ast.TableIdentifierFromSchema(s, name)
}
