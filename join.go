package sqlb

type joinType int

const (
	JOIN_INNER joinType = iota
	JOIN_OUTER
	JOIN_CROSS
)

type joinClause struct {
	joinType joinType
	left     selection
	right    selection
	on       *Expression
}

func (j *joinClause) argCount() int {
	ac := 0
	if j.on != nil {
		ac = j.on.argCount()
	}
	return ac + j.left.argCount() + j.right.argCount()
}

func (j *joinClause) size() int {
	size := 0
	switch j.joinType {
	case JOIN_INNER:
		size += len(Symbols[SYM_JOIN])
	case JOIN_OUTER:
		size += len(Symbols[SYM_LEFT_JOIN])
	case JOIN_CROSS:
		size += len(Symbols[SYM_CROSS_JOIN])
		// CROSS JOIN has no ON condition so just short-circuit here
		return size + j.right.size()
	}
	size += j.right.size()
	size += len(Symbols[SYM_ON])
	size += j.on.size()
	return size
}

func (j *joinClause) scan(b []byte, args []interface{}, curArg *int) int {
	bw := 0
	switch j.joinType {
	case JOIN_INNER:
		bw += copy(b[bw:], Symbols[SYM_JOIN])
	case JOIN_OUTER:
		bw += copy(b[bw:], Symbols[SYM_LEFT_JOIN])
	case JOIN_CROSS:
		bw += copy(b[bw:], Symbols[SYM_CROSS_JOIN])
	}
	bw += j.right.scan(b[bw:], args, curArg)
	if j.on != nil {
		bw += copy(b[bw:], Symbols[SYM_ON])
		bw += j.on.scan(b[bw:], args, curArg)
	}
	return bw
}

func Join(left selection, right selection, on *Expression) *joinClause {
	return &joinClause{left: left, right: right, on: on}
}

func OuterJoin(left selection, right selection, on *Expression) *joinClause {
	return &joinClause{
		joinType: JOIN_OUTER,
		left:     left,
		right:    right,
		on:       on,
	}
}

func CrossJoin(left selection, right selection) *joinClause {
	return &joinClause{joinType: JOIN_CROSS, left: left, right: right}
}
