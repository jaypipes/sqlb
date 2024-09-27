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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	envVarMySQLContainerIP = "SQLBTEST_MYSQL_IP"
)

func skipIfNoMySQL(t *testing.T) {
	if _, ok := os.LookupEnv(envVarMySQLContainerIP); !ok {
		t.Skip("No MySQL container found.")
	}
}

func getMySQLDSN() string {
	containerIP := os.Getenv(envVarMySQLContainerIP)
	return fmt.Sprintf("root@(%s:3306)/sqlbtest", containerIP)
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
