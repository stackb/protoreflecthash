package protoreflecthash

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func TestString(t *testing.T) {
	files := unmarshalProtoRegistryFiles(t, testProtoset)

	for name, tc := range map[string]struct {
		md   protoreflect.MessageDescriptor
		json string
		want string
	}{
		"Int32MessageZero": {
			md:   mdByPath(t, files, "test_protos/schema/proto3/integers.proto", "Int32Message"),
			json: `{"values": [0, 1, 2]}`,
			want: "ec28f92dbcce2dc9e38b48cd7725337ca7df40d729b8523a5b3512f7449e8156",
		},
	} {
		t.Run(name, func(t *testing.T) {
			msg := dynamicpb.NewMessage(tc.md)
			if err := protojson.Unmarshal([]byte(tc.json), msg); err != nil {
				t.Fatal(err)
			}

			got, err := String(msg)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("protohash (-want +got):\n%s", diff)
			}
		})
	}
}
