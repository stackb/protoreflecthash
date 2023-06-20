.PHONY: test
test:
	go test ./...

.PHONY: test_protoset
test_protoset:
	bazel build //test_protos:protoset
	cp -f bazel-bin/test_protos/protoset.pb testdata/

# vendor in the generated files for regular go build process
pb_go:
	bazel build //test_protos/schema/proto3:proto3_go_proto
	cp -f bazel-bin/test_protos/schema/proto3/proto3_go_proto_//github.com/stackb/protoreflecthash/test_protos/schema/proto3/*.pb.go test_protos/generated/latest/proto3
