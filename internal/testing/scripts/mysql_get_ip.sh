#!/usr/bin/env bash

DEBUG=${DEBUG:-0}
scripts_dir=$(cd $(dirname "$0") && pwd)
lib_dir="$scripts_dir/lib"

source $lib_dir/common.sh

check_is_installed docker

source $lib_dir/container.sh
source $lib_dir/mysql.sh

container_name=${1:-${MYSQL_CONTAINER_NAME:-"$DEFAULT_MYSQL_CONTAINER_NAME"}}

if container::is_running "$container_name"; then
  if container::get_ip "$container_name" container_ip; then
    echo $container_ip
  fi
fi
