//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"strconv"
	"strings"

	"github.com/jaypipes/sqlb/api"
	upg "github.com/jaypipes/sqlb/grammar"
	"github.com/jaypipes/sqlb/internal/grammar"
)

// Builder holds information about the formatting and dialect of the output SQL
// that sqlb writes to the output buffer
type Builder struct {
	sb   *strings.Builder
	opts api.Options
}

// StringArgs returns the built query string and a slice of interface{}
// representing the values of the query args used in the query string, if any.
func (b *Builder) StringArgs(target interface{}) (string, []interface{}) {
	switch el := target.(type) {
	case api.Element:
		argc := el.ArgCount()
		qargs := make([]interface{}, argc)
		curarg := 0
		return b.String(el, qargs, &curarg), qargs
	case *upg.InsertStatement:
		sb := &strings.Builder{}
		argc := len(el.Values)
		qargs := make([]interface{}, argc)
		curarg := 0
		sb.Write(grammar.Symbols[grammar.SYM_INSERT])
		// We don't add any table alias when outputting the table identifier
		sb.WriteString(el.TableName)
		sb.WriteRune(' ')
		sb.Write(grammar.Symbols[grammar.SYM_LPAREN])

		ncols := len(el.Columns)
		for x, c := range el.Columns {
			// We don't add the table identifier or use an alias when outputting
			// the column names in the <columns> element of the INSERT statement
			sb.WriteString(c)
			if x != (ncols - 1) {
				sb.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
			}
		}
		sb.Write(grammar.Symbols[grammar.SYM_VALUES])
		for x, v := range el.Values {
			sb.WriteString(InterpolationMarker(b.opts, curarg))
			qargs[curarg] = v
			curarg++
			if x != (ncols - 1) {
				sb.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
			}
		}
		sb.Write(grammar.Symbols[grammar.SYM_RPAREN])
		return sb.String(), qargs
	default:
		return "", []interface{}{}
	}
}

// InterpolationMarker returns a string with an interpolation marker of the
// specified dialect and position
func InterpolationMarker(opts api.Options, position int) string {
	b := &strings.Builder{}
	if opts.Dialect() == api.DialectPostgreSQL {
		b.Write(grammar.Symbols[grammar.SYM_DOLLAR])
		b.WriteString(strconv.Itoa(position + 1))
	} else {
		b.Write(grammar.Symbols[grammar.SYM_QUEST_MARK])
	}
	return b.String()
}

func (b *Builder) String(el api.Element, qargs []interface{}, curarg *int) string {
	b.sb.WriteString(b.opts.FormatPrefixWith())
	b.sb.WriteString(el.String(b.opts, qargs, curarg))
	return b.sb.String()
}

// New returns a builder for the supplied dialect
func New(
	mods ...api.Option,
) *Builder {
	opts := api.MergeOptions(mods)
	return &Builder{
		sb:   &strings.Builder{},
		opts: opts,
	}
}
