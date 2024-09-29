//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package testing_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jaypipes/sqlb"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	envVarMySQLHost          = "MYSQL_HOST"
	envVarMySQLPassword      = "MYSQL_PASSWORD"
	envVarPostgreSQLHost     = "POSTGRESQL_HOST"
	envVarPostgreSQLPassword = "POSTGRESQL_PASSWORD"
)

func skipIfNoMySQL(t *testing.T) {
	if _, ok := os.LookupEnv(envVarMySQLHost); !ok {
		t.Skip("No MySQL container found.")
	}
}

func getMySQLDSN() string {
	host := os.Getenv(envVarMySQLHost)
	pwd := os.Getenv(envVarMySQLPassword)
	return fmt.Sprintf("root:%s@tcp(%s:3306)/sqlbtest", pwd, host)
}

func skipIfNoPostgreSQL(t *testing.T) {
	if _, ok := os.LookupEnv(envVarPostgreSQLHost); !ok {
		t.Skip("No PostgreSQL container found.")
	}
}

func getPostgreSQLDSN() string {
	host := os.Getenv(envVarPostgreSQLHost)
	pwd := os.Getenv(envVarPostgreSQLPassword)
	if pwd != "" {
		pwd = fmt.Sprintf("password=%s", pwd)
	}
	return fmt.Sprintf(
		"host=%s port=5432 user=postgres %s dbname=sqlbtest sslmode=disable",
		host, pwd,
	)
}

func TestReflectMySQL(t *testing.T) {
	skipIfNoMySQL(t)
	db, err := sql.Open("mysql", getMySQLDSN())
	if err != nil {
		log.Fatal(err)
	}
	meta, err := sqlb.Reflect(db)
	if err != nil {
		log.Fatal(err)
	}
	require.Nil(t, err)
	require.NotNil(t, meta)
	assert.Equal(t, 3, len(meta.Tables))
	assert.Equal(t, "sqlbtest", meta.Name)

	users := meta.T("users")
	require.NotNil(t, users)
	userCols := users.Projections()
	assert.Equal(t, 3, len(userCols))

	userID := users.C("id")
	require.NotNil(t, userID)
	assert.Equal(t, "id", userID.Name())
	assert.Equal(t, users, userID.References())
}

func TestReflectPostgreSQL(t *testing.T) {
	skipIfNoPostgreSQL(t)
	db, err := sql.Open("postgres", getPostgreSQLDSN())
	if err != nil {
		log.Fatal(err)
	}
	meta, err := sqlb.Reflect(db)
	if err != nil {
		log.Fatal(err)
	}
	require.Nil(t, err)
	require.NotNil(t, meta)
	assert.Equal(t, 3, len(meta.Tables))
	assert.Equal(t, "sqlbtest", meta.Name)

	users := meta.T("users")
	require.NotNil(t, users)
	userCols := users.Projections()
	assert.Equal(t, 3, len(userCols))

	userID := users.C("id")
	require.NotNil(t, userID)
	assert.Equal(t, "id", userID.Name())
	assert.Equal(t, users, userID.References())
}
