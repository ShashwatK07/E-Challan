.PHONY: build run test fmt lint vet race clean

BINARY_NAME=echallan-services

build:
	go build -o $(BINARY_NAME) cmd/echallan-services/main.go

run:
	go run cmd/echallan-services/main.go

test:
	go test -v ./...

race:
	go test -race -v ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY_NAME)
