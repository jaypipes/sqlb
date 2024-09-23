//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package inspect

import (
	"github.com/jaypipes/sqlb/core/grammar"
	"github.com/jaypipes/sqlb/core/types"
)

// DerivedColumnFromAnyAndAlias returns a DerivedColumn from something that can
// be coerced into a ValueExpression and an optional alias.
//
// If the first argument cannot be coerced into a ValueExpression, returns nil.
func DerivedColumnFromAnyAndAlias(
	subject interface{},
	alias string,
) *grammar.DerivedColumn {
	var dc *grammar.DerivedColumn
	switch v := subject.(type) {
	case *grammar.DerivedColumn:
		dc = v
	case grammar.DerivedColumn:
		dc = &v
	case types.Projection:
		dc = v.DerivedColumn()
	case *grammar.AggregateFunction:
		dc = &grammar.DerivedColumn{
			Value: grammar.ValueExpression{
				Row: &grammar.RowValueExpression{
					Primary: &grammar.NonParenthesizedValueExpressionPrimary{
						SetFunction: &grammar.SetFunctionSpecification{
							Aggregate: v,
						},
					},
				},
			},
		}
	case *grammar.SubstringFunction:
		dc = &grammar.DerivedColumn{
			Value: grammar.ValueExpression{
				Common: &grammar.CommonValueExpression{
					String: &grammar.StringValueExpression{
						Character: &grammar.CharacterValueExpression{
							Factor: &grammar.CharacterFactor{
								Primary: grammar.CharacterPrimary{
									Function: &grammar.StringValueFunction{
										Character: &grammar.CharacterValueFunction{
											Substring: v,
										},
									},
								},
							},
						},
					},
				},
			},
		}
	case *grammar.RegexSubstringFunction:
		dc = &grammar.DerivedColumn{
			Value: grammar.ValueExpression{
				Common: &grammar.CommonValueExpression{
					String: &grammar.StringValueExpression{
						Character: &grammar.CharacterValueExpression{
							Factor: &grammar.CharacterFactor{
								Primary: grammar.CharacterPrimary{
									Function: &grammar.StringValueFunction{
										Character: &grammar.CharacterValueFunction{
											RegexSubstring: v,
										},
									},
								},
							},
						},
					},
				},
			},
		}
	case *grammar.FoldFunction:
		dc = &grammar.DerivedColumn{
			Value: grammar.ValueExpression{
				Common: &grammar.CommonValueExpression{
					String: &grammar.StringValueExpression{
						Character: &grammar.CharacterValueExpression{
							Factor: &grammar.CharacterFactor{
								Primary: grammar.CharacterPrimary{
									Function: &grammar.StringValueFunction{
										Character: &grammar.CharacterValueFunction{
											Fold: v,
										},
									},
								},
							},
						},
					},
				},
			},
		}
	case *grammar.TranscodingFunction:
		dc = &grammar.DerivedColumn{
			Value: grammar.ValueExpression{
				Common: &grammar.CommonValueExpression{
					String: &grammar.StringValueExpression{
						Character: &grammar.CharacterValueExpression{
							Factor: &grammar.CharacterFactor{
								Primary: grammar.CharacterPrimary{
									Function: &grammar.StringValueFunction{
										Character: &grammar.CharacterValueFunction{
											Transcoding: v,
										},
									},
								},
							},
						},
					},
				},
			},
		}
	case *grammar.TransliterationFunction:
		dc = &grammar.DerivedColumn{
			Value: grammar.ValueExpression{
				Common: &grammar.CommonValueExpression{
					String: &grammar.StringValueExpression{
						Character: &grammar.CharacterValueExpression{
							Factor: &grammar.CharacterFactor{
								Primary: grammar.CharacterPrimary{
									Function: &grammar.StringValueFunction{
										Character: &grammar.CharacterValueFunction{
											Transliteration: v,
										},
									},
								},
							},
						},
					},
				},
			},
		}
	case *grammar.TrimFunction:
		dc = &grammar.DerivedColumn{
			Value: grammar.ValueExpression{
				Common: &grammar.CommonValueExpression{
					String: &grammar.StringValueExpression{
						Character: &grammar.CharacterValueExpression{
							Factor: &grammar.CharacterFactor{
								Primary: grammar.CharacterPrimary{
									Function: &grammar.StringValueFunction{
										Character: &grammar.CharacterValueFunction{
											Trim: v,
										},
									},
								},
							},
						},
					},
				},
			},
		}
	}
	if dc != nil {
		if alias != "" {
			dc.As = &alias
		}
	}
	return dc
}
