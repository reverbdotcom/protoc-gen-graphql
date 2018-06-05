protoc-gen-graphql
---

Produces a graphql schema from your protos.

Usage:

`make` compiles Linux and MacOS binaries to `./target/protoc-gen-graphql.linux` and `./target/protoc-gen-graphql.darwin`.

To generate a `schema.graphql` file with `protoc`, copy or symlink the platform-specific binary to `protoc-gen-graphql`, then run:

```
protoc \
-I your_protos_dir \
--plugin ./protoc-gen-graphql \
--graphql_out=file_out=schema.graphql,root_query=GraphQL:your_schemas_dir \
your_protos_dir/message.proto
```
