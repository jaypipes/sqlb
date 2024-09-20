//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package fn

import (
	"fmt"
	"strings"

	"github.com/jaypipes/sqlb/core/types"
	"github.com/jaypipes/sqlb/grammar"
	"github.com/jaypipes/sqlb/internal/inspect"
)

// Substring returns a SubstringFunction that produces a SUBSTRING() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The first argument is the subject of the SUBSTRING function and must be
// coercible to a character value expression. The second argument is the FROM
// portion of the SUBSTRING function, which is the index in the subject from
// which to return a substring. The second argument must be coercible to a
// numeric value expression.
func Substring(
	subjectAny interface{},
	fromAny interface{},
) *SubstringFunction {
	var ref types.Relation
	switch subjectAny := subjectAny.(type) {
	case types.Projection:
		ref = subjectAny.References()
	}
	subject := inspect.CharacterValueExpressionFromAny(subjectAny)
	if subject == nil {
		msg := fmt.Sprintf(
			"expected coerceable CharacterValueExpression but got %+v(%T)",
			subjectAny, subjectAny,
		)
		panic(msg)
	}
	from := inspect.NumericValueExpressionFromAny(fromAny)
	if from == nil {
		msg := fmt.Sprintf(
			"expected coerceable NumericValueExpression but got %+v(%T)",
			fromAny, fromAny,
		)
		panic(msg)
	}
	return &SubstringFunction{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		SubstringFunction: &grammar.SubstringFunction{
			Subject: *subject,
			From:    *from,
		},
	}
}

// SubstringFunction wraps the SUBSTRING() SQL function grammar element
type SubstringFunction struct {
	BaseFunction
	*grammar.SubstringFunction
}

// CommonValueExpression returns the object as a
// `*grammar.CommonValueExpression`
func (f *SubstringFunction) CommonValueExpression() *grammar.CommonValueExpression {
	return &grammar.CommonValueExpression{
		String: &grammar.StringValueExpression{
			Character: &grammar.CharacterValueExpression{
				Factor: &grammar.CharacterFactor{
					Primary: grammar.CharacterPrimary{
						Function: &grammar.StringValueFunction{
							Character: &grammar.CharacterValueFunction{
								Substring: f.SubstringFunction,
							},
						},
					},
				},
			},
		},
	}
}

