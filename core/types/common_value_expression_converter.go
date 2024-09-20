//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package types

import "github.com/jaypipes/sqlb/core/grammar"

// CommonValueExpressionConverter knows how to convert itself into a
// `*grammar.CommonValueExpression`
type CommonValueExpressionConverter interface {
	// CommonValueExpression returns the object as a
	// `*grammar.CommonValueExpression`
	CommonValueExpression() *grammar.CommonValueExpression
}
