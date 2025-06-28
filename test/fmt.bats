#!/usr/bin/env bats
# bats file_tags=unit

setup() {
	load "test_helper/load"
	load_lib "bats-assert"
	load_lib "bats-support"
}

@test "oslo formats a single service" {
  run oslo fmt -f "${TEST_SUITE_INPUTS}/fmt/service.yaml"
  assert_success
  assert_output "$(cat "${TEST_SUITE_OUTPUTS}/fmt/service.yaml")"
}