// DerivedColumn returns the `*grammar.DerivedColumn` element representing
// the Projection
func (f *SubstringFunction) DerivedColumn() *grammar.DerivedColumn {
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
func (f *SubstringFunction) As(alias string) types.Projection {
	f.alias = alias
	return f
}

// Using modifies the SUBSTRING function with a character length units.
func (f *SubstringFunction) Using(
	using grammar.CharacterLengthUnits,
) *SubstringFunction {
	f.SubstringFunction.Using = using
	return f
}

// For modifies the SUBSTRING function with a string length. The supplied
// argument must be coercible into a Numeric Value Expression.
func (f *SubstringFunction) For(valAny interface{}) *SubstringFunction {
	v := inspect.NumericValueExpressionFromAny(valAny)
	if v == nil {
		msg := fmt.Sprintf(
			"expected coerceable NumericValueExpression but got %+v(%T)",
			valAny, valAny,
		)
		panic(msg)
	}
	f.SubstringFunction.For = v
	return f
}

// RegexSubstring returns a RegexSubstringFunction that produces a SUBSTRING()
// SQL function of the Regular Expression subtype that can be passed to sqlb
// constructs and functions like Select()
//
// The first argument is the subject of the SUBSTRING function and must be
// coercible to a character value expression. The second argument is the
// SIMILAR portion of the SUBSTRING function, which is the regular expression
// pattern to evaluate against the subject. The second argument must be
// coercible to a character value expression. The third argument is the ESCAPE
// portion of the SUBSTRING function, which is the characters that should be
// used as an escape sequence for the regular expression. The third argument
// must be coercible to a character value expression.
func RegexSubstring(
	subjectAny interface{},
	similarAny interface{},
	escapeAny interface{},
) *RegexSubstringFunction {
	var ref types.Relation
	switch subjectAny := subjectAny.(type) {
	case types.Projection:
		ref = subjectAny.References()
	}
	subject := inspect.CharacterValueExpressionFromAny(subjectAny)
	if subject == nil {
		msg := fmt.Sprintf(
			"expected coerceable CharacterValueExpression for "+
				"subject argument but got %+v(%T)",
			subjectAny, subjectAny,
		)
		panic(msg)
	}
	similar := inspect.CharacterValueExpressionFromAny(similarAny)
	if similar == nil {
		msg := fmt.Sprintf(
			"expected coerceable CharacterValueExpression for "+
				"similar argument but got %+v(%T)",
			similarAny, similarAny,
		)
		panic(msg)
	}
	escape := inspect.CharacterValueExpressionFromAny(escapeAny)
	if escape == nil {
		msg := fmt.Sprintf(
			"expected coerceable CharacterValueExpression for "+
				"escape argument but got %+v(%T)",
			escapeAny, escapeAny,
		)
		panic(msg)
	}
	return &RegexSubstringFunction{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		RegexSubstringFunction: &grammar.RegexSubstringFunction{
			Subject: *subject,
			Similar: *similar,
			Escape:  *escape,
		},
	}
}

// RegexSubstringFunction wraps the SUBSTRING() SQL function with a regular
// expression matching variant grammar element
type RegexSubstringFunction struct {
	BaseFunction
	*grammar.RegexSubstringFunction
}

// CommonValueExpression returns the object as a
// `*grammar.CommonValueExpression`
func (f *RegexSubstringFunction) CommonValueExpression() *grammar.CommonValueExpression {
	return &grammar.CommonValueExpression{
		String: &grammar.StringValueExpression{
			Character: &grammar.CharacterValueExpression{
				Factor: &grammar.CharacterFactor{
					Primary: grammar.CharacterPrimary{
						Function: &grammar.StringValueFunction{
							Character: &grammar.CharacterValueFunction{
								RegexSubstring: f.RegexSubstringFunction,
							},
						},
					},
				},
			},
		},
	}
}

// DerivedColumn returns the `*grammar.DerivedColumn` element representing
// the Projection
func (f *RegexSubstringFunction) DerivedColumn() *grammar.DerivedColumn {
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
func (f *RegexSubstringFunction) As(alias string) types.Projection {
	f.alias = alias
	return f
}

// Upper returns a FoldFunction that produces an UPPER() SQL function that can
// be passed to sqlb constructs and functions like Select()
//
// The only argument is the subject of the UPPER function and must be coercible
// to a character value expression.
func Upper(
	subjectAny interface{},
) *FoldFunction {
	return Fold(subjectAny, grammar.FoldCaseUpper)
}

// Lower returns a FoldFunction that produces a LOWER() SQL function that can
// be passed to sqlb constructs and functions like Select()
//
// The only argument is the subject of the LOWER function and must be coercible
// to a character value expression.
func Lower(
	subjectAny interface{},
) *FoldFunction {
	return Fold(subjectAny, grammar.FoldCaseLower)
}

// Fold returns a FoldFunction that produces an UPPER() or LOWER() SQL function
// that can be passed to sqlb constructs and functions like Select()
//
// The first argument is the subject of the UPPER or LOWER function and must be
// coercible to a character value expression.
func Fold(
	subjectAny interface{},
	foldCase grammar.FoldCase,
) *FoldFunction {
	var ref types.Relation
	switch subjectAny := subjectAny.(type) {
	case types.Projection:
		ref = subjectAny.References()
	}
	subject := inspect.CharacterValueExpressionFromAny(subjectAny)
	if subject == nil {
		msg := fmt.Sprintf(
			"expected coerceable CharacterValueExpression but got %+v(%T)",
			subjectAny, subjectAny,
		)
		panic(msg)
	}
	return &FoldFunction{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		FoldFunction: &grammar.FoldFunction{
			Case:    foldCase,
			Subject: *subject,
		},
	}
}

// FoldFunction wraps the UPPER() or LOWER() SQL function grammar element
type FoldFunction struct {
	BaseFunction
	*grammar.FoldFunction
}

// CommonValueExpression returns the object as a
// `*grammar.CommonValueExpression`
func (f *FoldFunction) CommonValueExpression() *grammar.CommonValueExpression {
	return &grammar.CommonValueExpression{
		String: &grammar.StringValueExpression{
			Character: &grammar.CharacterValueExpression{
				Factor: &grammar.CharacterFactor{
					Primary: grammar.CharacterPrimary{
						Function: &grammar.StringValueFunction{
							Character: &grammar.CharacterValueFunction{
								Fold: f.FoldFunction,
							},
						},
					},
				},
			},
		},
	}
}

// DerivedColumn returns the `*grammar.DerivedColumn` element representing
// the Projection
func (f *FoldFunction) DerivedColumn() *grammar.DerivedColumn {
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
func (f *FoldFunction) As(alias string) types.Projection {
	f.alias = alias
	return f
}

// Convert returns a TranscodingFunction that produces a CONVERT() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The first argument is the subject of the CONVERT function and must be
// coercible to a character value expression. The second argument is the USING
// portion of the CONVERT function.
func Convert(
	subjectAny interface{},
	using string,
) *TranscodingFunction {
	var ref types.Relation
	switch subjectAny := subjectAny.(type) {
	case types.Projection:
		ref = subjectAny.References()
	}
	subject := inspect.CharacterValueExpressionFromAny(subjectAny)
	if subject == nil {
		msg := fmt.Sprintf(
			"expected coerceable CharacterValueExpression but got %+v(%T)",
			subjectAny, subjectAny,
		)
		panic(msg)
	}
	return &TranscodingFunction{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		TranscodingFunction: &grammar.TranscodingFunction{
			Subject: *subject,
			Using: grammar.SchemaQualifiedName{
				Identifiers: grammar.IdentifierChain{
					Identifiers: strings.Split(using, "."),
				},
			},
		},
	}
}

// TranscodingFunction wraps the CONVERT() SQL function grammar element
type TranscodingFunction struct {
	BaseFunction
	*grammar.TranscodingFunction
}

// CommonValueExpression returns the object as a
// `*grammar.CommonValueExpression`
func (f *TranscodingFunction) CommonValueExpression() *grammar.CommonValueExpression {
	return &grammar.CommonValueExpression{
		String: &grammar.StringValueExpression{
			Character: &grammar.CharacterValueExpression{
				Factor: &grammar.CharacterFactor{
					Primary: grammar.CharacterPrimary{
						Function: &grammar.StringValueFunction{
							Character: &grammar.CharacterValueFunction{
								Transcoding: f.TranscodingFunction,
							},
						},
					},
				},
			},
		},
	}
}

// DerivedColumn returns the `*grammar.DerivedColumn` element representing
// the Projection
func (f *TranscodingFunction) DerivedColumn() *grammar.DerivedColumn {
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
func (f *TranscodingFunction) As(alias string) types.Projection {
	f.alias = alias
	return f
}

// Translate returns a TransliterationFunction that produces a TRANSLATE() SQL
// function that can be passed to sqlb constructs and functions like Select()
//
// The first argument is the subject of the TRANSLATE function and must be
// coercible to a character value expression. The second argument is the USING
// portion of the TRANSLATE function.
func Translate(
	subjectAny interface{},
	using string,
) *TransliterationFunction {
	var ref types.Relation
	switch subjectAny := subjectAny.(type) {
	case types.Projection:
		ref = subjectAny.References()
	}
	subject := inspect.CharacterValueExpressionFromAny(subjectAny)
	if subject == nil {
		msg := fmt.Sprintf(
			"expected coerceable CharacterValueExpression but got %+v(%T)",
			subjectAny, subjectAny,
		)
		panic(msg)
	}
	return &TransliterationFunction{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		TransliterationFunction: &grammar.TransliterationFunction{
			Subject: *subject,
			Using: grammar.SchemaQualifiedName{
				Identifiers: grammar.IdentifierChain{
					Identifiers: strings.Split(using, "."),
				},
			},
		},
	}
}

// TransliterationFunction wraps the TRANSLATE() SQL function grammar element
type TransliterationFunction struct {
	BaseFunction
	*grammar.TransliterationFunction
}

// CommonValueExpression returns the object as a
// `*grammar.CommonValueExpression`
func (f *TransliterationFunction) CommonValueExpression() *grammar.CommonValueExpression {
	return &grammar.CommonValueExpression{
		String: &grammar.StringValueExpression{
			Character: &grammar.CharacterValueExpression{
				Factor: &grammar.CharacterFactor{
					Primary: grammar.CharacterPrimary{
						Function: &grammar.StringValueFunction{
							Character: &grammar.CharacterValueFunction{
								Transliteration: f.TransliterationFunction,
							},
						},
					},
				},
			},
		},
	}
}

// DerivedColumn returns the `*grammar.DerivedColumn` element representing
// the Projection
func (f *TransliterationFunction) DerivedColumn() *grammar.DerivedColumn {
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
func (f *TransliterationFunction) As(alias string) types.Projection {
	f.alias = alias
	return f
}

// Trim returns a TrimFunction that produces a TRIM([LEADING|TRAILING] chars
// FROM col) SQL function that can be passed to sqlb constructs and functions
// like Select()
//
// The first argument is the subject of the TRIM function and must be coercible
// to a character value expression. The second argument is the character(s) you
// wish to trim from the subject. The second argument must be coercible to a
// character value expression. The third argument specifies whether the
// leading, trailing or both sides of the subject should be trimmed.
func Trim(
	subjectAny interface{},
	charsAny interface{},
	spec grammar.TrimSpecification,
) *TrimFunction {
	var ref types.Relation
	switch subjectAny := subjectAny.(type) {
	case types.Projection:
		ref = subjectAny.References()
	}
	subject := inspect.CharacterValueExpressionFromAny(subjectAny)
	if subject == nil {
		msg := fmt.Sprintf(
			"expected coerceable CharacterValueExpression for "+
				"subject argument but got %+v(%T)",
			subjectAny, subjectAny,
		)
		panic(msg)
	}
	chars := inspect.CharacterValueExpressionFromAny(charsAny)
	if chars == nil {
		msg := fmt.Sprintf(
			"expected coerceable CharacterValueExpression for "+
				"chars argument but got %+v(%T)",
			charsAny, charsAny,
		)
		panic(msg)
	}
	return &TrimFunction{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		TrimFunction: &grammar.TrimFunction{
			Subject:       *subject,
			Character:     chars,
			Specification: spec,
		},
	}
}

// TrimSpace returns a TrimFunction that produces a TRIM(col) SQL
// function that can be passed to sqlb constructs and functions like Select()
func TrimSpace(
	subjectAny interface{},
) *TrimFunction {
	return doTrimSpace(subjectAny, grammar.TrimSpecificationBoth)
}

// LTrimSpace returns a TrimFunction that produces a TRIM(LEADING col) SQL
// function that can be passed to sqlb constructs and functions like Select()
func LTrimSpace(
	subjectAny interface{},
) *TrimFunction {
	return doTrimSpace(subjectAny, grammar.TrimSpecificationLeading)
}

// RTrimSpace returns a TrimFunction that produces a TRIM(TRAILING col) SQL
// function that can be passed to sqlb constructs and functions like Select()
func RTrimSpace(
	subjectAny interface{},
) *TrimFunction {
	return doTrimSpace(subjectAny, grammar.TrimSpecificationTrailing)
}

func doTrimSpace(
	subjectAny interface{},
	spec grammar.TrimSpecification,
) *TrimFunction {
	var ref types.Relation
	switch subjectAny := subjectAny.(type) {
	case types.Projection:
		ref = subjectAny.References()
	}
	subject := inspect.CharacterValueExpressionFromAny(subjectAny)
	if subject == nil {
		msg := fmt.Sprintf(
			"expected coerceable CharacterValueExpression for "+
				"subject argument but got %+v(%T)",
			subjectAny, subjectAny,
		)
		panic(msg)
	}
	return &TrimFunction{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		TrimFunction: &grammar.TrimFunction{
			Specification: spec,
			Subject:       *subject,
		},
	}
}

// LTrim returns a TrimFunction that produces a TRIM(LEADING char FROM col) SQL
// function that can be passed to sqlb constructs and functions like Select()
func LTrim(
	subjectAny interface{},
	charsAny interface{},
) *TrimFunction {
	return Trim(subjectAny, charsAny, grammar.TrimSpecificationLeading)
}

var TrimPrefix = LTrim

// RTrim returns a TrimFunction that produces a TRIM(TRAILING char FROM col) SQL
// function that can be passed to sqlb constructs and functions like Select()
func RTrim(
	subjectAny interface{},
	charsAny interface{},
) *TrimFunction {
	return Trim(subjectAny, charsAny, grammar.TrimSpecificationTrailing)
}

var TrimSuffix = RTrim

// TrimFunction wraps the TRIM() SQL function with a regular
// expression matching variant grammar element
type TrimFunction struct {
	BaseFunction
	*grammar.TrimFunction
}

// CommonValueExpression returns the object as a
// `*grammar.CommonValueExpression`
func (f *TrimFunction) CommonValueExpression() *grammar.CommonValueExpression {
	return &grammar.CommonValueExpression{
		String: &grammar.StringValueExpression{
			Character: &grammar.CharacterValueExpression{
				Factor: &grammar.CharacterFactor{
					Primary: grammar.CharacterPrimary{
						Function: &grammar.StringValueFunction{
							Character: &grammar.CharacterValueFunction{
								Trim: f.TrimFunction,
							},
						},
					},
				},
			},
		},
	}
}

// DerivedColumn returns the `*grammar.DerivedColumn` element representing
// the Projection
func (f *TrimFunction) DerivedColumn() *grammar.DerivedColumn {
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
func (f *TrimFunction) As(alias string) types.Projection {
	f.alias = alias
	return f
}

/*
// Ascii returns a Projection that contains the ASCII() SQL function
func Ascii(p api.Projection) api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_ASCII),
		elements: []api.Element{p.(api.Element)},
		sel:      p.From(),
	}
}

// Reverse returns a Projection that contains the REVERSE() SQL function
func Reverse(p api.Projection) api.Projection {
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_REVERSE),
		elements: []api.Element{p.(api.Element)},
		sel:      p.From(),
	}
}

// Concat returns a Projection that contains the CONCAT() SQL function
func Concat(projs ...api.Projection) api.Projection {
	els := make([]api.Element, len(projs))
	for x, p := range projs {
		els[x] = p.(api.Element)
	}
	subjects := element.NewList(els...)
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CONCAT),
		elements: []api.Element{subjects},
		// TODO(jaypipes): Clearly we need to support >1 selection...
		sel: projs[0].From(),
	}
}

// ConcatWs returns a Projection that contains the CONCAT_WS() SQL function
func ConcatWs(sep string, projs ...api.Projection) api.Projection {
	els := make([]api.Element, len(projs))
	for x, p := range projs {
		els[x] = p.(api.Element)
	}
	subjects := element.NewList(els...)
	return &Function{
		ScanInfo: grammar.FunctionScanTable(grammar.FUNC_CONCAT_WS),
		elements: []api.Element{element.NewValue(nil, sep), subjects},
		// TODO(jaypipes): Clearly we need to support >1 selection...
		sel: projs[0].From(),
	}
}
*/

/*
// TRIM/BTRIM/LTRIM/RTRIM SQL function support
//
// For MySQL, the TRIM() SQL function takes the following forms.
//
// Remove whitespace from before and after string:
//		TRIM(string)
// Remove whitespace from before string:
//		LTRIM(string)
// Remove whitespace from after string:
//		RTRIM(string)
// Remove the string remstr from before and after string:
//		TRIM(remstr FROM string)
// Remove the string remstr from before string:
//		TRIM(LEADING remstr FROM string)
// Remove the string remstr from after string:
//		TRIM(TRAILING remstr FROM string)
//
// For PostgreSQL, the TRIM() SQL function takes the following forms.
//
// Remove whitespace from before and after string:
//		BTRIM(string)
// Remove whitespace from before string:
//		TRIM(LEADING FROM string)
// Remove whitespace from after string:
//		TRIM(TRAILING FROM string)
// Remove longest string containing any character in chars from before
// and after string:
//		BTRIM(string, chars)
// Remove longest string containing any character in chars from before
// string:
//		TRIM(LEADING chars FROM string)
// Remove longest string containing any character in chars from after
// string:
//		TRIM(TRAILING chars FROM string)

type TrimLocation int

const (
	TRIM_BOTH TrimLocation = iota
	TRIM_TRAILING
	TRIM_LEADING
)

type TrimFunction struct {
	sel      api.Selection
	alias    string
	subject  api.Element
	chars    string
	location TrimLocation
}

func (f *TrimFunction) From() api.Selection {
	return f.sel
}

func (f *TrimFunction) DisableAliasScan() func() {
	origAlias := f.alias
	f.alias = ""
	return func() { f.alias = origAlias }
}

func (f *TrimFunction) As(alias string) api.Projection {
	aliased := &TrimFunction{
		sel:      f.sel,
		alias:    alias,
		subject:  f.subject,
		location: f.location,
		chars:    f.chars,
	}
	return aliased
}

func (f *TrimFunction) ArgCount() int {
	argc := f.subject.ArgCount()
	if f.chars != "" {
		argc++
	}
	return argc
}

func (f *TrimFunction) Desc() api.Orderable {
	return sortcolumn.NewDesc(f)
}

func (f *TrimFunction) Asc() api.Orderable {
	return sortcolumn.NewAsc(f)
}

// Helper function that scans into the supplied SQL []byte buffer for the
// TRIM/BTRIM() SQL function for MySQL
// TRIM() function for MySQL variants
func TrimFunctionScanMySQL(
	f *TrimFunction,
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	switch f.location {
	case TRIM_LEADING:
		if f.chars == "" {
			b.Write(grammar.Symbols[grammar.SYM_LTRIM])
		} else {
			b.Write(grammar.Symbols[grammar.SYM_TRIM])
			b.Write(grammar.Symbols[grammar.SYM_LEADING])
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			b.WriteString(builder.InterpolationMarker(opts, *curarg))
			qargs[*curarg] = f.chars
			*curarg++
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			b.Write(grammar.Symbols[grammar.SYM_FROM])
		}
		b.WriteString(TrimFunctionScanSubject(f, opts, qargs, curarg))
	case TRIM_TRAILING:
		if f.chars == "" {
			b.Write(grammar.Symbols[grammar.SYM_RTRIM])
		} else {
			b.Write(grammar.Symbols[grammar.SYM_TRIM])
			b.Write(grammar.Symbols[grammar.SYM_TRAILING])
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			b.WriteString(builder.InterpolationMarker(opts, *curarg))
			qargs[*curarg] = f.chars
			*curarg++
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			b.Write(grammar.Symbols[grammar.SYM_FROM])
		}
		b.WriteString(TrimFunctionScanSubject(f, opts, qargs, curarg))
	case TRIM_BOTH:
		if f.chars == "" {
			b.Write(grammar.Symbols[grammar.SYM_TRIM])
		} else {
			b.Write(grammar.Symbols[grammar.SYM_TRIM])
			b.WriteString(builder.InterpolationMarker(opts, *curarg))
			qargs[*curarg] = f.chars
			*curarg++
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			b.Write(grammar.Symbols[grammar.SYM_FROM])
		}
		b.WriteString(TrimFunctionScanSubject(f, opts, qargs, curarg))
	}
	return b.String()
}

// Helper function that scans into the supplied SQL []byte buffer for the
// TRIM/BTRIM() SQL function for PostgreSQL
// TRIM() function for PostgreSQL variants
func TrimFunctionScanPostgreSQL(
	f *TrimFunction,
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	switch f.location {
	case TRIM_LEADING:
		b.Write(grammar.Symbols[grammar.SYM_TRIM])
		b.Write(grammar.Symbols[grammar.SYM_LEADING])
		if f.chars != "" {
			b.WriteRune(' ')
			b.WriteString(builder.InterpolationMarker(opts, *curarg))
			qargs[*curarg] = f.chars
			*curarg++
		}
		b.Write(grammar.Symbols[grammar.SYM_SPACE])
		b.Write(grammar.Symbols[grammar.SYM_FROM])
		b.WriteString(TrimFunctionScanSubject(f, opts, qargs, curarg))
	case TRIM_TRAILING:
		b.Write(grammar.Symbols[grammar.SYM_TRIM])
		b.Write(grammar.Symbols[grammar.SYM_TRAILING])
		if f.chars != "" {
			b.Write(grammar.Symbols[grammar.SYM_SPACE])
			b.WriteString(builder.InterpolationMarker(opts, *curarg))
			qargs[*curarg] = f.chars
			*curarg++
		}
		b.Write(grammar.Symbols[grammar.SYM_SPACE])
		b.Write(grammar.Symbols[grammar.SYM_FROM])
		b.WriteString(TrimFunctionScanSubject(f, opts, qargs, curarg))
	case TRIM_BOTH:
		b.Write(grammar.Symbols[grammar.SYM_BTRIM])
		b.WriteString(TrimFunctionScanSubject(f, opts, qargs, curarg))
		if f.chars != "" {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
			b.WriteString(builder.InterpolationMarker(opts, *curarg))
			qargs[*curarg] = f.chars
			*curarg++
		}
	}
	return b.String()
}

// Scan in the subject of the TRIM() function
func TrimFunctionScanSubject(
	f *TrimFunction,
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	// We need to disable alias output for elements that are
	// projections. We don't want to output, for example,
	// "ON users.id AS user_id = TRIM(articles.author)"
	switch f.subject.(type) {
	case api.Projection:
		reset := f.subject.(api.Projection).DisableAliasScan()
		defer reset()
	}
	b.WriteString(f.subject.String(opts, qargs, curarg))
	return b.String()
}

func (f *TrimFunction) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	switch opts.Dialect() {
	case api.DialectPostgreSQL:
		b.WriteString(TrimFunctionScanPostgreSQL(f, opts, qargs, curarg))
	default:
		b.WriteString(TrimFunctionScanMySQL(f, opts, qargs, curarg))
	}
	b.Write(grammar.Symbols[grammar.SYM_RPAREN])
	if f.alias != "" {
		b.Write(grammar.Symbols[grammar.SYM_AS])
		b.WriteString(f.alias)
	}
	return b.String()
}

// Returns a struct that will output the TRIM() SQL function, trimming leading
// and trailing whitespace from the supplied projection
func Trim(p api.Projection) api.Projection {
	return &TrimFunction{
		subject:  p.(api.Element),
		sel:      p.From(),
		location: TRIM_BOTH,
	}
}

// Returns a struct that will output the LTRIM() SQL function for MySQL and the
// TRIM(LEADING FROM column) SQL function for PostgreSQL. The SQL function in
// either case will remove whitespace from the start of the supplied projection
func LTrim(p api.Projection) api.Projection {
	return &TrimFunction{
		subject:  p.(api.Element),
		sel:      p.From(),
		location: TRIM_LEADING,
	}
}

// Returns a struct that will output the RTRIM() SQL function for MySQL and the
// TRIM(TRAILING FROM column) SQL function for PostgreSQL. The SQL function in
// either case will remove whitespace from the start of the supplied projection
func RTrim(p api.Projection) api.Projection {
	return &TrimFunction{
		subject:  p.(api.Element),
		sel:      p.From(),
		location: TRIM_TRAILING,
	}
}

// Returns a struct that will output the TRIM() SQL function, trimming leading
// and trailing specified characters from the supplied projection
func TrimChars(p api.Projection, chars string) api.Projection {
	return &TrimFunction{
		subject:  p.(api.Element),
		sel:      p.From(),
		location: TRIM_BOTH,
		chars:    chars,
	}
}

// Returns a struct that will output the TRIM(LEADING chars FROM column) SQL
// function, trimming leading specified characters from the supplied projection
func LTrimChars(p api.Projection, chars string) api.Projection {
	return &TrimFunction{
		subject:  p.(api.Element),
		sel:      p.From(),
		location: TRIM_LEADING,
		chars:    chars,
	}
}

// Returns a struct that will output the TRIM(TRAILING chars FROM column) SQL
// function, trimming trailing specified characters from the supplied
// projection
func RTrimChars(p api.Projection, chars string) api.Projection {
	return &TrimFunction{
		subject:  p.(api.Element),
		sel:      p.From(),
		location: TRIM_TRAILING,
		chars:    chars,
	}
}
*/
