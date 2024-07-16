PACKAGE=github.com/hossted/cli
SERVICE_COMMON_PACKAGE=${PACKAGE}/hossted/service/common
VERSION=v"$(shell git describe --tags --always --abbrev=0 --match='[0-9]*.[0-9]*.[0-9]*' 2> /dev/null )"
COMMIT_HASH="$(shell git rev-parse --short HEAD)"
BUILD_TIMESTAMP=$(shell date '+%Y-%m-%d')
LDFLAGS="-X '${PACKAGE}/cmd.VERSION=${VERSION}' \
         -X '${PACKAGE}/cmd.ENVIRONMENT=prod' \
         -X '${PACKAGE}/cmd.COMMITHASH=${COMMIT_HASH}' \
         -X '${PACKAGE}/cmd.BUILDTIME=${BUILD_TIMESTAMP}' \
         -X '${SERVICE_COMMON_PACKAGE}.LOKI_PASSWORD=${LOKI_PASSWORD}' \
         -X '${SERVICE_COMMON_PACKAGE}.LOKI_URL=${LOKI_URL}' \
         -X '${SERVICE_COMMON_PACKAGE}.LOKI_USERNAME=${LOKI_USERNAME}' \
         -X '${SERVICE_COMMON_PACKAGE}.MIMIR_PASSWORD=${MIMIR_PASSWORD}' \
         -X '${SERVICE_COMMON_PACKAGE}.MIMIR_URL=${MIMIR_URL}' \
         -X '${SERVICE_COMMON_PACKAGE}.MIMIR_USERNAME=${MIMIR_USERNAME}' \
         -X '${SERVICE_COMMON_PACKAGE}.HOSSTED_API_URL=${HOSSTED_API_URL}' \
         -X '${SERVICE_COMMON_PACKAGE}.HOSSTED_AUTH_TOKEN=${HOSSTED_AUTH_TOKEN}'"

DEVFLAGS="-X '${PACKAGE}/cmd.VERSION=dev' \
         -X '${PACKAGE}/cmd.ENVIRONMENT=dev' \
         -X '${PACKAGE}/cmd.COMMITHASH=${COMMIT_HASH}' \
         -X '${PACKAGE}/cmd.BUILDTIME=${BUILD_TIMESTAMP}' \
         -X '${SERVICE_COMMON_PACKAGE}.LOKI_PASSWORD=${LOKI_PASSWORD}' \
         -X '${SERVICE_COMMON_PACKAGE}.LOKI_URL=${LOKI_URL}' \
         -X '${SERVICE_COMMON_PACKAGE}.LOKI_USERNAME=${LOKI_USERNAME}' \
         -X '${SERVICE_COMMON_PACKAGE}.MIMIR_PASSWORD=${MIMIR_PASSWORD}' \
         -X '${SERVICE_COMMON_PACKAGE}.MIMIR_URL=${MIMIR_URL}' \
         -X '${SERVICE_COMMON_PACKAGE}.MIMIR_USERNAME=${MIMIR_USERNAME}' \
         -X '${SERVICE_COMMON_PACKAGE}.HOSSTED_API_URL=${HOSSTED_API_URL}' \
         -X '${SERVICE_COMMON_PACKAGE}.HOSSTED_AUTH_TOKEN=${HOSSTED_AUTH_TOKEN}'"

linux: main.go
	GOOS=linux GOARCH=amd64 go build -o bin/linux/hossted-linux-amd64 -v -ldflags=${LDFLAGS}

windows: main.go
	GOOS=windows GOARCH=386 go build -o bin/windows/hossted.exe -v -ldflags=${LDFLAGS} 

osx: main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/osx/hossted-darwin-amd64  -v -ldflags=${LDFLAGS}

dev: main.go
	go build -o bin/dev/hossted -v -ldflags=${DEVFLAGS}

test: main.go
	go test -v ./... -short
