.PHONY: build run test race fmt vet clean

BINARY_NAME=echallan-calculator

build:
	go build -o $(BINARY_NAME) cmd/echallan-calculator/main.go

run:
	go run cmd/echallan-calculator/main.go

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
