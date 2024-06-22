//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

type JoinType int

const (
	JoinInner JoinType = iota
	JoinOuter
	JoinCross
)
