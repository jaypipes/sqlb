//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/api"
)

type BuilderOption struct {
	Dialect *api.Dialect
	Format  *api.FormatOptions
}

// BuilderOptionModifier modifies a BuilderOption
type BuilderOptionModifier func(o *BuilderOption)

// mergeOpts joins any BuilderOptionModifiers into one
func mergeOpts(mods []BuilderOptionModifier) *BuilderOption {
	o := &BuilderOption{}
	for _, mod := range mods {
		mod(o)
	}
	return o
}

// WithDialect informs the builder of the Dialect
func WithDialect(
	dialect api.Dialect,
) BuilderOptionModifier {
	return func(o *BuilderOption) {
		o.Dialect = &dialect
	}
}

// WithFormatSeparateClauseWith instructs the builder to use a supplied string
// as the separator between clauses
func WithFormatSeparateClauseWith(
	with string,
) BuilderOptionModifier {
	return func(o *BuilderOption) {
		if o.Format == nil {
			o.Format = &api.FormatOptions{}
		}
		o.Format.SeparateClauseWith = with
	}
}

// WithFormatPrefixWith instructs the builder to use a supplied string
// as a prefix for the resulting SQL string
func WithFormatPrefixWith(
	with string,
) BuilderOptionModifier {
	return func(o *BuilderOption) {
		if o.Format == nil {
			o.Format = &api.FormatOptions{}
		}
		o.Format.PrefixWith = with
	}
}
