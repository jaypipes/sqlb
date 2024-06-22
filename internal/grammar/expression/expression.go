//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package expression

import (
	"strings"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/element"
)

// Expression represents a comparison expression in the SQL statement, e.g. "a
// = b"
type Expression struct {
	ScanInfo grammar.ScanInfo
	elements []api.Element
}

// Elements returns the expression's list of contained elements
func (e *Expression) Elements() []api.Element {
	return e.elements
}

func (e *Expression) Referrents() []api.Selection {
	res := make([]api.Selection, 0)
	for _, el := range e.elements {
		switch el.(type) {
		case api.Projection:
			p := el.(api.Projection)
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

func (e *Expression) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	elidx := 0
	for _, sym := range e.ScanInfo {
		if sym == grammar.SYM_ELEMENT {
			el := e.elements[elidx]
			// We need to disable alias output for elements that are
			// projections. We don't want to output, for example,
			// "ON users.id AS user_id = articles.author"
			switch el.(type) {
			case api.Projection:
				reset := el.(api.Projection).DisableAliasScan()
				defer reset()
			}
			elidx++
			b.WriteString(el.String(opts, qargs, curarg))
		} else {
			b.Write(grammar.Symbols[sym])
		}
	}
	return b.String()
}

// Given a slice of interface{} variables, returns a slice of element members.
// If any of the interface{} variables are *not* of type element already, we
// construct a Value{} for the variable.
func toElements(vars ...interface{}) []api.Element {
	els := make([]api.Element, len(vars))
	for x, v := range vars {
		switch v.(type) {
		case api.Element:
			els[x] = v.(api.Element)
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
	els := make([]api.Element, len(vars))
	for x, v := range vars {
		els[x] = element.NewValue(nil, v)
	}
	return element.NewList(els...)
}

// Equal accepts two things and returns an Element representing an equality
// expression that can be passed to a Join or Where clause.
func Equal(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_EQUAL),
		elements: els,
	}
}

// NotEqual accepts two things and returns an Element representing an
// inequality expression that can be passed to a Join or Where clause.
func NotEqual(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_NEQUAL),
		elements: els,
	}
}

// And accepts two things and returns an Element representing an AND expression
// that can be passed to a Join or Where clause.
func And(a *Expression, b *Expression) *Expression {
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_AND),
		elements: []api.Element{a, b},
	}
}

// Or accepts two things and returns an Element representing an OR expression
// that can be passed to a Join or Where clause.
func Or(a *Expression, b *Expression) *Expression {
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_OR),
		elements: []api.Element{a, b},
	}
}

// In accepts two things and returns an Element representing an IN expression
// that can be passed to a Join or Where clause.
func In(subject api.Element, values ...interface{}) *Expression {
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_IN),
		elements: []api.Element{subject, toValueList(values...)},
	}
}

// Between accepts an element and a start and end things and returns an Element
// representing a BETWEEN expression that can be passed to a Join or Where
// clause.
func Between(subject api.Element, start interface{}, end interface{}) *Expression {
	els := toElements(subject, start, end)
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_BETWEEN),
		elements: els,
	}
}

// IsNull accepts an element and returns an Element representing an IS NULL
// expression that can be passed to a Join or Where clause.
func IsNull(subject api.Element) *Expression {
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_IS_NULL),
		elements: []api.Element{subject},
	}
}

// IsNotNull accepts an element and returns an Element representing an IS NOT
// NULL expression that can be passed to a Join or Where clause.
func IsNotNull(subject api.Element) *Expression {
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_IS_NOT_NULL),
		elements: []api.Element{subject},
	}
}

// GreaterThan accepts two things and returns an Element representing a greater
// than expression that can be passed to a Join or Where clause.
func GreaterThan(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_GREATER),
		elements: els,
	}
}

// GreaterThanOrEqual accepts two things and returns an Element representing a
// greater than or equality expression that can be passed to a Join or Where
// clause.
func GreaterThanOrEqual(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_GREATER_EQUAL),
		elements: els,
	}
}

// LessThan accepts two things and returns an Element representing a less than
// expression that can be passed to a Join or Where clause.
func LessThan(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_LESS),
		elements: els,
	}
}

// LessThanOrEqual accepts two things and returns an Element representing a
// less than or equality expression that can be passed to a Join or Where
// clause.
func LessThanOrEqual(left interface{}, right interface{}) *Expression {
	els := toElements(left, right)
	return &Expression{
		ScanInfo: grammar.ExpressionScanTable(grammar.EXP_LESS_EQUAL),
		elements: els,
	}
}
