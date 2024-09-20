# `sqlb/grammar` package

The `grammar` package in `sqlb` contains structs that comprise the ANSI 2003
SQL syntax along with some structs that represent various SQL language
extensions implemented by databases like MySQL or PostgreSQL.

All structs in this package are *plain old data* (POD) structs. None of the
structs have methods and this package contains zero logic for building complex
SQL statements from these basic structs. Look in the `api` and
`internal/builder` packages for that.
