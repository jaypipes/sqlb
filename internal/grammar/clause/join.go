//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
)

type Join struct {
	joinType api.JoinType
	left     builder.Selection
	right    builder.Selection
	on       *expression.Expression
}

func (j *Join) Left() builder.Selection {
	return j.left
}

func (j *Join) Right() builder.Selection {
	return j.right
}

func (j *Join) ArgCount() int {
	ac := 0
	if j.on != nil {
		ac = j.on.ArgCount()
	}
	return ac + j.left.ArgCount() + j.right.ArgCount()
}

func (j *Join) Size(b *builder.Builder) int {
	size := 0
	size += len(b.Format.SeparateClauseWith)
	switch j.joinType {
	case api.JoinInner:
		size += len(grammar.Symbols[grammar.SYM_JOIN])
	case api.JoinOuter:
		size += len(grammar.Symbols[grammar.SYM_LEFT_JOIN])
	case api.JoinCross:
		size += len(grammar.Symbols[grammar.SYM_CROSS_JOIN])
		// CROSS JOIN has no ON condition so just short-circuit here
		return size + j.right.Size(b)
	}
	size += j.right.Size(b)
	size += len(grammar.Symbols[grammar.SYM_ON])
	size += j.on.Size(b)
	return size
}

func (j *Join) Scan(b *builder.Builder, args []interface{}, curArg *int) {
	b.WriteString(b.Format.SeparateClauseWith)
	switch j.joinType {
	case api.JoinInner:
		b.Write(grammar.Symbols[grammar.SYM_JOIN])
	case api.JoinOuter:
		b.Write(grammar.Symbols[grammar.SYM_LEFT_JOIN])
	case api.JoinCross:
		b.Write(grammar.Symbols[grammar.SYM_CROSS_JOIN])
	}
	j.right.Scan(b, args, curArg)
	if j.on != nil {
		b.Write(grammar.Symbols[grammar.SYM_ON])
		j.on.Scan(b, args, curArg)
	}
}

func InnerJoin(left builder.Selection, right builder.Selection, on *expression.Expression) *Join {
	return &Join{
		joinType: api.JoinInner,
		left:     left,
		right:    right,
		on:       on,
	}
}

func OuterJoin(left builder.Selection, right builder.Selection, on *expression.Expression) *Join {
	return &Join{
		joinType: api.JoinOuter,
		left:     left,
		right:    right,
		on:       on,
	}
}

func CrossJoin(left builder.Selection, right builder.Selection) *Join {
	return &Join{
		joinType: api.JoinCross,
		left:     left,
		right:    right,
	}
}

func NewJoin(
	jt api.JoinType,
	left builder.Selection,
	right builder.Selection,
	on *expression.Expression,
) *Join {
	switch jt {
	case api.JoinInner:
		return InnerJoin(left, right, on)
	case api.JoinOuter:
		return OuterJoin(left, right, on)
	case api.JoinCross:
		return CrossJoin(left, right)
	}
	return nil
}
