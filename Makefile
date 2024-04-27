.PHONY: build test

# I'm assuming you guys are using Mac, so this should hopefully make it easy to run
build: 
	 @env GOOS=windows GOARCH=amd64 go build -o bin/windows/app.exe cmd/main.go
	 @env GOOS=darwin GOARCH=amd64 go build -o bin/mac/app cmd/main.go
	 @env GOOS=linux GOARCH=amd64 go build -o bin/linux/app cmd/main.go

run:
	 @go run cmd/main.go

test:
	 go test -v ./...
