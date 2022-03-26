PACKAGE=github.com/hossted/cli
VERSION=v"$(shell git describe --tags --always --abbrev=0 --match='[0-9]*.[0-9]*.[0-9]*' 2> /dev/null )"
COMMIT_HASH="$(shell git rev-parse --short HEAD)"
BUILD_TIMESTAMP=$(shell date '+%Y-%m-%d')
LDFLAGS="-X '${PACKAGE}/cmd.VERSION=${VERSION}' -X '${PACKAGE}/cmd.COMMITHASH=${COMMIT_HASH}' -X '${PACKAGE}/cmd.BUILDTIME=${BUILD_TIMESTAMP}'"
DEVFLAGS="-X '${PACKAGE}/cmd.VERSION=dev' -X '${PACKAGE}/cmd.COMMITHASH=${COMMIT_HASH}' -X '${PACKAGE}/cmd.BUILDTIME=${BUILD_TIMESTAMP}'"

linux: main.go
	GOOS=linux GOARCH=amd64 go build -o bin/linux/hossted -v -ldflags=${LDFLAGS}

windows: main.go
	GOOS=windows GOARCH=386 go build -o bin/windows/hossted.exe main.go

osx: main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/osx/hossted main.go

dev: main.go
	go build -o bin/dev/hossted -v -ldflags=${DEVFLAGS}

test: main.go
        go test -v ./... -short
