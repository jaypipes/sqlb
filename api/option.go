//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

type Option struct {
	Dialect *Dialect
	Format  *FormatOptions
}

// OptionModifier modifies a Option
type OptionModifier func(o *Option)

// MergeOptions joins any OptionModifiers into one
func MergeOptions(mods []OptionModifier) *Option {
	o := &Option{}
	for _, mod := range mods {
		mod(o)
	}
	return o
}

// WithDialect informs sqlb of the Dialect
func WithDialect(
	dialect Dialect,
) OptionModifier {
	return func(o *Option) {
		o.Dialect = &dialect
	}
}

// WithFormatSeparateClauseWith instructs sqlb to use a supplied string as the
// separator between clauses
func WithFormatSeparateClauseWith(
	with string,
) OptionModifier {
	return func(o *Option) {
		if o.Format == nil {
			o.Format = &FormatOptions{}
		}
		o.Format.SeparateClauseWith = with
	}
}

// WithFormatPrefixWith instructs sqlb to use a supplied string as a prefix for
// the resulting SQL string
func WithFormatPrefixWith(
	with string,
) OptionModifier {
	return func(o *Option) {
		if o.Format == nil {
			o.Format = &FormatOptions{}
		}
		o.Format.PrefixWith = with
	}
}
