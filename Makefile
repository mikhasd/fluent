GOROOT ?= /usr/local/go

build:
	${GOROOT}/bin/go build -v ./...

test:
	${GOROOT}/bin/go test -race -v -coverprofile=cover.out -coverpkg=./... ./...