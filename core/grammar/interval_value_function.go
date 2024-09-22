//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <interval value function>    ::=   <interval absolute value function>
//
// <interval absolute value function>    ::=   ABS <left paren> <interval value expression> <right paren>

type IntervalValueFunction struct {
	Abs *IntervalValueExpression
}

func (f *IntervalValueFunction) ArgCount(count *int) {
	if f.Abs != nil {
		f.Abs.ArgCount(count)
	}
}
