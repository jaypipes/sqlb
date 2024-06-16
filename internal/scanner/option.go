//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package scanner

import (
	"github.com/jaypipes/sqlb/types"
)

type ScannerOption struct {
	Dialect types.Dialect
	Format  types.FormatOptions
}

// ScannerOptionModifier modifies a ScannerOption
type ScannerOptionModifier func(o *ScannerOption)

// mergeOpts joins any ScannerOptionModifiers into one
func mergeOpts(mods []ScannerOptionModifier) *ScannerOption {
	o := &ScannerOption{}
	for _, mod := range mods {
		mod(o)
	}
	return o
}

// WithDialect informs the scanner of the Dialect
func WithDialect(
	dialect types.Dialect,
) ScannerOptionModifier {
	return func(o *ScannerOption) {
		o.Dialect = dialect
	}
}

// WithFormatSeparateClauseWith instructs the scanner to use a supplied string
// as the separator between clauses
func WithFormatSeparateClauseWith(
	with string,
) ScannerOptionModifier {
	return func(o *ScannerOption) {
		o.Format.SeparateClauseWith = with
	}
}

// WithFormatPrefixWith instructs the scanner to use a supplied string
// as a prefix for the resulting SQL string
func WithFormatPrefixWith(
	with string,
) ScannerOptionModifier {
	return func(o *ScannerOption) {
		o.Format.PrefixWith = with
	}
}
