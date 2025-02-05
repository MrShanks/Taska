all: build test run

build:
	go build -v -ldflags="-X github.com/MrShanks/Taska/taskmgr/main.version=$(shell yq '.app.version' config.yaml)" -o taskmgr .

test:
	go test -v -count=1 ./...

run:
	go run ./...
