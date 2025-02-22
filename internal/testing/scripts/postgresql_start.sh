#!/usr/bin/env bash

DEBUG=${DEBUG:-0}
scripts_dir=$(cd $(dirname "$0") && pwd)
lib_dir="$scripts_dir/lib"

source $lib_dir/common.sh

check_is_installed docker

source $lib_dir/container.sh
source $lib_dir/postgresql.sh

container_name=${1:-${POSTGRESQL_CONTAINER_NAME:-"$DEFAULT_POSTGRESQL_CONTAINER_NAME"}}

postgresql::start "$container_name"

container::get_ip "$container_name" container_ip
if [ -z "$container_ip" ]; then
  print::error "failed to determine IP address for '$container_name' container."
  exit 1
fi

print::info "postgresql running in container '${container_name}' at ${container_ip}:5432."
