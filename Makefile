RELEASE ?= 1.1.0

default: build

build:
	go build -o bin/calculate cmd/calculate.go

install:
	cp bin/calculate /usr/local/bin

clean:
	rm -r ./bin/*

buildall:
	GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/calculate cmd/calculate.go
	GOOS=linux GOARCH=386 go build -o bin/linux/i386/calculate cmd/calculate.go
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/amd64/calculate cmd/calculate.go
	GOOS=windows GOARCH=amd64 go build -o bin/windows/amd64/calculate.exe cmd/calculate.go
	GOOS=windows GOARCH=386 go build -o bin/windows/i386/calculate.exe cmd/calculate.go
	GOOS=solaris GOARCH=amd64 go build -o bin/solaris/amd64/calculate cmd/calculate.go

release: buildall
	tar czvf shim-calculate_v$(RELEASE).tar.gz bin/
