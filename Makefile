.EXPORT_ALL_VARIABLES:

.PHONY: build all clean test
GOFLAGS=-mod=vendor
GOPROXY="off"

targets = \
	target/protoc-gen-graphql.linux.amd64 \
	target/protoc-gen-graphql.linux.arm64 \
	target/protoc-gen-graphql.darwin.arm64

all: clean
clean:
	@rm -rf target/

build: clean
	@env GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o target/protoc-gen-graphql.linux
	@env GOOS=darwin go build $(GOFLAGS) -o target/protoc-gen-graphql.darwin

$(targets): export TARGET_ARCH=$(patsubst .%,%, $(suffix $(notdir $@)))
$(targets): export TARGET_OS=$(patsubst .%,%, $(suffix $(basename $(notdir $@))))
$(targets):
	@env GOOS=$${TARGET_OS} GOARCH=$${TARGET_ARCH} go build $(GOFLAGS) -o target/protoc-gen-graphql.$${TARGET_OS}.$${TARGET_ARCH}

fixtures/money.pb: fixtures/money.proto
	$(info Generating fixtures...)
	@cd fixtures && go generate

test: fixtures/money.pb
	@go test -v -race -cover ./
