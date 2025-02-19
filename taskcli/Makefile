all: build test run lint

build:
	go build -v -ldflags="-X github.com/MrShanks/Taska/taskcli/cmd.version=$(shell yq '.version' config.yaml)" -o taskcli .

test:
	go test -v -count=1 ./...

run:
	go run ./...

lint:
	golangci-lint run --config=../.golangci.yml

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out