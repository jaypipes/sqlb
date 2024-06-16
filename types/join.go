//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package types

type JoinType int

const (
	JOIN_INNER JoinType = iota
	JOIN_OUTER
	JOIN_CROSS
)
