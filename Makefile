.EXPORT_ALL_VARIABLES:

.PHONY: build all clean test
GOFLAGS=-mod=vendor
GOPROXY="off"

all: build

clean:
	@rm -rf target/

build: clean
	@env GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o target/protoc-gen-graphql.linux.amd64
	@env GOOS=linux GOARCH=arm64 go build $(GOFLAGS) -o target/protoc-gen-graphql.linux.arm64
	@env GOOS=darwin go build $(GOFLAGS) -o target/protoc-gen-graphql.darwin.arm64

fixtures/money.pb: fixtures/money.proto
	$(info Generating fixtures...)
	@cd fixtures && go generate

test: fixtures/money.pb
	@go test -v -race -cover ./
