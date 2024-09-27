#!/usr/bin/env bash

this_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
scripts_dir="$this_dir/.."
lib_dir="$scripts_dir/lib"

source "$lib_dir/common.sh"

default_debug_prefix="[debug] "
default_info_prefix="[info] "

# print::debug prints out a supplied message if the DEBUG environs variable is
# set. An optional second argument indicates the "indentation level" for the
# message.
print::debug() {
  if !debug_enabled; then
    return 0
  fi
  local __msg="$1"
  local __indent_level=${2:-}
  local __debug_prefix="${DEBUG_PREFIX:-$default_debug_prefix}"
  _echo_with_indent "$__msg" "$__debug_prefix" no $__indent_level
}

# print::debug_inline prints out a supplied message with no trailing newline if
# the DEBUG environs variable is set.
print::debug_inline() {
  if !debug_enabled; then
    return 0
  fi
  local __msg="$1"
  print::inline "$__msg"
}

# print::info prints out a supplied message. An optional second argument
# indicates the "indentation level" for the message.
print::info() {
  local __msg="$1"
  local __indent_level=${2:-}
  local __info_prefix="${INFO_PREFIX:-$default_info_prefix}"
  _echo_with_indent "$__msg" "$(_info_prefix)" no $__indent_level
}

# print::inline_first prints out a supplied message with no trailing newline
# after prepending the info-level prefix.
print::inline_first() {
  local __msg="$1"
  print::inline "$(_info_prefix)$__msg"
}

# print::debug_inline_first prints out a supplied message with no trailing
# newline after prepending the debug-level prefix if the DEBUG environs
# variable is set.
print::debug_inline_first() {
  local __msg="$1"
  print::debug_inline "$(_debug_prefix)$__msg"
}

# print::inline prints out a supplied message with no trailing newline.
print::inline() {
  local __msg="$1"
  echo -n "$__msg"
}

# print::ok prints "ok." and a newline.
print::ok() {
  echo " ok."
}

# print::fail prints "fail." and a newline.
print::fail() {
  echo " fail."
}

# print::debug_ok prints "ok." and a newline if debugging is on.
print::debug_ok() {
  if !debug_enabled; then
    return 0
  fi
  echo " ok."
}

# print::debug_fail prints "fail." and a newline if debugging is on.
print::debug_fail() {
  if !debug_enabled; then
    return 0
  fi
  echo " fail."
}

_echo_with_indent() {
  local __msg="$1"
  local __prefix="$2"
  local __indent_level=${4:-}
  __indent=""
  if [ -n "$__indent_level" ]; then
    __indent="$( for _ in $( seq 0 "$__indent_level" ); do printf " "; done )"
  fi
  echo "$__prefix$__indent$__msg"
}

_debug_prefix() {
  local __debug_prefix="${DEBUG_PREFIX:-$default_debug_prefix}"
  echo "$(_timestamp) $__debug_prefix"
}

_info_prefix() {
  local __info_prefix="${INFO_PREFIX:-$default_info_prefix}"
  echo "$(_timestamp) $__info_prefix"
}

_timestamp() {
  date '+%H:%M:%S.%3N'
}
