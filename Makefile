.PHONY: build test run

build:
	go vet ./...
	go build -o anime-server

test:
	go test ./...

run: build
	./anime-server
