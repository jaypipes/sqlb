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
        tdefs: make(map[string]*TableDef, 0),
    }

    users := &TableDef{
        meta: meta,
        name: "users",
    }
    colUserId := &ColumnDef{
        name: "id",
        tdef: users,
    }
    colUserName := &ColumnDef{
        name: "name",
        tdef: users,
    }

    articles := &TableDef{
        meta: meta,
        name: "articles",
    }
    colArticleId := &ColumnDef{
        name: "id",
        tdef: articles,
    }
    colArticleAuthor := &ColumnDef{
        name: "author",
        tdef: articles,
    }
    colArticleState := &ColumnDef{
        name: "state",
        tdef: articles,
    }

    article_states := &TableDef{
        meta: meta,
        name: "article_states",
    }
    colArticleStateId := &ColumnDef{
        name: "id",
        tdef: article_states,
    }
    colArticleStateName := &ColumnDef{
        name: "name",
        tdef: article_states,
    }

    users.cdefs = []*ColumnDef{colUserId, colUserName}
    articles.cdefs = []*ColumnDef{colArticleId, colArticleAuthor, colArticleState}
    article_states.cdefs = []*ColumnDef{colArticleStateId, colArticleStateName}
    meta.tdefs["users"] = users
    meta.tdefs["articles"] = articles
    meta.tdefs["article_states"] = article_states
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

    assert.Equal(2, len(meta.tdefs))

    artTbl := meta.tdefs["articles"]
    userTbl := meta.tdefs["users"]

    assert.Equal("articles", artTbl.name)
    assert.Equal("users", userTbl.name)

    assert.Equal(7, len(userTbl.cdefs))
    assert.Equal(5, len(artTbl.cdefs))

    createdOnCol := userTbl.Column("created_on")
    assert.NotNil(createdOnCol)
    assert.Equal("created_on", createdOnCol.cdef.name)
}

func TestReflectErrors(t *testing.T) {
    assert := assert.New(t)

    err := Reflect("mysql", nil, nil)
    assert.NotNil(err)
    assert.Equal(ERR_NO_META_STRUCT, err)
}
