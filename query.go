//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sqlb

import (
	"context"
	"database/sql"

	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/types"
)

// Query accepts a `database/sql` `DB` handle and a `pkg/types.Element` and
// calls the `databases/sql.DB.Query` method on the SQL string produced by the
// `Element`.
func Query(
	db *sql.DB,
	el types.Element,
) (*sql.Rows, error) {
	return QueryContext(context.TODO(), db, el)
}

// QueryContext accepts a `databases/sql` `DB` handle and a
// `pkg/types.Element` and calls the `databases/sql.DB.Query` method on the
// SQL string produced by the `Element`.
func QueryContext(
	ctx context.Context,
	db *sql.DB,
	el types.Element,
) (*sql.Rows, error) {
	s := scanner.New(types.DIALECT_UNKNOWN)
	qs, qargs := s.StringArgs(el)
	return db.QueryContext(ctx, qs, qargs...)
}
