#!/usr/bin/env bats
# bats file_tags=unit

setup() {
	load "test_helper/load"
	load_lib "bats-assert"
	load_lib "bats-support"
}

@test "display help with no arguments" {
  run oslo
  assert_success
  assert_output --partial "Usage"
}

@test "has help option" {
  run oslo --help
  assert_success
  assert_output --partial "Usage"
}

@test "has version option" {
  run oslo --version
  assert_success
  assert_output --regexp "^oslo version v\d+\.\d+\.\d+"
}
