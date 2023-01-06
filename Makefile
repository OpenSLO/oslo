# Copyright © 2022 OpenSLO Team
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

VERSION := $(shell git describe --tags)
.PHONY: build
build:
	go build -ldflags="-X 'main.version=${VERSION}'" -o bin/oslo

.PHONY: install
install: run/tests build
	mkdir -p ${GOPATH}/bin
	mv ./bin/oslo ${GOPATH}/bin/

.PHONY: install/checks/spell-and-markdown
install/checks/spell-and-markdown:
	yarn

.PHONY: install/bats
install/bats:
	@yarn global add bats

.PHONY: install/bats/dependencies
install/bats/dependencies:
	@./test/cli/install-deps.sh

.PHONY: run/checks/spell-and-markdown
run/checks/spell-and-markdown:
	yarn check-trailing-whitespaces
	yarn check-word-lists
	yarn cspell --no-progress '**/**'
	yarn markdownlint '*.md'

.PHONY: run/checks/golangci-lint
run/checks/golangci-lint:
	golangci-lint run

.PHONY: run/tests
run/tests:
	go test -race -cover ./...

.PHONY: run/bats/tests
run/bats/tests:
	find test/cli -type f -name "*.bats" ! -path "test/cli/test_helper/*" | xargs bats
