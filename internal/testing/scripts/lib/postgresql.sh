#!/usr/bin/env bash

DEFAULT_POSTGRESQL_CONTAINER_NAME="sqlb-test-postgresql"
DEFAULT_POSTGRESQL_IMAGE="docker.io/postgres"
DEFAULT_POSTGRESQL_IMAGE_VERSION="17.4"
DEFAULT_POSTGRESQL_PASSWORD="mysecretpassword"

this_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
scripts_dir="$this_dir/.."
lib_dir="$scripts_dir/lib"

source "$lib_dir/container.sh"
source "$lib_dir/print.sh"

# postgresql::start starts a container (in daemon mode) running postgresql
#
# Usage:
#
#   postgresql::start [CONTAINER_NAME] [DATA_DIR] [BIND_ADDRESS] [VERSION] [PASSWORD]
#
#   CONTAINER_NAME: (optional) name for the container
#     Default: sqlbtestpostgresql
#   DATA_DIR: (optional) path to directory to use for postgresql state
#     Default: A tmpdir is created (/tmp/postgresql-XXXXXX)
#   BIND_ADDRESS: (optional) bind address that postgresql will use within the
#   container
#     Default: 0.0.0.0
#   IMAGE: (optional) the Docker Image ID of PostgreSQL to run
#     Default: docker.io/postgres
#   VERSION: (optional) the version of PostgreSQL to run
#     Default: 17.4
#   PASSWORD: (optional) the password for the 'postgres' user
#     Default: mysecretpassword
#
# Usage:
#
#   # Start a container called "postgresql-example" running postgresql using
#   # /opt/postgresql-data as the directory to save postgresql's state
#   postgresql::start "postgresql-example" /opt/postgresql-data
postgresql::start() {
  local __container_name="${1:-$DEFAULT_POSTGRESQL_CONTAINER_NAME}"
  local __data_dir="${2:-$(mktemp -d -t postgresql-XXXXXX)}"
  if [ ! -d $__data_dir ]; then
    echo "ERROR: cannot start postgresql container. Supplied data_dir $__data_dir does not exist." >&2
    return 1
  fi
  local __node_addr="${3:-"0.0.0.0"}"

  if container::is_running "$__container_name"; then
    print::info "postgresql container '$__container_name' already running"
    return 0
  fi

  local __postgresql_image="${4:-$DEFAULT_POSTGRESQL_IMAGE}"
  local __postgresql_image_version="${5:-$DEFAULT_POSTGRESQL_IMAGE_VERSION}"
  local __postgresql_password="${6:-$DEFAULT_POSTGRESQL_PASSWORD}"

  if [ -z "$(docker images -q $__postgresql_image:$__postgresql_image_version)" ]; then
    print::inline_first "postgresql image '$__postgresql_image:$__postgresql_image_version' not found. pulling..."
    docker pull "$__postgresql_image:$__postgresql_image_version" >/dev/null && print::ok || (print::fail && return 1)
  fi

  print::inline_first "starting postgresql container '$__container_name' ..."
  docker run -d \
    --rm \
    -p 54320:5432 \
    --volume="$__data_dir":/var/lib/postgresql/data \
    --name "$__container_name" \
    -e POSTGRES_PASSWORD=$__postgresql_password \
    "$__postgresql_image:$__postgresql_image_version" >/dev/null 2>&1
  if [ $? -eq 0 ]; then
    print::ok
  else
    echo "failed to start postgresql container."
    return 1
  fi

  # sleep just a bit to make sure the port forwarding is visible to the host
  sleep 2
}


# postgresql::stop stops the named container running postgresql
#
# Usage:
#
#   postgresql::stop [CONTAINER_NAME]
#
#   CONTAINER_NAME: (optional) name for the container
#     Default: sqlbtestpostgresql
#
# Usage:
#
#   postgresql::stop "postgresql-example"
postgresql::stop() {
  local __container_name="${1:-$DEFAULT_POSTGRESQL_CONTAINER_NAME}"
  print::inline_first "stopping postgresql container '$__container_name' ..."
  container::stop "$__container_name"
  print::ok
}
