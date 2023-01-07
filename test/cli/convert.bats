#!/usr/bin/env bats

setup_file() {
    mkdir -p bin
    # get the latest nobl9 cli tool for testing
    curl https://api.github.com/repos/nobl9/sloctl/releases/latest | \
      grep "browser_download_url.*linux*" | \
      cut -d : -f 2,3 | tr -d \" | wget -O sloctl -qi -
    chmod +x sloctl
    mv sloctl bin/
}

setup() {
    load 'test_helper/common-setup'
    _common_setup
}

teardown_file() {
    rm ./bin/sloctl
}

@test "oslo errors without output flag set" {
  run sloctl convert
  assert_equal $status 1
}

#--------------------
#
# Nobl9 Coverting
#
@test "nobl9 - oslo fails when file doesn't exist" {
  run oslo convert -f test/foo.yaml -o nobl9
  assert_equal $status 1
  assert_output --partial "no such file or directory"
}

@test "nobl9 - oslo converts successfully" {
  run oslo convert -f test/v1/service/service.yaml -o nobl9
  assert_equal $status 0
  assert_output
}

@test "nobl9 - oslo converts file with multiple kinds successfully" {
  run oslo convert -f test/v1/multi.yaml -o nobl9
  assert_equal $status 0
  assert_output --partial "apiVersion"
}

@test "nobl9 - oslo converts multiple files successfully" {
  run oslo convert \
    -f test/v1/slo/slo-indicatorRef-rolling-cloudwatch.yaml \
    -f test/v1/sli/sli-threshold-cloudwatch.yaml \
    -f test/v1/data-source/data-source-cloudwatch.yaml \
    -o nobl9
  assert_equal $status 0
  assert_output --partial "apiVersion"
}

@test "nobl9 - sloctl parses converted openslo files successfully" {
  skip "Skipping until sloctl can do offline validation"
}
