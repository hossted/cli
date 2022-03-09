linux: main.go
	go build -o bin/linux/hossted main.go
	chomod 744 bin/linux/hossted

windows: main.go
	GOOS=windows GOARCH=386 go build -o bin/windows/hossted.exe main.go
	chomod 744 bin/linux/hossted

macs: main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/osx/hossted main.go
	chomod 744 bin/linux/hossted
