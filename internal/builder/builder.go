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
func (b *Builder) StringArgs(el api.Element) (string, []interface{}) {
	argc := el.ArgCount()
	qargs := make([]interface{}, argc)
	curarg := 0
	return b.String(el, qargs, &curarg), qargs
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
