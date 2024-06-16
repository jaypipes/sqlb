//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sqlb

import (
	"context"
	"database/sql"

	"github.com/jaypipes/sqlb/internal/scanner"
)

// Query accepts a `database/sql` `DB` handle and a `pkg/scanner.Element` and
// calls the `databases/sql.DB.Query` method on the SQL string produced by the
// `Element`.
func Query(
	db *sql.DB,
	el scanner.Element,
) (*sql.Rows, error) {
	return QueryContext(context.TODO(), db, el)
}

// QueryContext accepts a `database/sql` `DB` handle and a `pkg/scanner.Element`
// and calls the `database/sql.DB.QueryContext` method on the SQL string
// produced by the `Element`.
func QueryContext(
	ctx context.Context,
	db *sql.DB,
	el scanner.Element,
) (*sql.Rows, error) {
	s := scanner.New()
	qs, qargs := s.StringArgs(el)
	return db.QueryContext(ctx, qs, qargs...)
}
