//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package inspect

import (
	"github.com/jaypipes/sqlb/core/grammar"
	"github.com/jaypipes/sqlb/core/types"
)

// TableReferenceFromAny evaluates the supplied interface argument and returns
// a *TableReference if the supplied argument can be converted into a
// TableReference, or nil if the conversion cannot be done.
func TableReferenceFromAny(
	subject interface{},
) *grammar.TableReference {
	switch v := subject.(type) {
	case *grammar.TableReference:
		return v
	case grammar.TableReference:
		return &v
	case types.TableReferenceConverter:
		return v.TableReference()
	case *grammar.TablePrimary:
		return &grammar.TableReference{Primary: v}
	case grammar.TablePrimary:
		return &grammar.TableReference{Primary: &v}
	case *grammar.JoinedTable:
		return &grammar.TableReference{Joined: v}
	case grammar.JoinedTable:
		return &grammar.TableReference{Joined: &v}
	}
	return nil
}

// TableReferenceByName returns a pointer to the TableReference that has a name
// or alias/correlation ID matching the supplied string.
func TableReferenceByName(
	refs []grammar.TableReference,
	search string,
) *grammar.TableReference {
	for _, ref := range refs {
		if ref.Primary != nil {
			p := ref.Primary
			if p.Correlation != nil && p.Correlation.Name == search {
				return &ref
			} else if p.TableName != nil && *p.TableName == search {
				return &ref
			} else if p.QueryName != nil && *p.QueryName == search {
				return &ref
			}
		} else if ref.Joined != nil {
			jt := ref.Joined
			if jt.Qualified != nil {
				found := TableReferenceByName(
					[]grammar.TableReference{
						jt.Qualified.Left,
						jt.Qualified.Right,
					},
					search,
				)
				if found != nil {
					return &ref
				}
			}
		}
	}
	return nil
}
