# `sqlb` library reference

This document has in-depth information for users of the `sqlb` library,
including examples of constructing complex SQL expressions in a variety of
ways.

1. [Schema and Metadata](#schema-and-metadata)
    1. [Manually specifying metadata](#manually-specifying-metadata)
    1. [Automatically discovering metadata](#automatically-discovering-metadata)
1. [Aliasables](#aliasables)
1. [SQL Functions](#sql-functions)

## Schema and Metadata

When constructing SQL expressions with the `sqlb` library, you will almost
always end up referring to information about the underlying database schema,
e.g. its table and column definitions. There are two primary methods by which a
`sqlb` user can reference this metadata information: **manually** or via
**reflection**.

For these examples, we will be assuming a MySQL database called "blog" has been
created with the following database schema:

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

### Manually specifying metadata

If you have a relatively small and/or static database schema, you might want to
specify metadata about that database schema by manually constructing a
collection of metadata elements. This section demonstrates how you can do so.

The `sqlb.Meta` struct can be used to house metadata about a database schema's
tables, those tables' columns and relation information. You can create a new
`sqlb.Meta` using the `sqlb.NewMeta()`function, passing in the name of the
database driver and the name of the database schema:

```go
import (
    "github.com/jaypipes/sqlb"
)

func main() {
    meta := sqlb.NewMeta("mysql", "blog")
}
```

The `sqlb.Meta` struct has a method `TableDef(string)` which returns a pointer
to a `sqlb.TableDef` struct for a table whose name matches the supplied string
argument. If no matching table is found, `nil` is returned:

```go
    users := meta.TableDef("users")
    // users == nil
```

Since we have yet not told the `meta` struct about any tables in our database,
the `users` variable above will be `nil`. Let's now tell the `meta` struct
about the tables in our schema. A `sqlb.TableDef` struct contains metadata
about a specific table in a database. The `sqlb.Meta.NewTable()` method
returns a pointer to a `sqlb.TableDef` struct that's been initialized with the
name of the table.

```go
    users = meta.NewTableDef("users")
```

A similar process is used to set metadata about a table's column definitions.
The `sqlb.ColumnDef` struct defines the column name and links a pointer to
the appropriate `TableDef` struct. To look up a particular column by its name,
use the `sqlb.TableDef.ColumnDef()` method. If no such column is known, `nil`
will be returned:

```go
    colUserId := users.ColumnDef("id")
    // colUserId == nil
```

Use the `sqlb.TableDef.NewColumnDef()` method to create and return a new column
definition:

```go
    colUserId = users.NewColumnDef("id")
```

### Automatically discovering metadata

The other method of establishing database metadata is to let `sqlb` do it for
you. The `sqlb.Reflect()` function accepts three arguments: a string describing
the `database/sql` driver name, a pointer to a `database/sql:DB` struct and a
pointer to a `sqlb.Meta` struct that you wish to fill with metadata
information.

The following code demonstrates how to use the `sqlb.Reflect()` function
properly. We use a MySQL database instance, however simply change the
DB-specific driver import and driver name passed to `sqlb.Reflect()` to use a
different database server.

```go
import (
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var (
    db *sql.DB
    meta *sqlb.Meta
)

func main() {
    // First, set up your database/sql:DB struct using database/sql:Open()
    if db, err := sql.Open("mysql", "user:pass@/blog"); err != nil {
        log.Fatal(err)
    }

    // Next, ask sqlb.Reflect() to populate the metadata for the DB
    if err := sqlb.Reflect("mysql", db, meta); err != nil {
        log.Fatal(err)
    }

    // meta variable is now populate with table and column information for the
    // "blog" database. Here's some example code that loops through the table
    // metadata printing out table and column names.
    for _, td := range meta.TableDefs() {
        fmt.Printf("Table: %s\n", td.name)
        for _, cd := range td.ColumnDefs() {
            fmt.Printf(" Column: %s", cd.name)
        }
    }
}
```

## Aliasables

When constructing SQL expressions, it's often useful to provide an alias for a
database object with a long name. For instance, let's say I have a table called
`organizations` and I want to `SELECT` the `id` and `display_name` columns from
that table.

In SQL, I would write:

```sql
SELECT organizations.id, organizations.display_name
FROM organizations
```

The SQL language allows the user to specify an **alias** for an object using
the `AS` keyword. So, if I didn't feel like typing out the word "organizations"
over and over again, I could alias the organizations table object as "o", like
so:

```sql
SELECT o.id, o.display_name
FROM organizations AS o
```

Shorter and cleaner, for sure.

In `sqlb`, there are a number of objects that are **aliasable**:

* `Column` structs
* `Table` structs
* `Function` structs

Aliasable things may have an alias applied to them with the `SetAlias()` struct
method:

```go
    t := meta.Table("organizations")
    t.SetAlias("o")
```

if I use the `t` variable in a `sqlb.Select()` call, the SQL output will
automatically include the table alias for any columns that reference the table.

That said, there's a more convenient way to get a pointer to an aliased
database object: using the `As()` method on another struct.

As you learned earlier, the `sqlb.Meta` struct contains definitions of tables
in a database. A `sqlb` user can grab a pointer to one of these `sqlb.TableDef`
structs by calling the `sqlb.Meta.TableDef()` method, passing in a string for
the table name to get a definition for:

```go
    orgsTableDef := meta.TableDef("organizations")
```

Use the `sqlb.TableDef.As()` method to get a pointer to a `sqlb.Table` struct
that has had its alias set:

```go
    orgs := orgsTableDef.As("o")
```

If you're sure that a particular table exists in the `sqlb.Meta`, you can
shorten all of the above to just one line:

```go
    orgs := meta.TableDef("organizations").As("o")
```

Short and sweet.

## SQL Functions

`sqlb` supports a number of common SQL functions.

### Aggregate functions

Aggregate functions apply an operation over a group of records. The following
sections show the `sqlb` library functions, how to use them in your code, and
the eventual SQL string produced when used in query expression.

#### `Count()` and `CountDistinct()`

When you want to express a count of the total number of records in a matching
query, use the `sqlb.Count()` function.

```go
    articles := meta.Table("articles")
    q := Select(Count(articles))
    qs, qargs := q.StringArgs()
```

the `qs` variable would contain the following SQL string:

```sql
SELECT COUNT(*) FROM articles
```

You can add an alias to the projected column name for a function using the
`.As()` method, like so:

```go
    q := Select(Count(articles).As("num_articles"))
```

would produce this SQL string:

```sql
SELECT COUNT(*) AS num_articles FROM articles
```

If you want to count the number of distinct values of a column, use the
`sqlb.CountDistinct()` function:

```go
    articles := meta.Table("articles")
    q := Select(CountDistinct(articles.Column("author")))
    qs, qargs := q.StringArgs()
```

which would produce:

```sql
SELECT COUNT(DISTINCT author) FROM articles
```

#### `Sum()`, `Avg()`, `Min()`, and `Max()`

The `sqlb.Sum()`, `sqlb.Avg()`, `sqlb.`Min()`, and `sqlb.Max()` functions
produce the associated SQL aggregate functions. They all take a single argument
which cam be a `Column`, `ColumnDef` or the result of another SQL function,
as these examples show:


```go
    articles := meta.Table("articles")
    q := Select(Min(articles.Column("created_on").As("earliest_article")))
    qs, qargs := q.StringArgs()
```

SQL produced:

```sql
SELECT MIN(created_on) AS earliest_article FROM articles
```

