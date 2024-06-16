//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package expression

import (
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/element"
)

// Expression represents a comparison expression in the SQL statement, e.g. "a
// = b"
type Expression struct {
	ScanInfo grammar.ScanInfo
	elements []builder.Element
}

// Elements returns the expression's list of contained elements
func (e *Expression) Elements() []builder.Element {
	return e.elements
}

func (e *Expression) Referrents() []builder.Selection {
	res := make([]builder.Selection, 0)
	for _, el := range e.elements {
		switch el.(type) {
		case builder.Projection:
			p := el.(builder.Projection)
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

func (e *Expression) Size(b *builder.Builder) int {
	size := 0
	elidx := 0
	for _, sym := range e.ScanInfo {
		if sym == grammar.SYM_ELEMENT {
			el := e.elements[elidx]
			// We need to disable alias output for elements that are
			// projections. We don't want to output, for example,
			// "ON users.id AS user_id = articles.author"
			switch el.(type) {
			case builder.Projection:
				reset := el.(builder.Projection).DisableAliasScan()
				defer reset()
			}
			elidx++
			size += el.Size(b)
		} else {
			size += len(grammar.Symbols[sym])
		}
	}
	return size
}

func (e *Expression) Scan(b *builder.Builder, args []interface{}, curArg *int) {
	elidx := 0
	for _, sym := range e.ScanInfo {
		if sym == grammar.SYM_ELEMENT {
			el := e.elements[elidx]
			// We need to disable alias output for elements that are
			// projections. We don't want to output, for example,
			// "ON users.id AS user_id = articles.author"
			switch el.(type) {
			case builder.Projection:
				reset := el.(builder.Projection).DisableAliasScan()
				defer reset()
			}
			elidx++
			el.Scan(b, args, curArg)
		} else {
			b.Write(grammar.Symbols[sym])
		}
	}
}

// Given a slice of interface{} variables, returns a slice of element members.
// If any of the interface{} variables are *not* of type element already, we
// construct a Value{} for the variable.
func toElements(vars ...interface{}) []builder.Element {
	els := make([]builder.Element, len(vars))
	for x, v := range vars {
		switch v.(type) {
		case builder.Element:
			els[x] = v.(builder.Element)
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
	els := make([]builder.Element, len(vars))
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
		elements: []builder.Element{a, b},
	}
}

func Or(a *Expression, b *Expression) *Expression {
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_OR),
		elements: []builder.Element{a, b},
	}
}

func In(subject builder.Element, values ...interface{}) *Expression {
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_IN),
		elements: []builder.Element{subject, toValueList(values...)},
	}
}

func Between(subject builder.Element, start interface{}, end interface{}) *Expression {
	els := toElements(subject, start, end)
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_BETWEEN),
		elements: els,
	}
}

func IsNull(subject builder.Element) *Expression {
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_IS_NULL),
		elements: []builder.Element{subject},
	}
}

func IsNotNull(subject builder.Element) *Expression {
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_IS_NOT_NULL),
		elements: []builder.Element{subject},
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
