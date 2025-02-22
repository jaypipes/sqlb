#!/usr/bin/env bash
DEBUG=${DEBUG:-0}

this_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
schema_dir="$this_dir/../schema"
scripts_dir="$this_dir"
lib_dir="$scripts_dir/lib"

source "$lib_dir/common.sh"

check_is_installed docker
check_is_installed mysql
check_is_installed psql

source "$lib_dir/container.sh"
source "$lib_dir/mysql.sh"
source "$lib_dir/postgresql.sh"

default_test_dbname="sqlbtest"

usage="Usage:
  $(basename "$0") [-dh]

Options:
  -d                                    Enables debug output.
  -h                                    Shows this help.

Environment variables:
  DEBUG:                                Toggles debug output. Set to any
                                        non-zero value to enable.
                                        Default: 0
  MYSQL_HOST                            Informs the script of a MySQL server to
                                        use for testing. If not set (the
                                        default), a Docker container running
                                        MySQL will automatically be started.
  MYSQL_ROOT_PASSWORD                   The root user password for the MySQL
                                        server.
                                        Default: ''
  MYSQL_CONTAINER_NAME                  Name of the Docker container to run
                                        MySQL. If MYSQL_HOST is supplied, this
                                        is ignored.
                                        Default: '$DEFAULT_MYSQL_CONTAINER_NAME'
  POSTGRESQL_HOST                       Informs the script of a PostgreSQL
                                        server to use for testing. If not set
                                        (the default), a Docker container
                                        running PostgreSQL will automatically
                                        be started.
  POSTGRESQL_PASSWORD                   The postgres user password for the
                                        PostgreSQL server.
                                        Default: 'mysecretpassword'
  POSTGRESQL_CONTAINER_NAME             Name of the Docker container to run
                                        PostgreSQL. If POSTGRESQL_HOST is
                                        supplied, this is ignored.
                                        Default: '$DEFAULT_POSTGRESQL_CONTAINER_NAME'
  TEST_DBNAME                           Name of the database/schema to use for
                                        testing.
                                        Default: '$default_test_dbname'
"

while getopts ":hd" opt; do
  case ${opt} in
    d)
      DEBUG=1
      ;;
    h)
      echo "$usage"
      exit 0
      ;;
    ?)
      echo "Invalid option: -${OPTARG}."
      echo "$usage"
      exit 1
      ;;
  esac
done

test_dbname="${TEST_DBNAME:-$default_test_dbname}"
mysql_container_name="${MYSQL_CONTAINER_NAME:-$DEFAULT_MYSQL_CONTAINER_NAME}"
mysql_host="${MYSQL_HOST}"
mysql_root_password="${MYSQL_ROOT_PASSWORD}"

if [ -z $mysql_host ]; then
  if ! container::is_running "$mysql_container_name"; then
    $scripts_dir/mysql_start.sh "$mysql_container_name"
  else
    print::info "mysql container '$mysql_container_name' already running"
  fi
  
  if ! container::get_ip "$mysql_container_name" mysql_container_ip; then
    echo "ERROR: could not get IP for mysql container"
    exit 1
  fi
  mysql_host="$mysql_container_ip"
fi

mysql_root_password_arg=""
if [ ! -z $mysql_root_password ]; then
  mysql_root_password_arg="--password=$mysql_root_password"
fi

print::inline_first "Dropping mysql test $test_dbname database ... "
mysql -uroot -P3306 $mysql_root_password_arg -h$mysql_host --protocol=tcp -e "DROP DATABASE IF EXISTS $test_dbname;"
if [ $? -eq 0 ]; then
  print::ok
else
  print::fail
fi

print::inline_first "Creating mysql test $test_dbname database ... "
tmpsql_path=$(mktemp)
dbname=$test_dbname envsubst < $schema_dir/mysql.sql > $tmpsql_path
mysql -uroot -P3306 $mysql_root_password_arg -h$mysql_host --protocol=tcp < $tmpsql_path
if [ $? -eq 0 ]; then
  print::ok
else
  print::fail
fi

postgresql_container_name="${POSTGRESQL_CONTAINER_NAME:-$DEFAULT_POSTGRESQL_CONTAINER_NAME}"
postgresql_host="${POSTGRESQL_HOST}"

if [ -z $postgresql_host ]; then
  if ! container::is_running "$postgresql_container_name"; then
    $scripts_dir/postgresql_start.sh "$postgresql_container_name"
  else
    print::info "postgresql container '$postgresql_container_name' already running"
  fi
  
  if ! container::get_ip "$postgresql_container_name" postgresql_container_ip; then
    echo "ERROR: could not get IP for postgresql container"
    exit 1
  fi
  postgresql_host="$postgresql_container_ip"
fi

postgresql_password="${POSTGRESQL_PASSWORD:-${POSTGRES_PASSWORD:-$DEFAULT_POSTGRESQL_PASSWORD}}"

print::inline_first "Dropping postgresql test $test_dbname database ... "
PGOPTIONS='--client-min-messages=warning' psql postgres://postgres:$postgresql_password@$postgresql_host:5432 -c "DROP DATABASE IF EXISTS $test_dbname;" >/dev/null 2>&1
if [ $? -eq 0 ]; then
  print::ok
else
  print::fail
fi

print::inline_first "Creating postgresql test $test_dbname database ... "
psql postgres://postgres:$postgresql_password@$postgresql_host:5432 -tc "SELECT 1 FROM pg_database WHERE datname = '$test_dbname'" | \
    grep -q 1 | \
    psql postgres://postgres:$postgresql_password@$postgresql_host:5432 -c "CREATE DATABASE $test_dbname" >/dev/null 2>&1
tmpsql_path=$(mktemp)
dbname=$test_dbname envsubst < $schema_dir/postgresql.sql > $tmpsql_path
psql postgres://postgres:$postgresql_password@$postgresql_host:5432/$test_dbname < $tmpsql_path >/dev/null 2>&1
if [ $? -eq 0 ]; then
  print::ok
else
  print::fail
fi
