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

	"github.com/jaypipes/sqlb/api"
)

var (
	m                 *api.Meta
	users             *api.Table
	articles          *api.Table
	articleStates     *api.Table
	userProfiles      *api.Table
	organizations     *api.Table
	organizationUsers *api.Table
)

func init() {
	m = &api.Meta{
		Dialect: api.DialectMySQL,
		Name:    "test",
	}

	users = &api.Table{
		Meta: m,
		Name: "users",
		Columns: map[string]*api.Column{
			"id": &api.Column{
				Table: users,
				Name:  "id",
			},
			"name": &api.Column{
				Table: users,
				Name:  "name",
			},
		},
	}

	articles = &api.Table{
		Name: "articles",
		Columns: map[string]*api.Column{
			"id": &api.Column{
				Table: articles,
				Name:  "id",
			},
			"author": &api.Column{
				Table: articles,
				Name:  "author",
			},
			"state": &api.Column{
				Table: articles,
				Name:  "state",
			},
		},
	}

	articleStates = &api.Table{
		Name: "article_states",
		Columns: map[string]*api.Column{
			"id": &api.Column{
				Table: articleStates,
				Name:  "id",
			},
			"name": &api.Column{
				Table: articleStates,
				Name:  "name",
			},
		},
	}

	userProfiles = &api.Table{
		Meta: m,
		Name: "user_profiles",
		Columns: map[string]*api.Column{
			"id": &api.Column{
				Table: userProfiles,
				Name:  "id",
			},
			"content": &api.Column{
				Table: userProfiles,
				Name:  "content",
			},
			"user": &api.Column{
				Table: userProfiles,
				Name:  "user",
			},
		},
	}

	organizations = &api.Table{
		Meta: m,
		Name: "organizations",
		Columns: map[string]*api.Column{
			"id": &api.Column{
				Table: organizations,
				Name:  "id",
			},
			"uuid": &api.Column{
				Table: organizations,
				Name:  "uuid",
			},
			"root_organization_id": &api.Column{
				Table: organizations,
				Name:  "root_organization_id",
			},
			"nested_set_left": &api.Column{
				Table: organizations,
				Name:  "nested_set_left",
			},
			"nested_set_right": &api.Column{
				Table: organizations,
				Name:  "nested_set_right",
			},
		},
	}

	organizationUsers = &api.Table{
		Meta: m,
		Name: "organization_users",
		Columns: map[string]*api.Column{
			"organization_id": &api.Column{
				Table: organizationUsers,
				Name:  "organization_id",
			},
			"user_id": &api.Column{
				Table: organizationUsers,
				Name:  "user_id",
			},
		},
	}

	m.Tables = map[string]*api.Table{
		"users":              users,
		"articles":           articles,
		"articleStates":      articleStates,
		"userProfiles":       userProfiles,
		"organizations":      organizations,
		"organization_users": organizationUsers,
	}
}

// M returns the Meta we use in testing
func M() *api.Meta {
	return m
}

// T returns the Table from the testing Meta with the supplied name
func T(name string) *api.Table {
	return m.T(name)
}

// C returns the Column from the testing Meta with the supplied table.column
// dotted notation
func C(name string) *api.Column {
	names := strings.Split(name, ".")
	tname := names[0]
	cname := names[1]
	return m.T(tname).C(cname)
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
