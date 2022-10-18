build:
	go build

build-osx-intel:
	GOOS=darwin GOARCH=amd64 go build -o bin/app-amd64-darwin .

build-osx-apple:
	GOOS=darwin GOARCH=arm64 go build -o bin/app-arm64-darwin .

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/app-amd64-linux .


tests:
	go test -v ./... # recursively