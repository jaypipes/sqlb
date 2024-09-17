//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package testutil

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jaypipes/sqlb/core/meta"
	"github.com/jaypipes/sqlb/core/types"
)

var (
	m                 *meta.Meta
	users             *meta.Table
	articles          *meta.Table
	articleStates     *meta.Table
	userProfiles      *meta.Table
	organizations     *meta.Table
	organizationUsers *meta.Table
)

func init() {
	m = &meta.Meta{
		Dialect: types.DialectMySQL,
		Name:    "test",
	}

	users = meta.NewTable(
		m, "users",
		"id",
		"name",
		"created_on",
	)

	articles = meta.NewTable(
		m, "articles",
		"id",
		"author",
		"state",
	)

	articleStates = meta.NewTable(
		m, "article_states",
		"id",
		"name",
	)

	userProfiles = meta.NewTable(
		m, "user_profiles",
		"id",
		"content",
		"user",
	)

	organizations = meta.NewTable(
		m, "organizations",
		"id",
		"uuid",
		"root_organization_id",
		"nested_set_left",
		"nested_set_right",
	)

	organizationUsers = meta.NewTable(
		m, "organization_users",
		"organization_id",
		"user_id",
	)

	m.Tables = map[string]*meta.Table{
		"users":              users,
		"articles":           articles,
		"articleStates":      articleStates,
		"userProfiles":       userProfiles,
		"organizations":      organizations,
		"organization_users": organizationUsers,
	}
}

// M returns the Meta we use in testing
func M() *meta.Meta {
	return m
}

// T returns the Table from the testing Meta with the supplied name
func T(name string) *meta.Table {
	return m.T(name)
}

// C returns the Column from the testing Meta with the supplied table.column
// dotted notation
func C(name string) types.Projection {
	names := strings.Split(name, ".")
	tname := names[0]
	cname := names[1]
	return m.T(tname).C(cname)
}

// ResetDB resets the testing database by dropping the database tables and
// recreating them.
func ResetDB(dialect types.Dialect, db *sql.DB) {
	var stmts []string
	switch dialect {
	case types.DialectMySQL:
		stmts = []string{
			"DROP TABLE IF EXISTS articles",
			"DROP TABLE IF EXISTS users",
			`CREATE TABLE users (
          id INT NOT NULL,
          email VARCHAR(100) NOT NULL,
          name VARCHAR(100) NOT NULL,
          is_author CHAR(1) NOT NULL,
          profile TEXT NULL,
          created_on DATETIME NOT NULL,
          updated_on DATETIME NOT NULL,
          PRIMARY KEY (id),
          UNIQUE INDEX (email)
        )`,
			`CREATE TABLE articles (
          id INT NOT NULL,
          title VARCHAR(200) NOT NULL,
          content TEXT NOT NULL,
          created_by INT NOT NULL,
          published_on DATETIME NULL,
          PRIMARY KEY (id),
          INDEX ix_title (title),
          FOREIGN KEY fk_users (created_by) REFERENCES users (id)
        );
        `,
		}
	case types.DialectPostgreSQL:
		stmts = []string{
			"BEGIN",
			"DROP TABLE IF EXISTS articles",
			"DROP TABLE IF EXISTS users",
			`CREATE TABLE users (
          id SERIAL NOT NULL,
          email VARCHAR(100) NOT NULL UNIQUE,
          name VARCHAR(100) NOT NULL,
          is_author CHAR(1) NOT NULL,
          profile TEXT NULL,
          created_on TIMESTAMP NOT NULL,
          updated_on TIMESTAMP NOT NULL,
          PRIMARY KEY (id)
        )`,
			`CREATE TABLE articles (
          id SERIAL NOT NULL,
          title VARCHAR(200) NOT NULL,
          content TEXT NOT NULL,
          created_by INT NOT NULL,
          published_on TIMESTAMP NULL,
          PRIMARY KEY (id),
          CONSTRAINT fk_users FOREIGN KEY (created_by) REFERENCES users (id)
        );
        `,
			"CREATE INDEX ix_title ON articles (title);",
			"COMMIT",
		}
	}
	for _, stmt := range stmts {
		_, err := db.Exec(stmt)
		if err != nil {
			fmt.Printf("ERROR: Failed resetting database: %v", err)
		}
	}
}
