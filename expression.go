package sqlb

type exprType int

const (
	EXP_EQUAL = iota
	EXP_NEQUAL
	EXP_AND
	EXP_OR
	EXP_IN
	EXP_BETWEEN
	EXP_IS_NULL
	EXP_IS_NOT_NULL
	EXP_GREATER
	EXP_GREATER_EQUAL
	EXP_LESS
	EXP_LESS_EQUAL
)

var (
	// A static table containing information used in constructing the
	// expression's SQL string during scan() calls
	exprScanTable = map[exprType]scanInfo{
		EXP_EQUAL: scanInfo{
			SYM_ELEMENT, SYM_EQUAL, SYM_ELEMENT,
		},
		EXP_NEQUAL: scanInfo{
			SYM_ELEMENT, SYM_NEQUAL, SYM_ELEMENT,
		},
		EXP_AND: scanInfo{
			SYM_LPAREN, SYM_ELEMENT, SYM_AND, SYM_ELEMENT, SYM_RPAREN,
		},
		EXP_OR: scanInfo{
			SYM_LPAREN, SYM_ELEMENT, SYM_OR, SYM_ELEMENT, SYM_RPAREN,
		},
		EXP_IN: scanInfo{
			SYM_ELEMENT, SYM_IN, SYM_ELEMENT, SYM_RPAREN,
		},
		EXP_BETWEEN: scanInfo{
			SYM_ELEMENT, SYM_BETWEEN, SYM_ELEMENT, SYM_AND, SYM_ELEMENT,
		},
		EXP_IS_NULL: scanInfo{
			SYM_ELEMENT, SYM_IS_NULL,
		},
		EXP_IS_NOT_NULL: scanInfo{
			SYM_ELEMENT, SYM_IS_NOT_NULL,
		},
		EXP_GREATER: scanInfo{
			SYM_ELEMENT, SYM_GREATER, SYM_ELEMENT,
		},
		EXP_GREATER_EQUAL: scanInfo{
			SYM_ELEMENT, SYM_GREATER_EQUAL, SYM_ELEMENT,
		},
		EXP_LESS: scanInfo{
			SYM_ELEMENT, SYM_LESS, SYM_ELEMENT,
		},
		EXP_LESS_EQUAL: scanInfo{
			SYM_ELEMENT, SYM_LESS_EQUAL, SYM_ELEMENT,
		},
	}
)

type Expression struct {
	scanInfo scanInfo
	elements []element
}

func (e *Expression) referrents() []selection {
	res := make([]selection, 0)
	for _, el := range e.elements {
		switch el.(type) {
		case projection:
			p := el.(projection)
			res = append(res, p.from())
		}
	}
	return res
}

func (e *Expression) argCount() int {
	ac := 0
	for _, el := range e.elements {
		ac += el.argCount()
	}
	return ac
}

func (e *Expression) size(scanner *sqlScanner) int {
	size := 0
	elidx := 0
	for _, sym := range e.scanInfo {
		if sym == SYM_ELEMENT {
			el := e.elements[elidx]
			// We need to disable alias output for elements that are
			// projections. We don't want to output, for example,
			// "ON users.id AS user_id = articles.author"
			switch el.(type) {
			case projection:
				reset := el.(projection).disableAliasScan()
				defer reset()
			}
			elidx++
			size += el.size(scanner)
		} else {
			size += len(Symbols[sym])
		}
	}
	return size
}

func (e *Expression) scan(scanner *sqlScanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	elidx := 0
	for _, sym := range e.scanInfo {
		if sym == SYM_ELEMENT {
			el := e.elements[elidx]
			// We need to disable alias output for elements that are
			// projections. We don't want to output, for example,
			// "ON users.id AS user_id = articles.author"
			switch el.(type) {
			case projection:
				reset := el.(projection).disableAliasScan()
				defer reset()
			}
			elidx++
			bw += el.scan(scanner, b[bw:], args, curArg)
		} else {
			bw += copy(b[bw:], Symbols[sym])
		}
	}
	return bw
}

func Equal(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		scanInfo: exprScanTable[EXP_EQUAL],
		elements: els,
	}
}

func NotEqual(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		scanInfo: exprScanTable[EXP_NEQUAL],
		elements: els,
	}
}

func And(a *Expression, b *Expression) *Expression {
	return &Expression{
		scanInfo: exprScanTable[EXP_AND],
		elements: []element{a, b},
	}
}

func Or(a *Expression, b *Expression) *Expression {
	return &Expression{
		scanInfo: exprScanTable[EXP_OR],
		elements: []element{a, b},
	}
}

func In(subject element, values ...interface{}) *Expression {
	return &Expression{
		scanInfo: exprScanTable[EXP_IN],
		elements: []element{subject, toValueList(values...)},
	}
}

func Between(subject element, start interface{}, end interface{}) *Expression {
	els := toElements(subject, start, end)
	return &Expression{
		scanInfo: exprScanTable[EXP_BETWEEN],
		elements: els,
	}
}

func IsNull(subject element) *Expression {
	return &Expression{
		scanInfo: exprScanTable[EXP_IS_NULL],
		elements: []element{subject},
	}
}

func IsNotNull(subject element) *Expression {
	return &Expression{
		scanInfo: exprScanTable[EXP_IS_NOT_NULL],
		elements: []element{subject},
	}
}

func GreaterThan(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		scanInfo: exprScanTable[EXP_GREATER],
		elements: els,
	}
}

func GreaterThanOrEqual(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		scanInfo: exprScanTable[EXP_GREATER_EQUAL],
		elements: els,
	}
}

func LessThan(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		scanInfo: exprScanTable[EXP_LESS],
		elements: els,
	}
}

func LessThanOrEqual(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		scanInfo: exprScanTable[EXP_LESS_EQUAL],
		elements: els,
	}
}
