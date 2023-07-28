.PHONY: all
all: build

.PHONY: build
build: generate
	go install .

.PHONY: generate
generate:
	go generate ./...

.PHONY: clean
clean:
	go clean
