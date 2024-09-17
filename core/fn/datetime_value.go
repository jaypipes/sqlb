//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package fn

import (
	"github.com/jaypipes/sqlb/core/types"
	"github.com/jaypipes/sqlb/grammar"
)

// CurrentDate returns a CurrentDateFunction that produces a CURRENT_DATE() SQL
// function that can be passed to sqlb constructs and functions like Select()
func CurrentDate() *CurrentDateFunction {
	return &CurrentDateFunction{
		DatetimeValueFunction: &grammar.DatetimeValueFunction{
			CurrentDate: true,
		},
	}
}

// CurrentDateFunction wraps the	CURRENT_DATE() SQL function grammar element
type CurrentDateFunction struct {
	BaseFunction
	*grammar.DatetimeValueFunction
}

// CommonValueExpression returns the object as a
// `*grammar.CommonValueExpression`
func (f *CurrentDateFunction) CommonValueExpression() *grammar.CommonValueExpression {
	return &grammar.CommonValueExpression{
		Datetime: &grammar.DatetimeValueExpression{
			Unary: &grammar.DatetimeTerm{
				Factor: grammar.DatetimeFactor{
					Primary: grammar.DatetimePrimary{
						Function: f.DatetimeValueFunction,
					},
				},
			},
		},
	}
}

// DerivedColumn returns the `*grammar.DerivedColumn` element representing
// the Projection
func (f *CurrentDateFunction) DerivedColumn() *grammar.DerivedColumn {
	dc := &grammar.DerivedColumn{
		ValueExpression: grammar.ValueExpression{
			Common: f.CommonValueExpression(),
		},
	}
	if f.alias != "" {
		dc.As = &f.alias
	}
	return dc
}

// As aliases the SQL function as the supplied column name
func (f *CurrentDateFunction) As(alias string) types.Projection {
	f.alias = alias
	return f
}

// CurrentTime returns a CurrentTimeFunction that produces a CURRENT_TIME() SQL
// function that can be passed to sqlb constructs and functions like Select()
func CurrentTime() *CurrentTimeFunction {
	return &CurrentTimeFunction{
		DatetimeValueFunction: &grammar.DatetimeValueFunction{
			CurrentTime: &grammar.CurrentTimeFunction{},
		},
	}
}

// CurrentTimeFunction wraps the CURRENT_TIME() SQL function grammar element
type CurrentTimeFunction struct {
	BaseFunction
	*grammar.DatetimeValueFunction
}

// Precision sets the function's time or timestamp precision value
func (f *CurrentTimeFunction) Precision(p uint) *CurrentTimeFunction {
	if f.DatetimeValueFunction == nil {
		f.DatetimeValueFunction = &grammar.DatetimeValueFunction{
			CurrentTime: &grammar.CurrentTimeFunction{},
		}
	}
	f.DatetimeValueFunction.CurrentTime.Precision = &p
	return f
}

// CommonValueExpression returns the object as a
// `*grammar.CommonValueExpression`
func (f *CurrentTimeFunction) CommonValueExpression() *grammar.CommonValueExpression {
	return &grammar.CommonValueExpression{
		Datetime: &grammar.DatetimeValueExpression{
			Unary: &grammar.DatetimeTerm{
				Factor: grammar.DatetimeFactor{
					Primary: grammar.DatetimePrimary{
						Function: f.DatetimeValueFunction,
					},
				},
			},
		},
	}
}

// DerivedColumn returns the `*grammar.DerivedColumn` element representing
// the Projection
func (f *CurrentTimeFunction) DerivedColumn() *grammar.DerivedColumn {
	dc := &grammar.DerivedColumn{
		ValueExpression: grammar.ValueExpression{
			Common: f.CommonValueExpression(),
		},
	}
	if f.alias != "" {
		dc.As = &f.alias
	}
	return dc
}

// As aliases the SQL function as the supplied column name
func (f *CurrentTimeFunction) As(alias string) types.Projection {
	f.alias = alias
	return f
}

// CurrentTimestamp returns a CurrentTimestampFunction that produces a
// CURRENT_TIME() SQL function that can be passed to sqlb constructs and
// functions like Select()
func CurrentTimestamp() *CurrentTimestampFunction {
	return &CurrentTimestampFunction{
		DatetimeValueFunction: &grammar.DatetimeValueFunction{
			CurrentTimestamp: &grammar.CurrentTimestampFunction{},
		},
	}
}

// CurrentTimestampFunction wraps the CURRENT_TIMESTAMP() SQL function grammar
// element
type CurrentTimestampFunction struct {
	BaseFunction
	*grammar.DatetimeValueFunction
}

// Precision sets the function's time or timestamp precision value
func (f *CurrentTimestampFunction) Precision(p uint) *CurrentTimestampFunction {
	if f.DatetimeValueFunction == nil {
		f.DatetimeValueFunction = &grammar.DatetimeValueFunction{
			CurrentTimestamp: &grammar.CurrentTimestampFunction{},
		}
	}
	f.DatetimeValueFunction.CurrentTimestamp.Precision = &p
	return f
}

// CommonValueExpression returns the object as a
// `*grammar.CommonValueExpression`
func (f *CurrentTimestampFunction) CommonValueExpression() *grammar.CommonValueExpression {
	return &grammar.CommonValueExpression{
		Datetime: &grammar.DatetimeValueExpression{
			Unary: &grammar.DatetimeTerm{
				Factor: grammar.DatetimeFactor{
					Primary: grammar.DatetimePrimary{
						Function: f.DatetimeValueFunction,
					},
				},
			},
		},
	}
}

// DerivedColumn returns the `*grammar.DerivedColumn` element representing
// the Projection
func (f *CurrentTimestampFunction) DerivedColumn() *grammar.DerivedColumn {
	dc := &grammar.DerivedColumn{
		ValueExpression: grammar.ValueExpression{
			Common: f.CommonValueExpression(),
		},
	}
	if f.alias != "" {
		dc.As = &f.alias
	}
	return dc
}

// As aliases the SQL function as the supplied column name
func (f *CurrentTimestampFunction) As(alias string) types.Projection {
	f.alias = alias
	return f
}

/*
// Now returns a Projection that contains the NOW() SQL function
func Now() api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_NOW),
	}
}

// CurrentTimestamp returns a Projection that contains the CURRENT_TIMESTAMP() SQL function
func CurrentTimestamp() api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_TIMESTAMP),
	}
}

// CurrentTime returns a Projection that contains the CURRENT_TIME() SQL function
func CurrentTime() api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_TIME),
	}
}

// CurrentDate returns a Projection that contains the CURRENT_DATE() SQL function
func CurrentDate() api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CURRENT_DATE),
	}
}
*/