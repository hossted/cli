linux: main.go
	go build -o bin/linux/hossted main.go

windows: main.go
	GOOS=windows GOARCH=386 go build -o bin/windows/hossted.exe main.go

macs: main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/osx/hossted main.go
