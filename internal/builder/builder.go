//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"strconv"
	"strings"

	"github.com/jaypipes/sqlb/core/types"
	"github.com/jaypipes/sqlb/grammar"
)

// Builder holds information about the formatting and dialect of the output SQL
// that sqlb writes to the output buffer
type Builder struct {
	strings.Builder
	opts types.Options
}

// StringArgs returns the built query string and a slice of interface{}
// representing the values of the query args used in the query string, if any.
func (b *Builder) StringArgs(target interface{}) (string, []interface{}) {
	b.WriteString(b.opts.FormatPrefixWith())
	switch el := target.(type) {
	case *grammar.UpdateStatementSearched:
		argc := 0
		ArgCount(el, &argc)
		qargs := make([]interface{}, argc)
		curarg := 0
		b.doUpdateStatementSearched(el, qargs, &curarg)
		return b.Builder.String(), qargs
	case *grammar.DeleteStatementSearched:
		argc := 0
		ArgCount(el, &argc)
		qargs := make([]interface{}, argc)
		curarg := 0
		b.doDeleteStatementSearched(el, qargs, &curarg)
		return b.Builder.String(), qargs
	case *grammar.InsertStatement:
		argc := len(el.Values)
		qargs := make([]interface{}, argc)
		curarg := 0
		b.doInsertStatement(el, qargs, &curarg)
		return b.Builder.String(), qargs
	case *grammar.QuerySpecification:
		argc := 0
		ArgCount(el, &argc)
		qargs := make([]interface{}, argc)
		curarg := 0
		b.doQuerySpecification(el, qargs, &curarg)
		return b.Builder.String(), qargs
	case *grammar.CursorSpecification:
		argc := 0
		ArgCount(el, &argc)
		qargs := make([]interface{}, argc)
		curarg := 0
		b.doCursorSpecification(el, qargs, &curarg)
		return b.Builder.String(), qargs
	default:
		return "", []interface{}{}
	}
}

// InterpolationMarker returns a string with an interpolation marker of the
// specified dialect and position
func InterpolationMarker(opts types.Options, position int) string {
	b := &strings.Builder{}
	if opts.Dialect() == types.DialectPostgreSQL {
		b.Write(grammar.Symbols[grammar.SYM_DOLLAR])
		b.WriteString(strconv.Itoa(position + 1))
	} else {
		b.Write(grammar.Symbols[grammar.SYM_QUEST_MARK])
	}
	return b.String()
}

// New returns a builder for the supplied dialect
func New(
	mods ...types.Option,
) *Builder {
	opts := types.MergeOptions(mods)
	return &Builder{
		opts: opts,
	}
}
