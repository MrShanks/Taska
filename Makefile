all: build test run

build:
	go build -v -ldflags="-X github.com/MrShanks/Taska/cmd.version=$(shell yq '.app.version' config.yaml)" -o taskcli .

test:
	go test -v -count=1 ./...

run:
	go run ./...
