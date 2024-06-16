//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package meta

import (
	"github.com/jaypipes/sqlb/types"
)

type MetaOption struct {
	Dialect types.Dialect
}

// MetaOptionModifier modifies a MetaOption
type MetaOptionModifier func(o *MetaOption)

// mergeOpts joins any MetaOptionModifiers into one
func mergeOpts(mods []MetaOptionModifier) *MetaOption {
	o := &MetaOption{}
	for _, mod := range mods {
		mod(o)
	}
	return o
}

// WithDialect informs the supplied function of the Dialect. If not passed, the
// `sql.DB` handle is queried for the dialect.
func WithDialect(
	dialect types.Dialect,
) MetaOptionModifier {
	return func(o *MetaOption) {
		o.Dialect = dialect
	}
}
