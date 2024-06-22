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

type ElementSizes struct {
	// The number of interface{} arguments that the element will add to the
	// slice of interface{} arguments that will eventually be constructed for
	// the Query
	ArgCount int
	// The number of bytes in the output buffer to represent this element
	BufferSize int
}

// Builder holds information about the formatting and dialect of the output SQL
// that sqlb writes to the output buffer
type Builder struct {
	*strings.Builder
	Dialect api.Dialect
	Format  api.FormatOptions
}

// Scan takes two slices and a pointer to an int. The first slice is a slice of
// bytes that the implementation should copy its string representation to and
// the other slice is a slice of interface{} values that the element should add
// its arguments to. The pointer to an int is the index of the current argument
// to be processed. The method returns a single int, the number of bytes
// written to the buffer.
func (b *Builder) Scan(args []interface{}, scannables ...Scannable) {
	curArg := 0
	b.WriteString(b.Format.PrefixWith)
	for _, scannable := range scannables {
		scannable.Scan(b, args, &curArg)
	}
}

func (b *Builder) Size(elements ...Element) *ElementSizes {
	buflen := 0
	argc := 0

	for _, el := range elements {
		argc += el.ArgCount()
		buflen += el.Size(b)
	}
	buflen += b.InterpolationLength(argc)
	buflen += len(b.Format.PrefixWith)

	return &ElementSizes{
		ArgCount:   argc,
		BufferSize: buflen,
	}
}

// StringArgs returns the built query string and a slice of interface{}
// representing the values of the query args used in the query string, if any.
func (b *Builder) StringArgs(el Element) (string, []interface{}) {
	sizes := b.Size(el)
	qargs := make([]interface{}, sizes.ArgCount)
	b.Grow(sizes.BufferSize)
	b.Scan(qargs, el)
	return b.String(), qargs
}

// AddInterpolationMarker adds an interpolation marker of the specified
// dialect and position into the built SQL string
func (b *Builder) AddInterpolationMarker(position int) {
	if b.Dialect == api.DialectPostgreSQL {
		b.Write(grammar.Symbols[grammar.SYM_DOLLAR])
		b.WriteString(strconv.Itoa(position + 1))
	} else {
		b.Write(grammar.Symbols[grammar.SYM_QUEST_MARK])
	}
}

// InterpolationLength returns the total length of the characters representing
// interpolation markers for query parameters. Different SQL dialects use
// different character sequences for marking query parameters during query
// preparation. For instance, MySQL and SQLite use the ? character. PostgreSQL
// uses a numbered $N schema with N starting at 1, SQL Server uses a :N scheme,
// etc.
func (b *Builder) InterpolationLength(argc int) int {
	if b.Dialect == api.DialectPostgreSQL {
		// $ character for each interpolated parameter plus ones digit of
		// number
		size := 2 * argc
		if argc > 9 {
			// tens digit
			size += argc - 9
		}
		if argc > 99 {
			// hundreds digit
			size += argc - 99
		}
		if argc > 999 {
			// thousands digit
			size += argc - 999
		}
		return size
	}
	return argc // Single question mark used as interpolation marker
}

// New returns a builder for the supplied dialect
func New(
	mods ...api.OptionModifier,
) *Builder {
	opts := api.MergeOptions(mods)
	if opts.Format == nil {
		opts.Format = &defaultFormatOptions
	}
	if opts.Dialect == nil {
		opts.Dialect = &defaultDialect
	}
	return &Builder{
		Builder: &strings.Builder{},
		Dialect: *opts.Dialect,
		Format:  *opts.Format,
	}
}
