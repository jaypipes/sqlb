//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

// DefaultFormatOptions is the set of formatting options used if not specified
// by the user
var DefaultFormatOptions = FormatOptions{
	SeparateClauseWith: " ",
	PrefixWith:         "",
}

// DefaultDialect is the Dialect to use when not specified or able to be
// discovered from the DB Driver
var DefaultDialect = DialectMySQL
