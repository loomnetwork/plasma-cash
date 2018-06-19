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
		github.com/spf13/viper \
		github.com/loomnetwork/go-loom \
		github.com/pkg/errors \
		github.com/go-kit/kit/log \
		github.com/ethereum/go-ethereum \
		github.com/gorilla/mux \
		github.com/inconshreveable/mousetrap \
		gopkg.in/check.v1

clean:
	go clean
	rm -f \
		src/types/types.pb.go \
		testdata/test.pb.go \