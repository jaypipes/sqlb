CREATE TABLE article_states (
  id INT NOT NULL
, name VARCHAR(20) NOT NULL
, PRIMARY KEY (id)
);

INSERT INTO article_states (id, name) VALUES
  (1, 'draft')
, (2, 'published')
, (3, 'archived')
;

CREATE TABLE users (
  id INT NOT NULL
, name VARCHAR(100) NOT NULL
, created_on TIMESTAMP NOT NULL
, PRIMARY KEY (id)
);

CREATE UNIQUE INDEX ix_name ON users (name);

INSERT INTO users (id, name, created_on) VALUES
  (1, 'Alice', '2018-01-18')
, (2, 'Bob', '2018-02-09')
, (3, 'Charlie', '2019-12-11')
;

CREATE TABLE articles (
  id INT NOT NULL
, title VARCHAR(200) NOT NULL
, author INT NOT NULL REFERENCES users (id)
, state INT NOT NULL
, created_on TIMESTAMP NOT NULL
, PRIMARY KEY (id)
);

CREATE INDEX ix_title ON articles (title);

INSERT INTO articles (id, title, author, state, created_on) VALUES
  (1, 'Alice''s list of grievances', 1, 2, '2018-01-19')
, (2, 'Bob''s list of accomplishments', 2, 1, '2018-06-19')
, (3, 'Charlie''s list of statements', 3, 3, '2019-12-13')
;
