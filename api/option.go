//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

type Options struct {
	dialect *Dialect
	format  *FormatOptions
}

// HasDialect returns true if the Options' Dialect has been set
func (o *Options) HasDialect() bool {
	return o != nil && o.dialect != nil
}

// Dialect returns the Options' Dialect or the default Dialect if not set
func (o *Options) Dialect() Dialect {
	if o == nil || o.dialect == nil {
		return DefaultDialect
	}
	return *o.dialect
}

// FormatPrefixWith returns the Options' FormatPrefixWith or the default if not
// set
func (o *Options) FormatPrefixWith() string {
	if o == nil || o.format == nil {
		return DefaultFormatOptions.PrefixWith
	}
	return (*o.format).PrefixWith
}

// FormatSeparateClauseWith returns the Options' FormatSeparateClauseWith or
// the default if not set
func (o *Options) FormatSeparateClauseWith() string {
	if o == nil || o.format == nil {
		return DefaultFormatOptions.SeparateClauseWith
	}
	return (*o.format).SeparateClauseWith
}

// Option modifies an Options
type Option func(o *Options)

// MergeOptions joins any OptionModifiers into one
func MergeOptions(mods []Option) Options {
	opts := Options{}
	for _, mod := range mods {
		mod(&opts)
	}
	return opts
}

// WithDialect informs sqlb of the Dialect
func WithDialect(
	dialect Dialect,
) Option {
	return func(o *Options) {
		o.dialect = &dialect
	}
}

// WithFormatSeparateClauseWith instructs sqlb to use a supplied string as the
// separator between clauses
func WithFormatSeparateClauseWith(
	with string,
) Option {
	return func(o *Options) {
		if o.format == nil {
			o.format = &FormatOptions{}
		}
		o.format.SeparateClauseWith = with
	}
}

// WithFormatPrefixWith instructs sqlb to use a supplied string as a prefix for
// the resulting SQL string
func WithFormatPrefixWith(
	with string,
) Option {
	return func(o *Options) {
		if o.format == nil {
			o.format = &FormatOptions{}
		}
		o.format.PrefixWith = with
	}
}
