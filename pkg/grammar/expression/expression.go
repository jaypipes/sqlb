//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package expression

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/grammar/element"
	"github.com/jaypipes/sqlb/pkg/types"
)

// Expression represents a comparison expression in the SQL statement, e.g. "a
// = b"
type Expression struct {
	ScanInfo grammar.ScanInfo
	elements []types.Element
}

// Elements returns the expression's list of contained elements
func (e *Expression) Elements() []types.Element {
	return e.elements
}

func (e *Expression) Referrents() []types.Selection {
	res := make([]types.Selection, 0)
	for _, el := range e.elements {
		switch el.(type) {
		case types.Projection:
			p := el.(types.Projection)
			res = append(res, p.From())
		}
	}
	return res
}

func (e *Expression) ArgCount() int {
	ac := 0
	for _, el := range e.elements {
		ac += el.ArgCount()
	}
	return ac
}

func (e *Expression) Size(scanner types.Scanner) int {
	size := 0
	elidx := 0
	for _, sym := range e.ScanInfo {
		if sym == grammar.SYM_ELEMENT {
			el := e.elements[elidx]
			// We need to disable alias output for elements that are
			// projections. We don't want to output, for example,
			// "ON users.id AS user_id = articles.author"
			switch el.(type) {
			case types.Projection:
				reset := el.(types.Projection).DisableAliasScan()
				defer reset()
			}
			elidx++
			size += el.Size(scanner)
		} else {
			size += len(grammar.Symbols[sym])
		}
	}
	return size
}

func (e *Expression) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	elidx := 0
	for _, sym := range e.ScanInfo {
		if sym == grammar.SYM_ELEMENT {
			el := e.elements[elidx]
			// We need to disable alias output for elements that are
			// projections. We don't want to output, for example,
			// "ON users.id AS user_id = articles.author"
			switch el.(type) {
			case types.Projection:
				reset := el.(types.Projection).DisableAliasScan()
				defer reset()
			}
			elidx++
			bw += el.Scan(scanner, b[bw:], args, curArg)
		} else {
			bw += copy(b[bw:], grammar.Symbols[sym])
		}
	}
	return bw
}

// Given a slice of interface{} variables, returns a slice of element members.
// If any of the interface{} variables are *not* of type element already, we
// construct a Value{} for the variable.
func toElements(vars ...interface{}) []types.Element {
	els := make([]types.Element, len(vars))
	for x, v := range vars {
		switch v.(type) {
		case types.Element:
			els[x] = v.(types.Element)
		default:
			els[x] = element.NewValue(nil, v)
		}
	}
	return els
}

// Given a variable number of interface{} variables, returns a List containing
// Value structs for the variables
// If any of the interface{} variables are *not* of type element already, we
// construct a Value{} for the variable.
func toValueList(vars ...interface{}) *element.List {
	els := make([]types.Element, len(vars))
	for x, v := range vars {
		els[x] = element.NewValue(nil, v)
	}
	return element.NewList(els...)
}

func Equal(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_EQUAL),
		elements: els,
	}
}

func NotEqual(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_NEQUAL),
		elements: els,
	}
}

func And(a *Expression, b *Expression) *Expression {
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_AND),
		elements: []types.Element{a, b},
	}
}

func Or(a *Expression, b *Expression) *Expression {
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_OR),
		elements: []types.Element{a, b},
	}
}

func In(subject types.Element, values ...interface{}) *Expression {
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_IN),
		elements: []types.Element{subject, toValueList(values...)},
	}
}

func Between(subject types.Element, start interface{}, end interface{}) *Expression {
	els := toElements(subject, start, end)
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_BETWEEN),
		elements: els,
	}
}

func IsNull(subject types.Element) *Expression {
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_IS_NULL),
		elements: []types.Element{subject},
	}
}

func IsNotNull(subject types.Element) *Expression {
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_IS_NOT_NULL),
		elements: []types.Element{subject},
	}
}

func GreaterThan(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_GREATER),
		elements: els,
	}
}

func GreaterThanOrEqual(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_GREATER_EQUAL),
		elements: els,
	}
}

func LessThan(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_LESS),
		elements: els,
	}
}

func LessThanOrEqual(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_LESS_EQUAL),
		elements: els,
	}
}
