# Copyright (c) 2024 rodneyosodo
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at:
# http://www.apache.org/licenses/LICENSE-2.0

DOCKER_IMAGE_NAME_PREFIX=ghcr.io/rodneyosodo/twiga
SERVICES = users posts notifications
DOCKERS = $(addprefix docker_,$(SERVICES))
DOCKERS_DEV = $(addprefix docker_dev_,$(SERVICES))
BUILD_DIR = build
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOARM ?= $(shell go env GOARM)
CGO_ENABLED ?= 0
VERSION ?= $(shell git describe --abbrev=0 --tags 2>/dev/null || echo "v0.0.0")
TIME ?= $(shell date +%F_%T)
COMMIT ?= $(shell git rev-parse --short HEAD)
COMMIT_DATE ?= $(shell git log -1 --date=format:"%F_%T" --format=%cd)

define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) \
	go build -ldflags "-s -w \
	-X 'go.szostok.io/version.version=$(VERSION)' \
	-X 'go.szostok.io/version.commit=$(COMMIT)' \
	-X 'go.szostok.io/version.commitDate=$(COMMIT_DATE)' \
	-X 'go.szostok.io/version.buildDate=$(TIME)'" \
	-o ${BUILD_DIR}/$(1) cmd/$(1)/main.go
endef

define make_docker
	$(eval svc=$(subst docker_,,$(1)))

	docker buildx build \
		--no-cache \
		--build-arg SVC=$(svc) \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg GOARM=$(GOARM) \
		--build-arg VERSION=$(VERSION) \
		--build-arg TIME=$(TIME) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg COMMIT_DATE=$(COMMIT_DATE) \
		--tag=$(DOCKER_IMAGE_NAME_PREFIX)/$(svc) \
		-f docker/Dockerfile .
endef

define make_docker_dev
	$(eval svc=$(subst docker_dev_,,$(1)))

	docker buildx build \
		--no-cache \
		--build-arg SVC=$(svc) \
		--tag=$(DOCKER_IMAGE_NAME_PREFIX)/$(svc) \
		-f docker/Dockerfile.dev ./build
endef

.PHONY: all
all: $(SERVICES)

.PHONY: $(SERVICES)
$(SERVICES):
	$(call compile_service,$(@))

.PHONY: $(DOCKERS)
$(DOCKERS):
	$(call make_docker,$(@))

.PHONY: $(DOCKERS_DEV)
$(DOCKERS_DEV):
	$(call make_docker_dev,$(@))

.PHONY: dockers
dockers: $(DOCKERS)

.PHONY: dockers_dev
dockers_dev: $(DOCKERS_DEV)

.PHONY: run
run:
	docker compose -f docker/docker-compose.yaml --env-file docker/.env up

.PHONY: stop
stop:
	docker compose -f docker/docker-compose.yaml down

.PHONY: down
down:
	docker compose -f docker/docker-compose.yaml down --volumes

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

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
	@go test -v --race -count 1 -coverprofile=cover.out $$(go list ./... | grep -v 'mocks\|internal\|middleware\|proto\|cmd\|consumer\|producer')
	@echo "Tests complete."

.PHONY: proto
proto:
	@echo "Generating protobuf files..."
	@protoc -I. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative users/proto/users.proto
	@echo "Protobuf files generated."

