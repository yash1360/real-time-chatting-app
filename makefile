# Define binary file name
BINARY_NAME=chat

# Default target: when you just type make, it executes build
.PHONY: all
all: build

# Build: compile to generate binary file
.PHONY: build
build:
	go build -o $(BINARY_NAME) .

# Run: compile first, then run
.PHONY: run
run: build
	./$(BINARY_NAME)

# Test: run tests
.PHONY: test
test:
	go test -v ./...

# Clean: delete generated binary files
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

# Format: format code
.PHONY: fmt
fmt:
	go fmt ./...