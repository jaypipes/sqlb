#!/usr/bin/env bash

this_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
scripts_dir="$this_dir/.."
lib_dir="$scripts_dir/lib"

# debug_enabled returns 0 if the DEBUG environs variable is set to anything
# other than "0", 1 otherwise.
debug_enabled() {
  if [ -z "$DEBUG" ]; then
    return 1
  fi
  if [[ "$DEBUG" != "0" ]]; then
    return 0
  else
    return 1
  fi
}

# check_is_installed checks to see if the supplied executable is installed and
# exits if not. An optional second argument is an extra message to display when
# the supplied executable is not installed.
#
# Usage:
#
#   check_is_installed PROGRAM [ MSG ]
#
# Example:
#
#   check_is_installed kind "You can install kind with the helper scripts/install-kind.sh"
check_is_installed() {
  local __name="$1"
  local __extra_msg="$2"
  if ! is_installed "$__name"; then
    echo "FATAL: Missing requirement '$__name'" >&2
    echo "Please install $__name before running this script." >&2
    if [[ -n $__extra_msg ]]; then
      echo "" >&2
      echo "$__extra_msg" >&2
      echo "" >&2
    fi
    exit 1
  fi
}

# is_installed returns 0 if the supplied program is installed on the system,
# 1 otherwise
#
# Usage:
#
#   is_installed PROGRAM
#
# Example:
#
#   if is_installed $PROGRAM; then
#     echo "$PROGRAM is installed."
#   fi
is_installed() {
  local __name="$1"
  if command -v "$__name" >/dev/null 2>&1; then
    return 0
  else
    return 1
  fi
}
