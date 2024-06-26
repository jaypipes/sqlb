# sqlb

[![Go Reference](https://pkg.go.dev/badge/github.com/jaypipes/sqlb.svg)](https://pkg.go.dev/github.com/jaypipes/sqlb)
[![Go Report Card](https://goreportcard.com/badge/github.com/jaypipes/sqlb)](https://goreportcard.com/report/github.com/jaypipes/sqlb)
[![Build Status](https://github.com/jaypipes/sqlb/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/jaypipes/sqlb/actions)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)

`sqlb` is a Go library designed for efficiently constructing SQL expressions in
a concise, readable fashion.

Instead of hand-constructing strings containing raw SQL, users of the `sqlb`
library instead construct query expressions and the `sqlb` library does the
work of producing the raw strings that get sent to a SQL database.

## Building SQL expressions, not strings

It's best to learn by example, so let's walk through a common way in which Go
applications might typically work with an underlying SQL database and transform
this application to instead work with the `sqlb` library, showing the resulting
gains in both code expressiveness, application speed and memory efficiency.

Our example will be a simple blogging application.

Imagine we have the following set of tables in our database:

```sql
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
);

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
```

Our blogging application's default home page might return information about the
last ten articles published. It's reasonable to believe that the following SQL
expression might be used to grab this information from the database:

```sql
SELECT
  articles.title,
  articles.content,
  articles.created_on
  users.name,
FROM articles
JOIN users
 ON articles.created_by = users.id
ORDER BY articles.created_on DESC
LIMIT 10
```

Our Go code for the server side of our application might look something like
this:

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
)

const (
    DSN = "root:password@/blogdb"
)

var db *sql.DB

type Article struct {
    Title string
    AuthorName string
    PublishedOn string
    Content string
}

func getArticles() []*Article {
    qs := `
SELECT
  articles.title,
  articles.content,
  articles.created_on
  users.name,
FROM articles
JOIN users
 ON articles.created_by = users.id
ORDER BY articles.created_on DESC
LIMIT 10
`
    res := make([]*Article, 0)
    rows, err := db.Query(qs)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        a := &Article{}
        err := rows.Scan(
            &a.Title,
            &a.Content,
            &a.PublishedOn,
            &a.AuthorName,
        )
        if err != nil {
            log.Fatal(err)
        }
        res = append(res, a)
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
    return res
}

func handler(w http.ResponseWriter, r *http.Request) {
    articleTemplate := `%s
-----------------------------------------------------
by %s on %s

%s
`
    articles := getArticles()
    for _, article := range articles {
        fmt.Fprintf(w, articleTemplate, article.Title, article.AuthorName,
                    article.PublishedOn, article.Content)
    }
}

func main() {
    if db, err := sql.Open("mysql", DSN); err != nil {
        log.Fatal(err)
    }
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

**Note**: Clearly, I'm not doing proper error handling and I'm hard-coding
things like the DSN that should be pulled from a configuration system in this
example code.

The above code works, but it's fragile in the face of inevitable change to the
application. What if we want to make the number of articles returned
configurable? What if we want to allow users to list only articles by a
particular author? In both of these cases, we will need to modify the
`getArticles()` function to modify the SQL query string that it constructs:

```go
func getArticles(numArticles int, byAuthor string) []*Articles {
    // Our collection of query arguments
    qargs := make([]interface{}, 0)
    qs := `
SELECT
  articles.title,
  articles.content,
  articles.created_on
  users.name,
FROM articles
JOIN users
 ON articles.created_by = users.id
`
    if byAuthor != "" {
        qs += "WHERE users.name = ? "
        qargs = append(qargs, byAuthor)
    }
    qs += `ORDER BY articles.created_on DESC
LIMIT ?`
    qargs = append(qargs, numArticles)
    res := make([]*Article, 0)
    rows, err := db.Query(qs, qargs...)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        a := &Article{}
        err := rows.Scan(
            &a.Title,
            &a.Content,
            &a.PublishedOn,
            &a.AuthorName,
        )
        if err != nil {
            log.Fatal(err)
        }
        res = append(res, a)
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
    return res
}
```

As you can see above, the minor enhancements to our application of allowing a
configurable number of articles and filtering by author have already begun to
make the `getArticles()` function unwieldy. The string being generated for our
SQL `SELECT` statement is difficult to read and hides the *intent* of the query
expression in the ugliness of string concatenation.

Adding more filtering capability brings more conditionals and more string
concatenation, leading to ever-increasing complexity and reduced code
readability.

`sqlb` is designed to solve this problem.

## Rewriting our application to use `sqlb`

Let's rewrite our example application above to use the `sqlb` library instead
of manually constructing SQL strings.

We start by initializing `sqlb`'s reflection system in our application's
`main()` entrypoint:

```go
import (
    "database/sql"

    "github.com/jaypipes/sqlb"
)


var (
    meta *sqlb.Meta
    articles *sqlb.Table
    users *sqlb.Table
)

func main() {
    if db, err := sql.Open("mysql", DSN); err != nil {
        log.Fatal(err)
    }
    if meta, err := sqlb.Reflect(db); err != nil {
        log.Fatal(err)
    }
    articles = meta.Table("articles")
    users = meta.Table("users")
}
```

The `sqlb.Meta` struct is now populated with information about the database,
including metadata about tables, columns, indexes, and relations. You use
`sqlb.Meta` when constructing `sqlb` Query Expressions.

We've set two package-scoped variables called `articles` and `users` that refer
to the "articles" and "users" database tables, respectively. We will refer to
these variables in our `getArticles()` function.

Let's transform our original `getArticles()` function -- before we added
support for a configurable number of articles and filtering by author -- to use
`sqlb`:

```go

func getArticles() []*Article {
    q := sqlb.Select(articles.C("title"), articles.C("content"),
                     articles.C("created_by"), users.C("name"))
    q.Join(users, sqlb.Equal(articles.C("author"), users.C("id")))
    q.OrderBy(articles.C("created_by").Desc())
    q.Limit(10)

    res := make([]*Article, 0)
    rows, err := sqlb.Query(db, q)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        a := &Article{}
        err := rows.Scan(
            &a.Title,
            &a.Content,
            &a.PublishedOn,
            &a.AuthorName,
        )
        if err != nil {
            log.Fatal(err)
        }
        res = append(res, a)
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
    return res
}
```

The above code ends up producing an identical SQL string as the original code
without any of the string concatenation.

Let's add in functionality to have a configurable number of returned articles
and optionally filter for a specific author's articles.

```go
func getArticles(numArticles int, byAuthor string) []*Articles {
    q := sqlb.Select(articles.C("title"), articles.C("content"),
                     articles.C("created_by"), users.C("name"))
    q.Join(users, sqlb.Equal(articles.C("author"), users.C("id")))
    if byAuthor != "" {
        q.Where(sqlb.Equal(users.C("name"), byAuthor))
    }
    q.OrderBy(articles.C("created_by").Desc())
    q.Limit(numArticles)

    res := make([]*Article, 0)
    rows, err := sqlb.Query(db, q)
    if err != nil {
        log.Fatal(err)
    }
    for rows.Next() {
        a := &Article{}
        err := rows.Scan(
            &a.Title,
            &a.Content,
            &a.PublishedOn,
            &a.AuthorName,
        )
        if err != nil {
            log.Fatal(err)
        }
        res = append(res, a)
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
    return res
}
```

No more manually constructing and reconstructing strings or tracking query
arguments. `sqlb` handles the SQL string construction for you as well as the
slice of query arguments, allowing you to write custom query code in a more
natural and efficient manner.

## License

`sqlb` is licensed under the Apache license version 2. See the
[COPYING](COPYING) file for more information.
