module github.com/jaypipes/sqlb/internal/testing

go 1.21.11

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/jaypipes/sqlb v0.0.0-20240927000255-1416c33ef9fb
	github.com/lib/pq v1.10.9
	github.com/stretchr/testify v1.9.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/samber/lo v1.47.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/jaypipes/sqlb => ../../
