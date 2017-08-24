# `sqlb` library reference

This document has in-depth information for users of the `sqlb` library,
including examples of constructing complex SQL expressions in a variety of
ways.

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

The `sqlb.Meta` struct has a method `Table(string)` which returns a pointer to
a `sqlb.TableDef` struct for a table whose name matches the supplied string
argument. If no matching table is found, `nil` is returned:

```go
    users := meta.Table("users")
    // users == nil
```

Since we have yet not told the `meta` struct about any tables in our database,
the `users` variable above will be `nil`. Let's now tell the `meta` struct
about the tables in our schema. A `sqlb.TableDef` struct contains metadata
about a specific table in a database. The `sqlb.Meta.NewTable()` method
returns a pointer to a `sqlb.TableDef` struct that's been initialized with the
name of the table.

```go
    users = meta.NewTable("users")
```

### Using `sqlb.Reflect()` to automatically gather metadata

TODO
