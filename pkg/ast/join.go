//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ast

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/types"
)

type JoinClause struct {
	joinType types.JoinType
	left     types.Selection
	right    types.Selection
	on       *Expression
}

func (j *JoinClause) Left() types.Selection {
	return j.left
}

func (j *JoinClause) Right() types.Selection {
	return j.right
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
	switch j.joinType {
	case types.JOIN_INNER:
		size += len(grammar.Symbols[grammar.SYM_JOIN])
	case types.JOIN_OUTER:
		size += len(grammar.Symbols[grammar.SYM_LEFT_JOIN])
	case types.JOIN_CROSS:
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
	switch j.joinType {
	case types.JOIN_INNER:
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_JOIN])
	case types.JOIN_OUTER:
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_LEFT_JOIN])
	case types.JOIN_CROSS:
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_CROSS_JOIN])
	}
	bw += j.right.Scan(scanner, b[bw:], args, curArg)
	if j.on != nil {
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_ON])
		bw += j.on.Scan(scanner, b[bw:], args, curArg)
	}
	return bw
}

func Join(left types.Selection, right types.Selection, on *Expression) *JoinClause {
	return &JoinClause{
		joinType: types.JOIN_INNER,
		left:     left,
		right:    right,
		on:       on,
	}
}

func OuterJoin(left types.Selection, right types.Selection, on *Expression) *JoinClause {
	return &JoinClause{
		joinType: types.JOIN_OUTER,
		left:     left,
		right:    right,
		on:       on,
	}
}

func CrossJoin(left types.Selection, right types.Selection) *JoinClause {
	return &JoinClause{
		joinType: types.JOIN_CROSS,
		left:     left,
		right:    right,
	}
}

func NewJoinClause(
	jt types.JoinType,
	left types.Selection,
	right types.Selection,
	on *Expression,
) *JoinClause {
	switch jt {
	case types.JOIN_INNER:
		return Join(left, right, on)
	case types.JOIN_OUTER:
		return OuterJoin(left, right, on)
	case types.JOIN_CROSS:
		return CrossJoin(left, right)
	}
	return nil
}
