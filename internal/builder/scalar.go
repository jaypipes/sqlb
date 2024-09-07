//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

func (b *Builder) doScalar(
	el interface{},
	qargs []interface{},
	curarg *int,
) {
	qargs[*curarg] = el
	b.WriteString(InterpolationMarker(b.opts, *curarg))
	*curarg++
}
