.PHONY: build all clean test

all: clean
clean:
	@rm -rf target/

build: clean
	@env GOOS=linux go build -o target/protoc-gen-graphql.linux
	@env GOOS=darwin go build -o target/protoc-gen-graphql.darwin

fixtures/money.pb: fixtures/money.proto
	$(info Generating fixtures...)
	@cd fixtures && go generate

test: fixtures/money.pb
	@go test -v -race -cover ./
