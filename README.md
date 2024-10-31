# granted-rpc-go

## Local development

To test the code generator, run:

```bash
go install ./cmd/protoc-gen-granted-rpc-go &&
          protoc --go_out=. --go_opt=paths=source_relative \
              --granted-rpc-go_out=. --granted-rpc-go_opt=paths=source_relative \
              example/example.proto
```
