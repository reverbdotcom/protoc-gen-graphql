.PHONY: build all clean

all: clean
clean:
	@rm -rf target/

build: clean
	@env GOOS=linux go build -o target/protoc-gen-graphql.linux
	@env GOOS=darwin go build -o target/protoc-gen-graphql.darwin
