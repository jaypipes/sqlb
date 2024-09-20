//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package types

import "github.com/jaypipes/sqlb/core/grammar"

// QuerySpecificationConverter knows how to convert itself into a
// `*grammar.QuerySpecification`
type QuerySpecificationConverter interface {
	// QuerySpecification returns the object as a `*grammar.QuerySpecification`
	QuerySpecification() *grammar.QuerySpecification
}
