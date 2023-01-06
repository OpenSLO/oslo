#!/bin/bash

LIBS_BATS_SUPPORT_VERSION="0.3.0"
LIBS_BATS_ASSERT_VERSION="2.1.0"

echo -n "Installing bats dependencies"
# change our working dir to the test folder that this lives in
cd "$(dirname $0)" || exit
pwd
mkdir -p ./test_helper/bats-support
echo -n "."
curl -sSL https://github.com/ztombol/bats-support/archive/refs/tags/v${LIBS_BATS_SUPPORT_VERSION}.tar.gz -o /tmp/bats-support.tgz
echo -n "."
tar -zxf /tmp/bats-support.tgz -C ./test_helper/bats-support --strip 1
echo -n "."
rm -rf /tmp/bats-support.tgz
echo -n "."

mkdir -p ./test_helper/bats-assert
echo -n "."
curl -sSL https://github.com/bats-core/bats-assert/archive/refs/tags/v${LIBS_BATS_ASSERT_VERSION}.tar.gz -o /tmp/bats-assert.tgz
echo -n "."
tar -zxf /tmp/bats-assert.tgz -C ./test_helper/bats-assert --strip 1
echo -n "."
rm -rf /tmp/bats-assert.tgz
echo -n "."

echo "Done!"
