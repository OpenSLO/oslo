#!/usr/bin/env bats

setup_suite() {
    # build the bin to test
    make build
}

teardown_suite() {
    # cleanup
    rm -rf "${OSLO_BIN}"
}
