//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/types"
)

type JoinType int

const (
	JOIN_INNER JoinType = iota
	JOIN_OUTER
	JOIN_CROSS
)

type JoinClause struct {
	JoinType JoinType
	left     types.Selection
	right    types.Selection
	on       *ast.Expression
}

func (j *JoinClause) ArgCount() int {
	ac := 0
	if j.on != nil {
		ac = j.on.ArgCount()
	}
	return ac + j.left.ArgCount() + j.right.ArgCount()
}

func (j *JoinClause) Size(scanner types.Scanner) int {
	size := 0
	size += len(scanner.FormatOptions().SeparateClauseWith)
	switch j.JoinType {
	case JOIN_INNER:
		size += len(grammar.Symbols[grammar.SYM_JOIN])
	case JOIN_OUTER:
		size += len(grammar.Symbols[grammar.SYM_LEFT_JOIN])
	case JOIN_CROSS:
		size += len(grammar.Symbols[grammar.SYM_CROSS_JOIN])
		// CROSS JOIN has no ON condition so just short-circuit here
		return size + j.right.Size(scanner)
	}
	size += j.right.Size(scanner)
	size += len(grammar.Symbols[grammar.SYM_ON])
	size += j.on.Size(scanner)
	return size
}

func (j *JoinClause) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], scanner.FormatOptions().SeparateClauseWith)
	switch j.JoinType {
	case JOIN_INNER:
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_JOIN])
	case JOIN_OUTER:
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_LEFT_JOIN])
	case JOIN_CROSS:
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_CROSS_JOIN])
	}
	bw += j.right.Scan(scanner, b[bw:], args, curArg)
	if j.on != nil {
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_ON])
		bw += j.on.Scan(scanner, b[bw:], args, curArg)
	}
	return bw
}

func Join(left types.Selection, right types.Selection, on *ast.Expression) *JoinClause {
	return &JoinClause{left: left, right: right, on: on}
}

func OuterJoin(left types.Selection, right types.Selection, on *ast.Expression) *JoinClause {
	return &JoinClause{
		JoinType: JOIN_OUTER,
		left:     left,
		right:    right,
		on:       on,
	}
}

func CrossJoin(left types.Selection, right types.Selection) *JoinClause {
	return &JoinClause{JoinType: JOIN_CROSS, left: left, right: right}
}
