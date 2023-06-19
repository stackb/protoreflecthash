.PHONY: test
test:
	go test ./...

.PHONY: test_protoset
test_protoset:
	bazel build //test_protos/schema/proto3
	cp -f bazel-bin/test_protos/schema/proto3/proto3-descriptor-set.proto.bin testdata/test_protos.protoset.pb