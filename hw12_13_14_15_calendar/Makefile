include deployments/.env

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

build-calendar:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-scheduler:
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler

build-sender:
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

build: build-calendar build-scheduler build-sender

run-calendar: build-calendar
	$(BIN_CALENDAR) -config ./configs/config_calendar.toml

run-scheduler: build-scheduler
	$(BIN_SCHEDULER) -config ./configs/config_scheduler.toml

run-sender: build-sender
	$(BIN_SENDER) -config ./configs/config_sender.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

run-docker-postgres:
	docker run -d --name $(POSTGRES_CONTAINER) \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-e POSTGRES_DB=$(POSTGRES_DB) \
	-p $(POSTGRES_PORT):5432 \
	-v postgres-data:/var/lib/postgresql/data \
	-v $(shell pwd)/deployments/init.sql:/docker-entrypoint-initdb.d/init.sql \
	postgres:latest

run-docker-rabbitmq:
	docker run -d --name $(RABBITMQ_CONTAINER) \
	-p $(RABBITMQ_PORT):5672 \
	-p $(RABBITMQ_PORT_UI):15672 \
	-e RABBITMQ_DEFAULT_USER=$(RABBITMQ_USER) \
	-e RABBITMQ_DEFAULT_PASS=$(RABBITMQ_PASSWORD) \
	rabbitmq:management

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
