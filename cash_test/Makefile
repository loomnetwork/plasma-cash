CURRENT_DIRECTORY = $(shell pwd)



.PHONY: all clean test lint deps proto

cli:
	go build -o plasmascash_tester src/cmd/main.go

test: 
	go test ./...

lint:
	golint ./...

deps:
	go get \
		github.com/gogo/protobuf/jsonpb \
		github.com/gogo/protobuf/proto \
		github.com/spf13/cobra \
		github.com/gomodule/redigo/redis \
		github.com/loomnetwork/go-loom \
		github.com/hashicorp/go-plugin \
		github.com/pkg/errors \
		github.com/grpc-ecosystem/go-grpc-prometheus \
		github.com/go-kit/kit/log

clean:
	go clean
	rm -f \
		src/types/types.pb.go \
		testdata/test.pb.go \