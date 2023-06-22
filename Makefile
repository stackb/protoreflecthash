define BAZELWORKSPACE
workspace(name = "com_github_stackb_protoreflecthash")
endef

define BAZELBUILD
load("@rules_go//go:def.bzl", "go_library", "go_test")
load("@gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/stackb/protoreflecthash
# gazelle:resolve go go github.com/stackb/protoreflecthash/test_protos/generated/latest/proto2 @com_github_stackb_protoreflecthash//test_protos/generated/latest/proto2
# gazelle:resolve go go github.com/stackb/protoreflecthash/test_protos/generated/latest/proto3 @com_github_stackb_protoreflecthash//test_protos/generated/latest/proto3

gazelle(name = "gazelle")
endef

define BAZELVERSION
6.2.1
endef

define BAZELMODULE
module(
    name = "protoreflecthash",
    version = "0.1.0",
    repo_name = "com_github_stackb_protoreflecthash",
)

bazel_dep(name = "protobuf", version = "21.7", repo_name = "com_google_protobuf")
bazel_dep(name = "rules_proto", version = "4.0.0")
bazel_dep(name = "rules_go", version = "0.39.1")
bazel_dep(name = "gazelle", version = "0.31.1")
endef

define BAZELRC
build --experimental_enable_bzlmod
endef

.PHONY: test
test:
	go test github.com/stackb/protoreflecthash

.PHONY: test_protoset
test_protoset:
	bazel build //test_protos:protoset
	cp -f bazel-bin/test_protos/protoset.pb testdata/

.PHONY: bazel
bazel:
	echo "$$BAZELWORKSPACE" > WORKSPACE
	echo "$$BAZELMODULE" > MODULE.bazel
	echo "$$BAZELBUILD" > BUILD.bazel
	echo "$$BAZELVERSION" > .bazelversion
	echo "$$BAZELRC" > .bazelrc
	# bazel build //test_protos/schema/proto3:proto3_go_proto //test_protos/schema/proto2:proto2_go_proto
	# cp -f bazel-bin/test_protos/schema/proto3/proto3_go_proto_/github.com/stackb/protoreflecthash/test_protos/generated/latest/proto3/*.pb.go test_protos/generated/latest/proto3
	# cp -f bazel-bin/test_protos/schema/proto2/proto2_go_proto_/github.com/stackb/protoreflecthash/test_protos/generated/latest/proto2/*.pb.go test_protos/generated/latest/proto2
