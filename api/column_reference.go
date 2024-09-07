//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import (
	"github.com/jaypipes/sqlb/grammar"
)

// ColumnReferenceFromAny evaluates the supplied interface argument and returns
// a *ColumnReference if the supplied argument can be converted into a
// ColumnReference, or nil if the conversion cannot be done.
func ColumnReferenceFromAny(
	subject interface{},
) *grammar.ColumnReference {
	switch v := subject.(type) {
	case *grammar.ColumnReference:
		return v
	case grammar.ColumnReference:
		return &v
	case *Column:
		return &grammar.ColumnReference{
			BasicIdentifierChain: &grammar.IdentifierChain{
				Identifiers: []string{v.TableName(), v.name},
			},
		}
	}
	return nil
}
