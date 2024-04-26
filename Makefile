.PHONY: build test

build: 
	 @go build -o bin/app.exe cmd/main.go

run:
	 @go run cmd/main.go

test:
	 go test -v ./...
