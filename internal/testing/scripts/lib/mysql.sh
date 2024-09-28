#!/usr/bin/env bash

DEFAULT_MYSQL_CONTAINER_NAME="sqlb-test-mysql"
DEFAULT_MYSQL_IMAGE_VERSION="8.39.0"

this_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
scripts_dir="$this_dir/.."
lib_dir="$scripts_dir/lib"

source "$lib_dir/container.sh"
source "$lib_dir/print.sh"

# mysql::start starts a container (in daemon mode) running mysql
#
# Usage:
#
#   mysql::start [CONTAINER_NAME] [DATA_DIR] [BIND_ADDRESS]
#
#   CONTAINER_NAME: (optional) name for the container
#     Default: sqlbtestmysql
#   DATA_DIR: (optional) path to directory to use for mysql state
#     Default: A tmpdir is created (/tmp/mysql-XXXXXX)
#   BIND_ADDRESS: (optional) bind address that mysql will use within the
#   container
#     Default: 0.0.0.0
#   VERSION: (optional) the version of MySQL to run
#     Default: 8.39.0
#
# Usage:
#
#   # Start a container called "mysql-example" running mysql using
#   # /opt/mysql-data as the directory to save mysql's state
#   mysql::start "mysql-example" /opt/mysql-data
mysql::start() {
  local __container_name="${1:-$DEFAULT_MYSQL_CONTAINER_NAME}"
  local __data_dir="${2:-$(mktemp -d -t mysql-XXXXXX)}"
  if [ ! -d $__data_dir ]; then
    echo "ERROR: cannot start mysql container. Supplied data_dir $__data_dir does not exist." >&2
    return 1
  fi
  local __node_addr="${3:-"0.0.0.0"}"

  if container::is_running "$__container_name"; then
    print::info "mysql container '$__container_name' already running"
    return 0
  fi

  mysql_image_version="${4:-$DEFAULT_MYSQL_IMAGE_VERSION}"

  print::inline_first "starting mysql container '$__container_name' ..."
  docker run -d \
    --rm \
    -p 33060:3306 \
    --volume="$__data_dir":/var/lib/mysql \
    --name "$__container_name" \
    -e MYSQL_ALLOW_EMPTY_PASSWORD=1 \
    -e MYSQL_ROOT_HOST="%" \
    mysql/mysql-server:$mysql_image_version \
    --bind-address "$__node_addr" >/dev/null 2>&1
  if [ $? -eq 0 ]; then
    print::ok
  else
    echo "failed to start mysql container."
    return 1
  fi
    

  # We need to wait until we see "/usr/sbin/mysqld: ready for connections" in
  # the docker logs *and* the server was started on port 3306. There are two
  # times that "mysqld: ready for connections" will appear. The first is when
  # the server is initially started on port 0 during bootstrapping and the
  # second is the "normal" mysqld startup process. If we don't do this, we'll
  # just get connection errors when trying to connect to the server :(

  print::inline_first "waiting for mysql to be ready for connections ..."
  local __sleep_time=0
  ready=0
  until [ $__sleep_time -eq 120 ]; do
    sleep 4
    __sleep_time=$(( __sleep_time + 3 ))
    found=$(docker logs --tail 20 "$__container_name" 2>&1 | grep -A2 "mysqld: ready for connections" | grep "port: 3306")
    if [ $? -eq 0 ]; then
      ready=0
      break
    else
      print::inline "."
    fi
  done

  if [ $ready -eq 0 ]; then
    print::ok
  else
    print::fail
    echo "failed to detect mysql ready for connections after 2 minutes."
    return 1
  fi
}


# mysql::stop stops the named container running mysql
#
# Usage:
#
#   mysql::stop [CONTAINER_NAME]
#
#   CONTAINER_NAME: (optional) name for the container
#     Default: sqlbtestmysql
#
# Usage:
#
#   mysql::stop "mysql-example"
mysql::stop() {
  local __container_name="${1:-$DEFAULT_MYSQL_CONTAINER_NAME}"
  print::inline_first "stopping mysql container '$__container_name' ..."
  container::stop "$__container_name"
  print::ok
}
