.PHONY: proto
proto: $(addprefix internal/sbproto/,$(basename $(notdir $(wildcard sb/protolib/*.proto))))

internal/sbproto/%: sb/protolib/%.proto
	mkdir -p internal/sbproto/$(basename $(notdir $<))
	(protoc >/dev/null) || (echo "Please install protoc https://grpc.io/docs/protoc-installation/" && exit 1)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
	PATH="$(PATH):$(shell go env GOPATH)/bin" protoc -I=sb --go_out=$@ --go_opt=module=go.snowblossom/$@ --go-grpc_out=$@ --go-grpc_opt=module=go.snowblossom/$@ $<

