//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// LIMIT <count> [ <offset> ]

// LimitClause represents the SQL MySQL/PostgreSQL extension LIMIT clause
type LimitClause struct {
	Count  int
	Offset *int
}

func (c *LimitClause) ArgCount(count *int) {
	if c.Offset != nil {
		*count = *count + 2
	} else {
		*count++
	}
}
