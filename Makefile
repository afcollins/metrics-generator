BINARY := metrics-generator

.PHONY: build test clean

build:
	go build -o $(BINARY) ./cmd/$(BINARY)

test:
	go test ./...

clean:
	rm -f $(BINARY)
