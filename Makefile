# Copyright (c) 2024 rodneyosodo
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at:
# http://www.apache.org/licenses/LICENSE-2.0

.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run
	@echo "Linting complete."

.PHONY: mocks
mocks:
	@echo "Generating mocks..."
	@go generate ./...
	@echo "Mocks generated."

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v --race -count 1 -coverprofile=cover.out $$(go list ./... | grep -v 'mocks\|internal\|middleware\|proto')
	@echo "Tests complete."

.PHONY: proto
proto:
	@echo "Generating protobuf files..."
	@protoc -I. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative users/proto/users.proto
	@echo "Protobuf files generated."

