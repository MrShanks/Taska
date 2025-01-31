all: build test run

build:
	go build -v -ldflags="-X github.com/MrShanks/Taska/cmd.version=$(shell yq '.app.version' config.yaml)" -o taskcli .

test:
<<<<<<< HEAD
	go test -v -count=1 ./...
=======
	go test -v ./...
>>>>>>> 7717e80 (feat: TAS-10 add faulty test)

run:
	go run ./...
