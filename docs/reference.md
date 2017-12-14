# `sqlb` library reference

This document has in-depth information for users of the `sqlb` library,
including examples of constructing complex SQL expressions in a variety of
ways.

1. [Schema and Metadata](#schema-and-metadata)
    1. [Manually specifying metadata](#manually-specifying-metadata)
    1. [Automatically discovering metadata](#automatically-discovering-metadata)
    1. [SQL Dialects](#sql-dialects)
1. [Modifying data](#modifying-data)
    1. [Inserting new rows](#inserting-data-into-the-database)
    1. [Deleting rows](#deleting-data-from-a-table)
    1. [Updating rows](#updating-data-in-a-table)
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
`sqlb.Meta` using the `sqlb.NewMeta()`function, passing in the database dialect
and the name of the database schema:

```go
import (
    "github.com/jaypipes/sqlb"
)

func main() {
    meta := sqlb.NewMeta(sqlb.DIALECT_MYSQL, "blog")
}
```

The `sqlb.Meta` struct has a method `Table(string)` which returns a pointer
to a `sqlb.Table` struct for a table whose name matches the supplied string
argument. If no matching table is found, `nil` is returned:

```go
    users := meta.Table("users")
    // users == nil
```

Since we have yet not told the `meta` struct about any tables in our database,
the `users` variable above will be `nil`. Let's now tell the `meta` struct
about the tables in our schema. A `sqlb.Table` struct contains metadata
about a specific table in a database. The `sqlb.Meta.NewTable()` method
returns a pointer to a `sqlb.Table` struct that's been initialized with the
name of the table.

```go
    users = meta.NewTable("users")
```

A similar process is used to set metadata about a table's column definitions.
The `sqlb.Column` struct defines the column name and links a pointer to
the appropriate `Table` struct. To look up a particular column by its name,
use the `sqlb.Table.C()` method. If no such column is known, `nil`
will be returned:

```go
    colUserId := users.C("id")
    // colUserId == nil
```

Use the `sqlb.Table.NewColumn()` method to create and return a new column
definition:

```go
    colUserId = users.NewColumn("id")
```

### Automatically discovering metadata

The other method of establishing database metadata is to let `sqlb` do it for
you. The `sqlb.Reflect()` function accepts three arguments: the database
dialect in use, a pointer to a `database/sql:DB` struct and a pointer to a
`sqlb.Meta` struct that you wish to fill with metadata information.

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
    if err := sqlb.Reflect(sqlb.DIALECT_MYSQL, db, meta); err != nil {
        log.Fatal(err)
    }

    // meta variable is now populate with table and column information for the
    // "blog" database. Here's some example code that loops through the table
    // metadata printing out table and column names.
    for _, td := range meta.Tables() {
        fmt.Printf("Table: %s\n", td.name)
        for _, cd := range td.Columns() {
            fmt.Printf(" Column: %s", cd.name)
        }
    }
}
```

### SQL Dialects

`sqlb` supports outputting multiple SQL dialects. The two currently-supported
SQL dialects are the MySQL and PostgreSQL dialects.

When constructing SQL expressions with `sqlb`, you will make use of a
`sqlb.Meta` struct, typically by referencing tables and columns. When you
create a `sqlb.Meta` struct using either the `sqlb.NewMeta()` or
`sqlb.Reflect()` functions, the first argument you will supply to those
functions is of type `sqlb.Dialect`, which is a type-alias for an `int`. The
two dialects currently supported are:

* `sqlb.DIALECT_MYSQL`
* `sqlb.DIALECT_POSTGRESQL`

Pass either of those values as the first parameter to `sqlb.NewMeta()` or
`sqlb.Reflect()` and the SQL string generated by `sqlb` will use that database
server's particular dialect. Dialect affects things like how SQL functions are
output -- for example, `TRIM()` vs `BTRIM()` -- as well as other SQL language
constructs.

## Modifying data

The `INSERT`, `DELETE` and `UPDATE` SQL statements are used to add, remove and
modify records in a table in an RDBMS. The following sections provide example
code showing how to use various `sqlb` functions to modify data in your
database.

### Inserting data into the database

Data is added to a table in an RDBMS via the `INSERT` SQL statement, which
follows the form:

```
INSERT INTO <table> (<column_list>) VALUES (<value_list>)
```

While you could manually build up the string for an `INSERT` statement, you can
also use the `sqlb` library to generate that string for you.

The `Insert()` function accepts two arguments: a pointer to a `Table` struct
and a map containing the values to be inserted into the table. The keys of the
map should correspond to the names of the columns the value should be inserted
for.

Imagine I want to add a new record to the `users` table for a new author named
"Fred Flintstone". I might come up with a function named `AddAuthor()` that
looks something like this:

```go
import (
    "database/sql"
    "time"

    "github.com/jaypipes/sqlb"
)

var (
    db *sql.DB
    meta *sqlb.Meta
)

func AddAuthor(name string, email string, profile string) error {
    values := map[string]interface{}{
        "name": name,
        "email": email,
        "profile": profile,
        "is_author": 1,
        "created_on": time.Now().UTC(),
    }
    q := sqlb.Insert(meta.Table("users"), values)
    qs, qargs := q.StringArgs()

    tx, err := db.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    stmt, err := tx.Prepare(qs)
    if err != nil {
        err
    }
    defer stmt.Close()
    res, err := stmt.Exec(qargs...)
    if err != nil {
        return err
    }
    return nil
}
```

The variable `q` in the above `AddAuthor()` function would contain a pointer to
a `sqlb.InsertQuery` struct. This struct's `StringArgs()` method returns two
variables: a string and a slice of `interface{}`. The string (the `qs` variable
above) contains the following SQL query:

```sql

INSERT INTO users (name, email, profile, is_author, created_on) VALUES (?, ?, ?, ?, ?)
```

the slice of `interface{}` (the `qargs` variable above) contains the query
parameters that back the "?" interpolation markers in the SQL string.

The `qs` and `qargs` variables are then used in the calls to the `Tx.Prepare()`
and `Prepare()` method of the returned value from `Tx.Prepare()`.

### Deleting data from a table

Rows in a table may be deleted from an RDBMS using the `DELETE` SQL statement, which has the form:

```
DELETE FROM <table>[ WHERE <predicate>]
```

The `sqlb.Delete()` function can be used to generate a `DELETE` SQL statement.
`sqlb.Delete()` takes a single argument representing the target table to delete
rows from. The function returns a pointer to a `DeleteQuery` struct.

The `sqlb.DeleteQuery` struct has a `Where()` method that can be used to add
the optional `WHERE` clause to the `DELETE` statement, as this example
demonstrates:

```go
import (
    "database/sql"
    "time"

    "github.com/jaypipes/sqlb"
)

var (
    db *sql.DB
    meta *sqlb.Meta
)

func DeleteCommentsOlderThan(dt time.Time) error {
    comments := meta.Table("comments")
    q := sqlb.Delete(comments)
    q.Where(sqlb.LessThan(comments.C("created_on"), dt))
    qs, qargs := q.StringArgs()

    tx, err := db.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    stmt, err := tx.Prepare(qs)
    if err != nil {
        err
    }
    defer stmt.Close()
    res, err := stmt.Exec(qargs...)
    if err != nil {
        return err
    }
    return nil
}
```

The variable `q` in the `DeleteCommentsOlderThan()` function above will contain
a pointer to a `DeleteQuery` struct. We call the `DeleteQuery.Where()` method
to winnow the records to delete from the comments table to only those comments
that were created before the supplied `dt` parameter.

Calling `q.StringArgs()` results in two variables being populated: `qs` and
`qargs`. The `qs` is the SQL string we will pass to our transaction context
(`tx` variable) `Prepare()` method. The SQL produced will be:

```sql
DELETE FROM comments WHERE created_on < ?
```

and the `qargs` variable, which we pass to the `Exec()` method of the prepared
statement will be a slice of `interface{}` and contain a single element: the
`dt` variable.

### Updating data in a table

Use the `sqlb.Update()` function to construct an `UPDATE` SQL statement that
can be executed to modify columnar data within a table in your RDBMS.

The `sqlb.Update()` function takes two arguments. The first argument is a
pointer to a `Table` struct and indicates the target table in the RDBMS. The
second argument is a `map[string]interface{}` that contains the values that you
wish to update for named columns in the target table.

```go
import (
    "database/sql"

    "github.com/jaypipes/sqlb"
)

var (
    db *sql.DB
    meta *sqlb.Meta
)

func UpdateUserProfile(userId int64, profile string) error {
    users := meta.Table("users")
    newVals := map[string]interface{}{
        "profile": profile,
    }
    q := sqlb.Update(users, newVals)
    q.Where(sqlb.Equal(users.C("id"), userId)
    qs, qargs := q.StringArgs()

    tx, err := db.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    stmt, err := tx.Prepare(qs)
    if err != nil {
        err
    }
    defer stmt.Close()
    res, err := stmt.Exec(qargs...)
    if err != nil {
        return err
    }
    return nil
}
```

The `qs` variable in the above code will contain the following string:

```sql
UPDATE users SET profile = ? WHERE users.id = ?
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
in a database. A `sqlb` user can grab a pointer to one of these `sqlb.Table`
structs by calling the `sqlb.Meta.Table()` method, passing in a string for
the table name to get a definition for:

```go
    orgsTable := meta.Table("organizations")
```

Use the `sqlb.Table.As()` method to get a pointer to a `sqlb.Table` struct
that has had its alias set:

```go
    orgs := orgsTable.As("o")
```

If you're sure that a particular table exists in the `sqlb.Meta`, you can
shorten all of the above to just one line:

```go
    orgs := meta.Table("organizations").As("o")
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
    q := sqlb.Select(Count(articles))
    qs, qargs := q.StringArgs()
```

the `qs` variable would contain the following SQL string:

```sql
SELECT COUNT(*) FROM articles
```

You can add an alias to the projected column name for a function using the
`.As()` method, like so:

```go
    q := sqlb.Select(sqlb.Count(articles).As("num_articles"))
```

would produce this SQL string:

```sql
SELECT COUNT(*) AS num_articles FROM articles
```

If you want to count the number of distinct values of a column, use the
`sqlb.CountDistinct()` function:

```go
    articles := meta.Table("articles")
    q := sqlb.Select(sqlb.CountDistinct(articles.C("author")))
    qs, qargs := q.StringArgs()
```

which would produce:

```sql
SELECT COUNT(DISTINCT author) FROM articles
```

#### `Sum()`, `Avg()`, `Min()`, and `Max()`

The `sqlb.Sum()`, `sqlb.Avg()`, `sqlb.`Min()`, and `sqlb.Max()` functions
produce the associated SQL aggregate functions. They all take a single argument
which cam be a `Column`, `Column` or the result of another SQL function,
as these examples show:


```go
    articles := meta.Table("articles")
    q := sqlb.Select(sqlb.Min(articles.C("created_on").As("earliest_article")))
    qs, qargs := q.StringArgs()
```

SQL produced:

```sql
SELECT MIN(created_on) AS earliest_article FROM articles
```

Here's an example that involves a `JOIN` between an `orders` and
`order\_details` table using the `sqlb.Sum()` call to produce an aliased
projection of `SUM(od.amount) AS `total\_amount` in the resulting SQL string.

```go
    o := meta.Table("orders").As("o")
    od := meta.Table("order_details").As("od")
    oId := o.C("id")
    oCustomer := o.C("customer")
    odOrderId := od.C("order_id")
    odAmount := od.C("amount")

    q := sqlb.Select(
        oCustomer,
        sqlb.Sum(odAmount).As("order_total"),
    )
    q.Join(od, sqlb.Equal(oId, odOrderId))
    q.GroupBy(oCustomer)
    qs, qargs := q.StringArgs()
```

would produce the following SQL:

```sql
SELECT o.customer, SUM(od.amount) AS order_total
FROM orders AS o
JOIN order_details AS od
ON o.id = od.order_id
GROUP BY o.customer
```

### String functions

`sqlb` supports a number of SQL functions that operate on string data.

#### `sqlb.Trim()` function variants

* `Trim()` -- Remove whitespace from before and after a string
* `LTrim()` -- Remove whitespace from before a string
* `RTrim()` -- Remove whitespace from after a string
* `TrimChars()` -- Remove specified characters from before and after a string
* `LTrimChars()` -- Remove specified characters from before a string
* `RTrimChars()` -- Remove specified characters from after a string

The SQL function generated by these `sqlb` functions depends on the dialect of
the RDBMS in use. Examples of output are included below for supported dialects.

```go
    users := meta.Table("users")
    colName := users.Column("name")

    sel := sqlb.Select(sqlb.TrimChars(colName, "$#@"))
    qs := sel.String()
```

The following table shows the SQL string that is returned from `sel.String()`
for each `sqlb` string function and the supported RDBMS dialects.

| `sqlb` function | RDBMS dialect | `qs` contents |
| --------------- | ------------- | ------------- |
| `Trim(colName)` | MySQL         | `SELECT TRIM(users.name) FROM users` |
| `Trim(colName)` | PostgreSQL    | `SELECT BTRIM(users.name) FROM users` |
| `LTrim(colName)` | MySQL         | `SELECT LTRIM(users.name) FROM users` |
| `LTrim(colName)` | PostgreSQL    | `SELECT TRIM(LEADING FROM users.name) FROM users` |
| `RTrim(colName)` | MySQL         | `SELECT RTRIM(users.name) FROM users` |
| `RTrim(colName)` | PostgreSQL    | `SELECT TRIM(TRAILING FROM users.name) FROM users` |
| `TrimChars(colName, "#$%")` | MySQL         | `SELECT TRIM(? FROM users.name) FROM users` |
| `TrimChars(colName, "#$%")` | PostgreSQL    | `SELECT TRIM($1 FROM users.name) FROM users` |
| `LTrimChars(colName, "#$%")` | MySQL         | `SELECT TRIM(LEADING ? FROM users.name) FROM users` |
| `LTrimChars(colName, "#$%")` | PostgreSQL    | `SELECT TRIM(LEADING $1 FROM users.name) FROM users` |
| `RTrimChars(colName, "#$%")` | MySQL         | `SELECT TRIM(TRAILING ? FROM users.name) FROM users` |
| `RTrimChars(colName, "#$%")` | PostgreSQL    | `SELECT TRIM(TRAILING $1 FROM users.name) FROM users` |
