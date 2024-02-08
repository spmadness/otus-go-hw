include .env

BIN_CALENDAR := "./bin/calendar"
BIN_SCHEDULER := "./bin/scheduler"
BIN_SENDER := "./bin/sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate:
	rm -rf internal/server/grpc/pb
	mkdir -p internal/server/grpc/pb

	protoc \
		--proto_path=api \
		--go_out=internal/server/grpc/pb \
		--go-grpc_out=internal/server/grpc/pb \
		api/*.proto

build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(BIN_CALENDAR) -config ./configs/config_calendar.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN_CALENDAR) version
	$(BIN_SCHEDULER) version
	$(BIN_SENDER) version

test:
	go test -race -count=1 ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

lint: install-lint-deps
	golangci-lint run ./...

migrate:
	goose -dir ./migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable" up

.PHONY: generate build run build-img run-img version test lint migrate
