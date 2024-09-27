#!/usr/bin/env bash

this_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
scripts_dir="$this_dir/.."
lib_dir="$scripts_dir/lib"

source "$lib_dir/print.sh"

# container::get_ip sets a variable with a supplied name to the IP address of a
# named container.
#
# Usage:
#
#   container::get_ip CONTAINER_NAME IP_VAR_NAME
#
#   CONTAINER_NAME: (required) name for the container
#   IP_VAR_NAME: (required) name for variable to store the IP address in
#
# Example:
#
#   # Get the IP address of the container named "etcd-testing" and store that
#   # IP address in a variable named "etcd_container_ip"
#   container::get_ip "etcd-testing" etcd_container_ip
#   echo $etcd_container_ip
container::get_ip() {
  local __container_name="$1"
  local __store_result=$2
  local __sleep_time=0
  local __found_ip=""

  until [ $__sleep_time -eq 8 ]; do
    sleep $(( __sleep_time++ ))
    __found_ip=$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$__container_name")
    if [[ "$__found_ip" != "" ]]; then
      break
    fi
  done
  eval $__store_result="'$__found_ip'"
}

# container::is_running Returns 0 if a container with the given name is
# running, 1 otherwise.
#
# Usage:
#
#   container::is_running CONTAINER_NAME
#
#   CONTAINER_NAME: (required) name for the container
#
# Usage:
#
#   if container::is_running "etcd-example"; then
#     echo "etcd-example container is running."
#   else
#     echo "etcd-example container is not running."
#   fi
container::is_running() {
  local __container_name="$1"

  __running=$(
    docker inspect \
        --format='{{.State.Running}}' \
        "$__container_name" 2>/dev/null
  )
  if [ $? -eq 1 ]; then
    return 1
  fi
  if [ "$__running" = "true" ]; then
    return 0
  fi
  return 1
}

# container::exists returns 0 if the container with the given name exists, 1
# otherwise.
#
# Note that a stopped container still exists on the system. So, while
# container::is_running returns 1 for a stopped container, container::exists
# will return 0 for a stopped container.
#
# Usage:
#
#   container::exists CONTAINER_NAME
#
#   CONTAINER_NAME: (required) name for the container
#
# Usage:
#
#   if container::exists "etcd-example"; then
#     echo "etcd-example container exists."
#   else
#     echo "etcd-example container does not exist."
#   fi
container::exists() {
  local __container_name="$1"

  # Obviously all containers that are running also exist...
  if container::is_running "$__container_name"; then
    return 0
  fi

  __status=$(docker inspect --format='{{.State.Status}}' "$__container_name" 2>/dev/null )
  if [ $? -eq 1 ]; then
    return 1
  fi
  if [ "$__status" = "created" ]; then
    return 0
  fi
  return 1
}

# container::stop gracefully stops the container with the given name but does
# not destroy the container.
#
# Usage:
#
#   container::stop CONTAINER_NAME
#
#   CONTAINER_NAME: (required) name for the container
#
# Usage:
#
#   container::stop "etcd-example"
container::stop() {
  local __container_name="$1"

  if container::is_running "$__container_name"; then
    docker stop "$__container_name" --time 2 1>/dev/null
  fi
}

# container::destroy destroys the container with the given name.
#
# If the container with the given name is active, this function first attempts
# to gracefully stop the container.
#
# Usage:
#
#   container::destroy CONTAINER_NAME
#
#   CONTAINER_NAME: (required) name for the container
#
# Usage:
#
#   container::destroy "etcd-example"
container::destroy() {
  local __container_name="$1"

  container::stop "$__container_name"
  if container::exists "$__container_name"; then
    docker rm "$__container_name"  1>/dev/null
  fi
}

# container::image_exists returns 0 if an image with the given name:tag exists,
# 1 otherwise.
#
# Usage:
#
#   container::image_exists CONTAINER_NAME
#
#   CONTAINER_NAME: (required) name for the container
#
# Usage:
#
#   if container::image_exists "etcd-example"; then
#     echo "etcd-example image exists."
#   fi
container::image_exists() {
  local __image_name_tag="$1"

  docker image inspect "$__image_name_tag" >/dev/null 2>&1
}
