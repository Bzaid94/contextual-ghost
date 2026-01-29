BINARY_NAME=contextual-ghost

.PHONY: build clean run install

build:
	go build -o $(BINARY_NAME) ./cmd/ghost

run:
	go run ./cmd/ghost

install:
	go install ./cmd/ghost

clean:
	rm -f $(BINARY_NAME)
	go clean
