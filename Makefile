build:
	mkdir -p bin
	go build -o bin/qr src/*.go

build-all:
	mkdir -p bin

	env GOOS=linux GOARCH=386 go build -o bin/qr-linux32 src/*.go
	env GOOS=linux GOARCH=amd64 go build -o bin/qr-linux64 src/*.go

	env GOOS=windows GOARCH=386 go build -o bin/qr-windows32.exe src/*.go
	env GOOS=windows GOARCH=amd64 go build -o bin/qr-windows64.exe src/*.go

	env GOOS=darwin GOARCH=amd64 go build -o bin/qr-darwin64 src/*.go
