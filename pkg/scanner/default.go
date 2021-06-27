//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package scanner

import (
	"github.com/jaypipes/sqlb/pkg/types"
)

// DefaultFormatOptions is the set of formatting options used if not specified
// by the user
var DefaultFormatOptions = &types.FormatOptions{
	SeparateClauseWith: " ",
	PrefixWith:         "",
}

// DefaultScanner is the default scanner used if not specified by the user
var DefaultScanner = &sqlScanner{
	dialect: types.DIALECT_MYSQL,
	format:  DefaultFormatOptions,
}
