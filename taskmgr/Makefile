all: build test lint migration_up

build:
	go build -v -ldflags="-X github.com/MrShanks/Taska/taskmgr/main.version=$(shell yq '.version' config.yaml)" -o taskmgr .

test:
	go test -v -count=1 ./...

lint:
	golangci-lint run --config=../.golangci.yml

migration_up:
	goose -dir goose/ postgres "postgres://taskauser:$(POSTGRES_PWD)@localhost:5432/taskadb" up
	
migration_down:
	goose -dir goose/ postgres "postgres://taskauser:$(POSTGRES_PWD)@localhost:5432/taskadb" down
	
run:
	go run ./...

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
