//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package testutil

import (
	"database/sql"
	"fmt"

	"github.com/jaypipes/sqlb/api"
)

// Meta returns the Meta we use in testing
func Meta() *api.Meta {
	m := &api.Meta{
		Dialect: api.DialectMySQL,
		Name:    "test",
	}
	users := m.AddTable("users")
	users.AddColumn("id")
	users.AddColumn("name")

	articles := m.AddTable("articles")
	articles.AddColumn("id")
	articles.AddColumn("author")
	articles.AddColumn("state")

	articleStates := m.AddTable("article_states")
	articleStates.AddColumn("id")
	articleStates.AddColumn("name")

	userProfiles := m.AddTable("user_profiles")
	userProfiles.AddColumn("id")
	userProfiles.AddColumn("content")
	userProfiles.AddColumn("user")

	orgs := m.AddTable("organizations")
	orgs.AddColumn("id")
	orgs.AddColumn("uuid")
	orgs.AddColumn("root_organization_id")
	orgs.AddColumn("nested_set_left")
	orgs.AddColumn("nested_set_right")

	orgUsers := m.AddTable("organization_users")
	orgUsers.AddColumn("organization_id")
	orgUsers.AddColumn("user_id")

	return m
}

// ResetDB resets the testing database by dropping the database tables and
// recreating them.
func ResetDB(dialect api.Dialect, db *sql.DB) {
	var stmts []string
	switch dialect {
	case api.DialectMySQL:
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
	case api.DialectPostgreSQL:
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
