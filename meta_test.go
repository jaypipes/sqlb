package sqlb

import (
    "database/sql"
    "fmt"
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
)

var (
    _MYSQL_DB_INIT = []string{
        "DROP TABLE IF EXISTS articles",
        "DROP TABLE IF EXISTS users",
        `
        CREATE TABLE users (
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
        `
        CREATE TABLE articles (
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
)

func testFixtureMeta() *Meta {
    meta := &Meta{
        schemaName: "test",
        tables: make(map[string]*Table, 0),
    }

    users := &Table{
        meta: meta,
        name: "users",
    }
    colUserId := &Column{
        name: "id",
        tbl: users,
    }
    colUserName := &Column{
        name: "name",
        tbl: users,
    }

    articles := &Table{
        meta: meta,
        name: "articles",
    }
    colArticleId := &Column{
        name: "id",
        tbl: articles,
    }
    colArticleAuthor := &Column{
        name: "author",
        tbl: articles,
    }
    colArticleState := &Column{
        name: "state",
        tbl: articles,
    }

    articleStates := &Table{
        meta: meta,
        name: "article_states",
    }
    colArticleStateId := &Column{
        name: "id",
        tbl: articleStates,
    }
    colArticleStateName := &Column{
        name: "name",
        tbl: articleStates,
    }

    userProfiles := &Table{
        meta: meta,
        name: "user_profiles",
    }
    colUserProfileId := &Column{
        name: "id",
        tbl: userProfiles,
    }
    colUserProfileContent := &Column{
        name: "content",
        tbl: userProfiles,
    }
    colUserProfileUser := &Column{
        name: "user",
        tbl: userProfiles,
    }

    users.columns = []*Column{colUserId, colUserName}
    articles.columns = []*Column{colArticleId, colArticleAuthor, colArticleState}
    articleStates.columns = []*Column{colArticleStateId, colArticleStateName}
    userProfiles.columns = []*Column{colUserProfileId, colUserProfileUser, colUserProfileContent}
    meta.tables["users"] = users
    meta.tables["articles"] = articles
    meta.tables["article_states"] = articleStates
    meta.tables["user_profiles"] = userProfiles
    return meta
}

func resetDB(driver string, db *sql.DB) {
    var stmts []string
    switch driver {
    case "mysql":
        stmts = _MYSQL_DB_INIT
    }
    for _, stmt := range stmts {
        _, err := db.Exec(stmt)
        if err != nil {
            fmt.Printf("ERROR: Failed resetting database: %v", err)
        }
    }
}

func TestReflectMySQL(t *testing.T) {
    dsn, found := os.LookupEnv("SQLB_TESTING_MYSQL_DSN")
    if ! found {
        t.Skip("No SQLB_TESTING_MYSQL_DSN environ set")
    }
    assert := assert.New(t)

    db, err := sql.Open("mysql", dsn);
    assert.Nil(err)

    resetDB("mysql", db)

    var meta Meta
    err = Reflect("mysql", db, &meta)
    assert.Nil(err)

    assert.Equal(2, len(meta.tables))

    artTbl := meta.tables["articles"]
    userTbl := meta.tables["users"]

    assert.Equal("articles", artTbl.name)
    assert.Equal("users", userTbl.name)

    assert.Equal(7, len(userTbl.columns))
    assert.Equal(5, len(artTbl.columns))

    createdOnCol := userTbl.Column("created_on")
    assert.NotNil(createdOnCol)
    assert.Equal("created_on", createdOnCol.name)
}

func TestReflectErrors(t *testing.T) {
    assert := assert.New(t)

    err := Reflect("mysql", nil, nil)
    assert.NotNil(err)
    assert.Equal(ERR_NO_META_STRUCT, err)
}
