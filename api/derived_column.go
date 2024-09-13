//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import "github.com/jaypipes/sqlb/grammar"

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
	case *grammar.AggregateFunction:
		dc = &grammar.DerivedColumn{
			ValueExpression: grammar.ValueExpression{
				Row: &grammar.RowValueExpression{
					Primary: &grammar.NonParenthesizedValueExpressionPrimary{
						SetFunction: &grammar.SetFunctionSpecification{
							Aggregate: v,
						},
					},
				},
			},
		}
	case *AggregateFunction:
		dc = &grammar.DerivedColumn{
			ValueExpression: grammar.ValueExpression{
				Row: &grammar.RowValueExpression{
					Primary: &grammar.NonParenthesizedValueExpressionPrimary{
						SetFunction: &grammar.SetFunctionSpecification{
							Aggregate: v.AggregateFunction,
						},
					},
				},
			},
		}
	case *SubstringFunction:
		dc = &grammar.DerivedColumn{
			ValueExpression: grammar.ValueExpression{
				Common: &grammar.CommonValueExpression{
					String: &grammar.StringValueExpression{
						Character: &grammar.CharacterValueExpression{
							Factor: &grammar.CharacterFactor{
								Primary: grammar.CharacterPrimary{
									Function: &grammar.StringValueFunction{
										Character: &grammar.CharacterValueFunction{
											Substring: v.SubstringFunction,
										},
									},
								},
							},
						},
					},
				},
			},
		}
	case *grammar.SubstringFunction:
		dc = &grammar.DerivedColumn{
			ValueExpression: grammar.ValueExpression{
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
	case *RegexSubstringFunction:
		dc = &grammar.DerivedColumn{
			ValueExpression: grammar.ValueExpression{
				Common: &grammar.CommonValueExpression{
					String: &grammar.StringValueExpression{
						Character: &grammar.CharacterValueExpression{
							Factor: &grammar.CharacterFactor{
								Primary: grammar.CharacterPrimary{
									Function: &grammar.StringValueFunction{
										Character: &grammar.CharacterValueFunction{
											RegexSubstring: v.RegexSubstringFunction,
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
			ValueExpression: grammar.ValueExpression{
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
	case *FoldFunction:
		dc = &grammar.DerivedColumn{
			ValueExpression: grammar.ValueExpression{
				Common: &grammar.CommonValueExpression{
					String: &grammar.StringValueExpression{
						Character: &grammar.CharacterValueExpression{
							Factor: &grammar.CharacterFactor{
								Primary: grammar.CharacterPrimary{
									Function: &grammar.StringValueFunction{
										Character: &grammar.CharacterValueFunction{
											Fold: v.FoldFunction,
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
			ValueExpression: grammar.ValueExpression{
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
	case *TranscodingFunction:
		dc = &grammar.DerivedColumn{
			ValueExpression: grammar.ValueExpression{
				Common: &grammar.CommonValueExpression{
					String: &grammar.StringValueExpression{
						Character: &grammar.CharacterValueExpression{
							Factor: &grammar.CharacterFactor{
								Primary: grammar.CharacterPrimary{
									Function: &grammar.StringValueFunction{
										Character: &grammar.CharacterValueFunction{
											Transcoding: v.TranscodingFunction,
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
			ValueExpression: grammar.ValueExpression{
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
	}
	if dc != nil {
		if alias != "" {
			dc.As = &alias
		}
	}
	return dc
}
