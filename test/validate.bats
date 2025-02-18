#!/usr/bin/env bats
# bats file_tags=unit

setup() {
	load "test_helper/load"
	load_lib "bats-assert"
	load_lib "bats-support"
}

@test "validate an invalid api version" {
  run oslo fmt -f "${TEST_SUITE_INPUTS}/invalid-apiversion.yaml"
  assert_failure
  assert_output "Error: issue parsing objects: error unmarshaling JSON: while decoding JSON: failed to decode object: unsupported openslo.Version: foo"
}
