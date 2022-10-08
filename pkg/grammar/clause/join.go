//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
	"strings"

	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/grammar/expression"
	"github.com/jaypipes/sqlb/pkg/types"
)

type Join struct {
	joinType types.JoinType
	left     types.Selection
	right    types.Selection
	on       *expression.Expression
}

func (j *Join) Left() types.Selection {
	return j.left
}

func (j *Join) Right() types.Selection {
	return j.right
}

func (j *Join) ArgCount() int {
	ac := 0
	if j.on != nil {
		ac = j.on.ArgCount()
	}
	return ac + j.left.ArgCount() + j.right.ArgCount()
}

func (j *Join) Size(scanner types.Scanner) int {
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

func (j *Join) Scan(scanner types.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	b.WriteString(scanner.FormatOptions().SeparateClauseWith)
	switch j.joinType {
	case types.JOIN_INNER:
		b.Write(grammar.Symbols[grammar.SYM_JOIN])
	case types.JOIN_OUTER:
		b.Write(grammar.Symbols[grammar.SYM_LEFT_JOIN])
	case types.JOIN_CROSS:
		b.Write(grammar.Symbols[grammar.SYM_CROSS_JOIN])
	}
	j.right.Scan(scanner, b, args, curArg)
	if j.on != nil {
		b.Write(grammar.Symbols[grammar.SYM_ON])
		j.on.Scan(scanner, b, args, curArg)
	}
}

func InnerJoin(left types.Selection, right types.Selection, on *expression.Expression) *Join {
	return &Join{
		joinType: types.JOIN_INNER,
		left:     left,
		right:    right,
		on:       on,
	}
}

func OuterJoin(left types.Selection, right types.Selection, on *expression.Expression) *Join {
	return &Join{
		joinType: types.JOIN_OUTER,
		left:     left,
		right:    right,
		on:       on,
	}
}

func CrossJoin(left types.Selection, right types.Selection) *Join {
	return &Join{
		joinType: types.JOIN_CROSS,
		left:     left,
		right:    right,
	}
}

func NewJoin(
	jt types.JoinType,
	left types.Selection,
	right types.Selection,
	on *expression.Expression,
) *Join {
	switch jt {
	case types.JOIN_INNER:
		return InnerJoin(left, right, on)
	case types.JOIN_OUTER:
		return OuterJoin(left, right, on)
	case types.JOIN_CROSS:
		return CrossJoin(left, right)
	}
	return nil
}
