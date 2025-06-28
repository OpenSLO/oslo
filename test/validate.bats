#!/usr/bin/env bats
# bats file_tags=unit

setup() {
  load "test_helper/load"
  load_lib "bats-assert"
  load_lib "bats-support"
}

@test "validate an invalid api version" {
  run oslo validate -f "${TEST_SUITE_INPUTS}/validate/invalid-apiversion.yaml"
  assert_failure
  assert_output "Error: failed to read objects from ${TEST_SUITE_INPUTS}/validate/invalid-apiversion.yaml: error unmarshaling JSON: while decoding JSON: failed to decode object: unsupported openslo.Version: foo"
}

@test "validate unknown field" {
  run oslo validate -f "${TEST_SUITE_INPUTS}/validate/unknown-field.yaml"
  assert_failure
  assert_output "Error: failed to read objects from ${TEST_SUITE_INPUTS}/validate/unknown-field.yaml: failed to decode openslo/v1 Service: json: unknown field \"markdown\""
}

@test "v1alpha" {
  run oslo validate -f "${TEST_SUITE_INPUTS}/validate/v1alpha.yaml"
  assert_failure
  assert_output "$(cat "${TEST_SUITE_OUTPUTS}/validate/v1alpha")"
}

@test "v1" {
  run oslo validate -f "${TEST_SUITE_INPUTS}/validate/v1.yaml"
  assert_failure
  assert_output "$(cat "${TEST_SUITE_OUTPUTS}/validate/v1")"
}

@test "v2alpha" {
  run oslo validate -f "${TEST_SUITE_INPUTS}/validate/v2alpha.yaml"
  assert_failure
  assert_output "$(cat "${TEST_SUITE_OUTPUTS}/validate/v2alpha")"
}

@test "mix of files" {
  run oslo validate -f "${TEST_SUITE_INPUTS}/validate/mix"
  assert_failure
  assert_output "$(cat "${TEST_SUITE_OUTPUTS}/validate/mix")"
}

@test "recursive directory read" {
  run oslo validate -R -f "${TEST_SUITE_INPUTS}/validate/recursive"
  assert_failure
  assert_output "$(cat "${TEST_SUITE_OUTPUTS}/validate/recursive")"
}
