//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

// Package grammar contains structs that comprise the ANSI 2003
// SQL syntax along with some structs that represent various SQL language
// extensions implemented by databases like MySQL or PostgreSQL.
//
// All structs in this package are *plain old data* (POD) structs. None of the
// structs have methods and this package contains zero logic for building complex
// SQL statements from these basic structs. Look in [core/expr] and
// [internal/builder] packages for that.
package grammar
