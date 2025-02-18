#!/usr/bin/env bats

setup_suite() {
  export TEST_SUITE_INPUTS="$BATS_TEST_DIRNAME/inputs"
  export TEST_SUITE_OUTPUTS="$BATS_TEST_DIRNAME/outputs"
}
