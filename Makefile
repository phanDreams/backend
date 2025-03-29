include .env

LOCAL_BIN:=$(CURDIR)/bin

build:
	go build -o bin/main ./cmd/.

run:
	make build
	./bin/main

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

install-goose:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

local-migration-status:
	${LOCAL_BIN}/goose -v postgres status -dir ${MIGRATION_DIR}

local-migration-up:
	${LOCAL_BIN}/goose -v postgres up -dir ${MIGRATION_DIR}

local-migration-down:
	${LOCAL_BIN}/goose -v postgres down -dir ${MIGRATION_DIR}