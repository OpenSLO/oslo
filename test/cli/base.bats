#!/usr/bin/env bats

setup() {
    load 'test_helper/common-setup'
    _common_setup
}

@test "oslo displays help with no arguments" {
  run oslo
  assert_equal $status 0
  assert_output --partial "Usage"
}

@test "oslo has help option" {
  run oslo --help
  assert_equal $status 0
  assert_output --partial "Usage"
}

@test "oslo has version option" {
  run oslo --version
  assert_equal $status 0
  assert_output
}
