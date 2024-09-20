//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package inspect

import "github.com/jaypipes/sqlb/grammar"

// ValueSpecificationFromAny evaluates the supplied interface argument and
// returns a *ValueSpecification if the supplied argument can be converted into
// a ValueSpecification, or nil if the conversion cannot be done.
func ValueSpecificationFromAny(
	subject interface{},
) *grammar.ValueSpecification {
	switch v := subject.(type) {
	case *grammar.ValueSpecification:
		return v
	case grammar.ValueSpecification:
		return &v
	case grammar.UnsignedValueSpecification:
		return &grammar.ValueSpecification{
			UnsignedValue: &v,
		}
	case *grammar.UnsignedValueSpecification:
		return &grammar.ValueSpecification{
			UnsignedValue: v,
		}
	case uint, uint8, uint16, uint64:
		return &grammar.ValueSpecification{
			UnsignedValue: &grammar.UnsignedValueSpecification{
				UnsignedLiteral: &grammar.UnsignedLiteral{
					UnsignedNumericLiteral: &grammar.UnsignedNumericLiteral{
						Value: v,
					},
				},
			},
		}
	case int, int8, int16, int64, float64:
		return &grammar.ValueSpecification{
			Literal: &grammar.Literal{
				SignedNumericLiteral: &grammar.SignedNumericLiteral{
					Value: v,
				},
			},
		}
	case string, rune, bool:
		return &grammar.ValueSpecification{
			UnsignedValue: &grammar.UnsignedValueSpecification{
				UnsignedLiteral: &grammar.UnsignedLiteral{
					GeneralLiteral: &grammar.GeneralLiteral{
						Value: v,
					},
				},
			},
		}
	}
	return nil
}

// UnsignedValueSpecificationFromAny evaluates the supplied interface argument
// and returns an *UnsignedValueSpecification if the supplied argument can be
// converted into a UnsignedValueSpecification, or nil if the conversion cannot
// be done.
func UnsignedValueSpecificationFromAny(
	subject interface{},
) *grammar.UnsignedValueSpecification {
	switch v := subject.(type) {
	case *grammar.UnsignedValueSpecification:
		return v
	case grammar.UnsignedValueSpecification:
		return &v
	case grammar.ValueSpecification:
		if v.UnsignedValue != nil {
			return v.UnsignedValue
		}
	case *grammar.ValueSpecification:
		if v.UnsignedValue != nil {
			return v.UnsignedValue
		}
	case uint, uint8, uint16, uint64, int, int8, int16, int64:
		return &grammar.UnsignedValueSpecification{
			UnsignedLiteral: &grammar.UnsignedLiteral{
				UnsignedNumericLiteral: &grammar.UnsignedNumericLiteral{
					Value: v,
				},
			},
		}
	case string, rune, bool:
		return &grammar.UnsignedValueSpecification{
			UnsignedLiteral: &grammar.UnsignedLiteral{
				GeneralLiteral: &grammar.GeneralLiteral{
					Value: v,
				},
			},
		}
	}
	return nil
}
