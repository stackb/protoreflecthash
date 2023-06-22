#!/bin/bash

write_workspace() {
    cat << EOF > WORKSPACE
workspace(name = "com_github_stackb_protoreflecthash")
EOF
}

write_modulebazel() {
    cat << EOF > MODULE.bazel
module(
    name = "protoreflecthash",
    version = "0.1.0",
    repo_name = "com_github_stackb_protoreflecthash",
)

bazel_dep(name = "protobuf", version = "21.7", repo_name = "com_google_protobuf")
bazel_dep(name = "rules_proto", version = "4.0.0")
bazel_dep(name = "rules_go", version = "0.39.1")
bazel_dep(name = "gazelle", version = "0.31.1")
EOF
}

write_bazelversion() {
    cat << EOF > .bazelversion
6.2.1
EOF
}

write_bazelrc() {
    cat << EOF > .bazelrc
build --experimental_enable_bzlmod
EOF
}

write_bazelbuild() {
    cat << EOF > BUILD.bazel
load("@rules_go//go:def.bzl", "go_library", "go_test")
load("@gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/stackb/protoreflecthash
# gazelle:resolve go go github.com/stackb/protoreflecthash/test_protos/generated/latest/proto2 @com_github_stackb_protoreflecthash//test_protos/generated/latest/proto2
# gazelle:resolve go go github.com/stackb/protoreflecthash/test_protos/generated/latest/proto3 @com_github_stackb_protoreflecthash//test_protos/generated/latest/proto3

gazelle(name = "gazelle")
endef

define BAZELVERSION
endef
EOF
}

run_gazelle() {
    bazel run gazelle
}

build_targets() {
	bazel build \
		//test_protos/schema/proto3:proto3_go_proto \
		//test_protos/schema/proto2:proto2_go_proto \
		//test_protos:protoset
}

copy_pb_go() {
	cp -f bazel-bin/test_protos/schema/proto3/proto3_go_proto_/github.com/stackb/protoreflecthash/test_protos/generated/latest/proto3/*.pb.go \
		test_protos/generated/latest/proto3

	cp -f bazel-bin/test_protos/schema/proto2/proto2_go_proto_/github.com/stackb/protoreflecthash/test_protos/generated/latest/proto2/*.pb.go \
		test_protos/generated/latest/proto2
}

copy_protoset() {
	cp -f bazel-bin/test_protos/protoset.pb \
		testdata/	
}

clean() {
	bazel clean
	rm .bazelversion .bazelrc MODULE.bazel WORKSPACE
	find . -name 'BUILD.bazel' | xargs rm
}

main() {
    write_workspace
    write_modulebazel
    write_bazelversion
    write_bazelrc
    write_bazelbuild
    run_gazelle
    build_targets
    copy_pb_go
    copy_protoset
}

clean