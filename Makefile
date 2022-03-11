PACKAGE=github.com/hossted/cli
VERSION=v"$(shell git describe --tags --always --abbrev=0 --match='[0-9]*.[0-9]*.[0-9]*' 2> /dev/null)"
COMMIT_HASH="$(shell git rev-parse --short HEAD)"
BUILD_TIMESTAMP=$(shell date '+%Y-%m-%d')
LDFLAGS="-X '${PACKAGE}/cmd.VERSION=${VERSION}' -X '${PACKAGE}/cmd.COMMITHASH=${COMMIT_HASH}' -X '${PACKAGE}/cmd.BUILDTIME=${BUILD_TIMESTAMP}'"

linux: main.go
	go build -o bin/linux/hossted main.go

windows: main.go
	GOOS=windows GOARCH=386 go build -o bin/windows/hossted.exe main.go

macs: main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/osx/hossted main.go
