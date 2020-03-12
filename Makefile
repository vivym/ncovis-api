.PHONY: all clean build

all: build

clean:
	rm -f build/ncovis-server

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/ncovis-server -ldflags -s -a -installsuffix cgo cmd/ncovis-server/*.go
