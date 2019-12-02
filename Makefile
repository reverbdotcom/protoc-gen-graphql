.PHONY: build all clean test

.EXPORT_ALL_VARIABLES:

GOFLAGS=-mod=vendor
GOPROXY="off"

all: clean
clean:
	@rm -rf target/

build: clean
	@GOOS=linux GOPROXY=off go build $(GOFLAGS) -o target/protoc-gen-graphql.linux
	@GOOS=darwin GOPROXY=off go build $(GOFLAGS) -o target/protoc-gen-graphql.darwin

fixtures/money.pb: fixtures/money.proto
	$(info Generating fixtures...)
	@cd fixtures && go generate

test: fixtures/money.pb
	@go test -v -race -cover ./
