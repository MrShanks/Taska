MAKEFILE_DIR := $(dir $(MAKEFILE_LIST))

all: build test lint migration_up

build:
	go build -v -ldflags="-X github.com/MrShanks/Taska/taskmgr/main.version=$(shell yq '.version' $(MAKEFILE_DIR)config.yaml)" -o taskmgr $(realpath $(MAKEFILE_DIR))

test:
	go test -v -count=1 ./...

lint:
	golangci-lint run --config=../.golangci.yml

migration_up:
	goose -dir "$(MAKEFILE_DIR)postgresdb/goose/" postgres "postgres://taskauser:$(POSTGRES_PWD)@localhost:5432/taskadb" up
	
run:
	go run ./...

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
