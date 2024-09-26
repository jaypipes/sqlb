#!/usr/bin/env bash

# Used for "resetting" a functional testing environment back to a clean start
# state.
#
# This script stops the following containers if they are running:
#  * sqlb-test-mysql
#
# And then proceeds to clear out the SQL database. We do not attempt to
# stop/start the sqlb-test-mysql container because this container takes a
# stupid long time to start up due to the init scripts used by the MySQL Docker
# container. Instead, we just DROP and re-CREATE the database.

DEBUG=${DEBUG:-0}

this_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
scripts_dir="$this_dir"
lib_dir="$scripts_dir/lib"

default_dbname="sqlbtest"

source "$lib_dir/common.sh"

check_is_installed docker
check_is_installed mysql

source "$lib_dir/container.sh"
source "$lib_dir/mysql.sh"

dbname="${TESTDB_NAME:-$default_dbname}"
mysql_container_name=${MYSQL_CONTAINER_NAME:-"$DEFAULT_MYSQL_CONTAINER_NAME"}

if ! container::is_running "$mysql_container_name"; then
  $scripts_dir/mysql_start.sh "$mysql_container_name"
else
  print::info "mysql container '$mysql_container_name' already running"
fi

if ! container::get_ip "$mysql_container_name" mysql_container_ip; then
  echo "ERROR: could not get IP for mysql container"
  exit 1
fi

print::inline_first "Dropping mysql test $dbname database ... "
mysql -uroot -P3306 -h$mysql_container_ip -e "DROP DATABASE IF EXISTS $dbname;"
print::ok

print::inline_first "Creating mysql test $dbname database ... "
tmpsql_path=$(mktemp)
dbname=$dbname envsubst < $scripts_dir/sqlbtest_mysql.sql > $tmpsql_path
mysql -uroot -P3306 -h$mysql_container_ip < $tmpsql_path
print::ok
