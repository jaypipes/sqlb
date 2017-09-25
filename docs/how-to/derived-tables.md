# Working with derived tables

Derived tables are subqueries in the `FROM` clause. They are useful when you
need to join to sets of information that are grouped or ordered differently
than other data in your `SELECT` expression.

For example, let's say you have the following table in your database schema
representing comments that readers have left on your blog:

```sql
CREATE TABLE comments (
  id INT NOT NULL,
  article_id INT NOT NULL,
  title VARCHAR(200) NOT NULL,
  content TEXT NOT NULL,
  created_on DATETIME NOT NULL,
  commenter VARCHAR(200) NOT NULL,
  PRIMARY KEY (id),
  INDEX ix_article_id (article_id)
);
```

Now let's imagine that you wish to show comments, ordered by the time they were
created, for articles written by the three most prolific authors.

We could solve this problem by doing two `SELECT` requests, one that grabs the
IDs of the three most prolific authors:

```sql
SELECT u.id
FROM users AS u
JOIN articles AS a
ON u.id = articles.author
WHERE u.is_author = 1
GROUP u.id
ORDER BY COUNT(*) DESC
LIMIT 3
```

And then do a second `SELECT` query that passes in the returned author IDs to
grab comments, like so:

```sql
SELECT c.*
FROM comments AS c
JOIN articles AS a
ON c.article_id = a.id
WHERE a.author IN ($AUTHORS)
ORDER BY c.created_on DESC
```

However, doing multiple queries is often less efficient than doing a single
query. We can use a `JOIN` to a derived table to accomplish the above in a
single query, like so:

```sql
SELECT c.*
FROM comments AS c
JOIN articles AS a
ON c.article_id = a.id
JOIN (
    SELECT u.id
    FROM users AS u
    JOIN articles AS a
    ON u.id = articles.author
    WHERE u.is_author = 1
    GROUP u.id
    ORDER BY COUNT(*) DESC
    LIMIT 3
) AS top_authors
ON a.author = top_authors.id
ORDER BY c.created_on DESC
```

But, how do we ask `sqlb` to construct such an expression? Well, it's fairly
simple. The `sqlb.SelectQuery` struct that is returned from the `sqlb.Select()`
function can be joined to another `sqlb.SelectQuery`, as this example shows:

```go
u := meta.TableDef("users")
a := meta.TableDef("articles")
c := meta.TableDef("comments")

usersId := u.ColumnDef("id")
articlesId := a.ColumnDef("id")
articlesAuthor := a.ColumnDef("author")
articlesIsAuthor := a.ColumnDef("is_author")
commentsArticleId := c.ColumnDef("article_id")
commentsCreatedOn := c.ColumnDef("created_on")

// First, build the subquery in the FROM clause (the derived table)
subq := sqlb.Select(usersId).Join(a, sqlb.Equal(usersId, articlesAuthor))
subq.Where(sqlb.Equal(articlesIsAuthor, 1))
subq.GroupBy(usersId).OrderBy(Count().Desc())
subq.Limit(3)
subq.As("top_authors")

// Next, build the outer SELECT on the comments table and join to the subselect
q := sqlb.Select(c).Join(a, sqlb.Equal(commentsArticleId, articlesId))
q.Join(subq, sqlb.Equal(articlesAuthor, subq.Column("id")))
q.OrderBy(commentsCreatedOn.Desc())
```
