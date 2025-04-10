#!/usr/bin/env bats
# bats file_tags=unit

setup() {
	load "test_helper/load"
	load_lib "bats-assert"
	load_lib "bats-support"
}

@test "validate an invalid api version" {
  run oslo validate -f "${TEST_SUITE_INPUTS}/invalid-apiversion.yaml"
  assert_failure
  assert_output "Error: failed to read objects from /oslo/test/inputs/invalid-apiversion.yaml: error unmarshaling JSON: while decoding JSON: failed to decode object: unsupported openslo.Version: foo"
}

@test "validate unknown field" {
  run oslo validate -f "${TEST_SUITE_INPUTS}/unknown-field.yaml"
  assert_failure
  assert_output "Error: failed to read objects from ${TEST_SUITE_INPUTS}/unknown-field.yaml: failed to decode openslo/v1 Service: json: unknown field \"markdown\""
}

@test "v1alpha" {
  run oslo validate -f "${TEST_SUITE_INPUTS}/v1alpha"
  assert_failure
  assert_output "Error: issue parsing objects: error unmarshaling JSON: while decoding JSON: failed to decode object: unsupported openslo.Version: foo"
}
