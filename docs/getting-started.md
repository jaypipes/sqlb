# Getting started with `sqlb`

OK, so you've decided to stop manually constructing SQL strings in your
application and instead want to give the `sqlb` library a try. The first thing
you'll want to do is grab the `sqlb` source code for use in your project.

## Grabbing the `sqlb` library code

If you're just doing some one-off work or trying things out, you can grab the
latest `sqlb` code using `go get`:

```
go get github.com/jaypipes/sqlb
```

For anything more serious, best practice is to use some sort of
dependency/vendor package management in your projects. One such dependency
management system is [`govendor`](https://github.com/kardianos/govendor/).

Use `govendor` to add `sqlb` as a dependency to your own Golang project by
calling `govendor fetch` from the root of your own source tree:

```
govendor fetch github.com/jaypipes/sqlb
```

**Tip**: We recommend specifying a release of `sqlb` using the syntax
"@$RELEASE" when adding a dependency, like so:

```
govendor fetch github.com/jaypipes/sqlb@0.5
```

This will pin your installation of `sqlb` to a specific version, ensuring that
when you do something like `govendor sync` you will *not* bring in any new code
for `sqlb`. When you need some functionality or bug fix in `sqlb` that is a
newer release, you will explicitly fetch that newer version of `sqlb`. This
makes vendor/dependency management an explicit and deterministic activity.

## Reflecting database schema metadata

Once you have pulled the `sqlb` library source code as a dependency for your
own Golang project, the first step in using `sqlb` will be to have `sqlb`
discover metadata about the tables and columns in your database.

Most Golang applications that access an RDBMS will have a single place where
the application initializes a `database/sql:DB` struct for use in querying the
application's database.

Typical code might look like this:

```go
package myproject

import (
    "database/sql"
    "log"
)

const (
    DSN = "user:password@/blog"
)

var (
    db *sql.DB
)

func main() {
    if db, err := sql.Open("mysql", DSN); err != nil {
        log.Fatal(err)
    }
    // Do something with the "db" variable...
}
```

**Note**: Clearly, the `DSN` variable wouldn't be hard-coded; the above is for
illustrative purposes.

Immediately after the call to `sql.Open()` is a good time to ask `sqlb` to
discover metadata about the schema's tables and those tables' columns. Use the
`sqlb.Reflect()` method to discover this metadata:

```go
package myproject

import (
    "database/sql"
    "log"

    "github.com/jaypipes/sqlb"
)

const (
    DSN = "user:password@/blog"
)

var (
    db *sql.DB
    meta *sqlb.Meta
)

func main() {
    if db, err := sql.Open("mysql", DSN); err != nil {
        log.Fatal(err)
    }
    // Grab the schema metadata for tables and columns
    if err = sqlb.Reflect("mysql", db, meta); err != nil {
        log.Fatal(err)
    }
}
```

The `meta` variable will now be populated with metadata about the database
represented by the `db` variable. This metadata will be essential for
constructing SQL expressions using the `sqlb` library.

## Constructing SQL expressions

The idea behind the `sqlb` library is to allow the user of the library to build
up a SQL expression instead of constructing a string variable.

The `sqlb.Select()` function constructs a `SELECT ... FROM` SQL expression. It
takes one or more variadic arguments (tables, columns, or SQL functions) that
the user wants to "project" in the `SELECT` statement.

**Note:**: Projections are all columns or column-like things that end up in the
`<projections>` part of the `SELECT <projections> FROM <sources>` expression.

Our "blog" database has only a single table, called "articles", at the moment:


```sql
CREATE TABLE articles (
  id INT NOT NULL,
  title VARCHAR(200) NOT NULL,
  published_on DATETIME NULL,
  content TEXT NOT NULL,
  PRIMARY KEY (id),
  INDEX ix_title (title)
);
```

Let's suppose that our application wants to output the ten newest articles in
our blog.

We might create a simple function `getRecentArticles()` like so:

```go
struct Article {
    Id int64
    Title string
    PublishedOn uint64
    Content string
}

func getRecentArticles() []*Article {
    // Get last 10 articles from our DB
    articles := meta.Table("articles")
    q := sqlb.Select(articles)
    q.OrderBy(articles.Column("published_on").Desc())
    q.Limit(10)

    qs, qargs := q.StringArgs()
    articles := make([]*Article, 0)
    rows, err := db.Query(qs, qargs...)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        article := &Article{}
        err := rows.Scan(&article.Id, &article.Title,
                         &article.PublishedOn, &article.Content)
        if err != nil {
            log.Fatal(err)
        }
        articles = append(articles, article)
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
    return articles
}
```

Let's analyze the above code line by line.

```go
    articles := meta.Table("articles")
```

The above line shows how the `sqlb.Meta` struct we populated above with
`sqlb.Reflect()` can be queried for a particular `sqlb.Table` struct that
represents an underlying database table (in this case, our only table
"articles").

```go
    q := sqlb.Select(articles)
```

We call the `sqlb.Select()` function to return a `sqlb.SelectQuery` struct. We
pass the `articles` variable (a `sqlb.Table` struct) to the `sqlb.Select()`
function. Doing this automatically adds all columns from the "articles"
database table into the `SELECT` SQL expression represented by the
`sqlb.SelectQuery` struct.

```go
    q.OrderBy(articles.Column("published_on").Desc())
```

The above instructs `sqlb` to add an `ORDER BY` clause to the SQL expression,
using the "published\_on" column as a sort column and ordering the sort in
reverse.

```go
    q.Limit(10)
```

We indicate that the SQL expression should limit the resulting rows to only the
top 10.

```go
    qs, qargs := q.StringArgs()
```

Finally, ask the `sqlb.SelectQuery` to construct a string representing the
`SELECT` statement and return a `[]interface{}` representing the interpolated
arguments that the query uses.

The `qs` variable would contain the following string:

```sql
SELECT id, title, published_on, content FROM articles ORDER BY published_on DESC LIMIT ?
```

and the `qargs` variable would be the following:

```go
[]interface{}{10}
```

The value of `10` is the interpolated query parameter for the question mark in
the `LIMIT` clause.

**Note**: It's important to note that up until now, our code has not needed to
allocate any byte arrays, strings or arrays of `interface{}`. The
`sqlb.SelectQuery.StringArgs()` method allocates a `[]byte` of an exact length
needed for the SQL query and a `[]interface{}` of an exact size needed for
interpolated query parameters.

The remaining lines of code:

```go
    articles := make([]*Article, 0)
    rows, err := db.Query(qs, qargs...)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        article := &Article{}
        err := rows.Scan(&article.Id, &article.Title,
                         &article.PublishedOn, &article.Content)
        if err != nil {
            log.Fatal(err)
        }
        articles = append(articles, article)
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
    return articles
```

simply passes the SQL string and array of `interface{}` generated by the call
to `sqlb.SelectQuery.StringArgs()` to the `database/sql:DB` struct and populate
the returned list of `Article` structs as you would find in any Golang
application that manually constructed SQL strings.
