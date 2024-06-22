//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/api"
)

// defaultFormatOptions is the set of formatting options used if not specified
// by the user
var defaultFormatOptions = api.FormatOptions{
	SeparateClauseWith: " ",
	PrefixWith:         "",
}

// defaultDialect is the Dialect to use when not specified or able to be
// discovered from the DB Driver
var defaultDialect = api.Dialect(api.DialectMySQL)
