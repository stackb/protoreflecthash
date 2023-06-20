package protoreflecthash

import (
	_ "embed"
	"fmt"
	"log"
	"math"
	"testing"

	"github.com/benlaurie/objecthash/go/objecthash"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/pcj/objecthash-proto/test_protos/generated/latest/proto3"
)

//go:embed testdata/test_protos.protoset.pb
var testProtoset []byte

func TestHashNil(t *testing.T) {
	for name, tc := range map[string]struct {
		value interface{}
		want  string
	}{
		"nil": {
			want: "1b16b1df538ba12dc3f97edbb85caa7050d46c148134290feba80f8236c83db9",
		},
	} {
		t.Run(name, func(t *testing.T) {
			h := hasher{}

			got := getHash(t, func() ([]byte, error) {
				return h.hashNil()
			})

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("protohash (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.want, objectHash(t, tc.value)); diff != "" {
				t.Errorf("objecthash (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHashBool(t *testing.T) {
	for name, tc := range map[string]struct {
		value bool
		want  string
	}{
		"false": {
			value: false,
			want:  "c02c0b965e023abee808f2b548d8d5193a8b5229be6f3121a6f16e2d41a449b3",
		},
		"true": {
			value: true,
			want:  "7dc96f776c8423e57a2785489a3f9c43fb6e756876d6ad9a9cac4aa4e72ec193",
		},
	} {
		t.Run(name, func(t *testing.T) {
			h := hasher{}

			got := getHash(t, func() ([]byte, error) {
				return h.hashBool(tc.value)
			})

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("protohash (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.want, objectHash(t, tc.value)); diff != "" {
				t.Errorf("objecthash (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHashInt(t *testing.T) {
	for name, tc := range map[string]struct {
		value int64
		want  string
	}{
		"zero": {
			value: 0,
			want:  "a4e167a76a05add8a8654c169b07b0447a916035aef602df103e8ae0fe2ff390",
		},
		"positive": {
			value: 1,
			want:  "4cd9b7672d7fbee8fb51fb1e049f690342035f543a8efe734b7b5ffb0c154a45",
		},
		"min": {
			value: math.MinInt,
			want:  "2df43a3eaece5bb912a43ce16ebdf392e1dd9ce14c16255783ca1f5456d7d04f",
		},
		"max": {
			value: math.MaxInt,
			want:  "eda7a99bc44462f5181f63a88e2ab9d8d318d68c2c2bf9ff70d9e4ecc2a99df3",
		},
	} {
		t.Run(name, func(t *testing.T) {
			h := hasher{}

			got := getHash(t, func() ([]byte, error) {
				return h.hashInt(tc.value)
			})

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("protohash (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.want, objectHash(t, tc.value)); diff != "" {
				t.Errorf("objecthash (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHashFloat(t *testing.T) {
	for name, tc := range map[string]struct {
		value float64
		want  string
	}{
		"zero": {
			value: 0,
			want:  "60101d8c9cb988411468e38909571f357daa67bff5a7b0a3f9ae295cd4aba33d",
		},
		"neg": {
			value: -1.0,
			want:  "f706daa44d7e40e21ea202c36119057924bb28a49949d8ddaa9c8c3c9367e602",
		},
		"pos": {
			value: +1.0,
			want:  "f01adc732390ab024d64080e0b173f0ee3a1610efbdd4ce2a13bbf8d9b26c639",
		},
	} {
		t.Run(name, func(t *testing.T) {
			h := hasher{}

			got := getHash(t, func() ([]byte, error) {
				return h.hashFloat(tc.value)
			})

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("protohash (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.want, objectHash(t, tc.value)); diff != "" {
				t.Errorf("objecthash (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHashString(t *testing.T) {
	for name, tc := range map[string]struct {
		value string
		want  string
	}{
		"zero": {
			value: "",
			want:  "0bfe935e70c321c7ca3afc75ce0d0ca2f98b5422e008bb31c00c6d7f1f1c0ad6",
		},
		"ascii": {
			value: "bob",
			want:  "5ef421eb52293e5e3919d3c6f08413b873129dd859f4d0ff8273e13a494b9e9e",
		},
		"unicode": {
			value: "你好",
			want:  "462b68f5e3d75aed5f02841b4ffee068d4cf33ce1b155105b71a9e5f358026df",
		},
	} {
		t.Run(name, func(t *testing.T) {
			h := hasher{}

			got := getHash(t, func() ([]byte, error) {
				return h.hashString(tc.value)
			})

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("protohash (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.want, objectHash(t, tc.value)); diff != "" {
				t.Errorf("objecthash (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHashBytes(t *testing.T) {
	for name, tc := range map[string]struct {
		value []byte
		want  string
	}{
		"zero": {
			value: []byte{},
			want:  "454349e422f05297191ead13e21d3db520e5abef52055e4964b82fb213f593a1",
		},
		"ascii": {
			value: []byte("bob"),
			want:  "aa75ac53926e8b0711ee730b4c0d8b232b167180f843da40d6e75871cd0785a5",
		},
		"unicode": {
			value: []byte("你好"),
			want:  "39fafdc74a5ee3ff86bd0b982304e58f4685767e87f5176307df9c9e1cf50925",
		},
	} {
		t.Run(name, func(t *testing.T) {
			h := hasher{}

			got := getHash(t, func() ([]byte, error) {
				return h.hashBytes(tc.value)
			})

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("protohash (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.want, objectHash(t, tc.value)); diff != "" {
				t.Errorf("objecthash (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHashList(t *testing.T) {
	for name, tc := range map[string]struct {
		kind  protoreflect.Kind
		value protoreflect.List
		want  string
	}{
		"zero": {
			kind:  protoreflect.StringKind,
			value: stringList{},
			want:  "acac86c0e609ca906f632b0e2dacccb2b77d22b0621f20ebece1a4835b93f6f0",
		},
		"foobar": {
			kind:  protoreflect.StringKind,
			value: stringList{"foo", "bar"},
			want:  "32ae896c413cfdc79eec68be9139c86ded8b279238467c216cf2bec4d5f1e4a2",
		},
	} {
		t.Run(name, func(t *testing.T) {
			h := hasher{}

			got := getHash(t, func() ([]byte, error) {
				return h.hashList(tc.kind, tc.value)
			})

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("protohash (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.want, objectHash(t, tc.value)); diff != "" {
				t.Errorf("objecthash (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHashMap(t *testing.T) {
	intMap := proto3.IntMaps{}

	for name, tc := range map[string]struct {
		kind  protoreflect.Kind
		value protoreflect.List
		want  string
	}{
		"zero": {
			kind:  protoreflect.StringKind,
			value: stringList{},
			want:  "acac86c0e609ca906f632b0e2dacccb2b77d22b0621f20ebece1a4835b93f6f0",
		},
		"foobar": {
			kind:  protoreflect.StringKind,
			value: stringList{"foo", "bar"},
			want:  "32ae896c413cfdc79eec68be9139c86ded8b279238467c216cf2bec4d5f1e4a2",
		},
	} {
		t.Run(name, func(t *testing.T) {
			h := hasher{}

			got := getHash(t, func() ([]byte, error) {
				return h.hashList(tc.kind, tc.value)
			})

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("protohash (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.want, objectHash(t, tc.value)); diff != "" {
				t.Errorf("objecthash (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHashMessage(t *testing.T) {
	files := unmarshalProtoRegistryFiles(t, testProtoset)

	t.Run("integers.proto", func(t *testing.T) {
		for name, tc := range map[string]struct {
			md              protoreflect.MessageDescriptor
			json            string
			want            string
			skipEquivalence bool
		}{
			"Int32MessageZero": {
				md:              mdByPath(t, files, "test_protos/schema/proto3/integers.proto", "Int32Message"),
				json:            `{"values": [0, 1, 2]}`,
				want:            "48594d0a83a5b3c003ec7399e6b87d946649a3893e03a166cb947be1d0754bc9",
				skipEquivalence: true, // No equivalent JSON: JSON does not have an "integer" type. All numbers are floats.
			},
		} {
			t.Run(name, func(t *testing.T) {
				h := hasher{}

				got := getHash(t, func() ([]byte, error) {
					return h.hashMessage(unmarshalJson(t, tc.md, tc.json))
				})

				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("protohash (-want +got):\n%s", diff)
				}
				if !tc.skipEquivalence {
					if diff := cmp.Diff(tc.want, jsonHash(t, tc.json)); diff != "" {
						t.Errorf("jsonhash (-want +got):\n%s", diff)
					}
				}
			})
		}
	})

	t.Run("floats.proto", func(t *testing.T) {
		for name, tc := range map[string]struct {
			md              protoreflect.MessageDescriptor
			json            string
			want            string
			skipEquivalence bool
		}{
			"FloatMessage": {
				md:              mdByPath(t, files, "test_protos/schema/proto3/floats.proto", "FloatMessage"),
				json:            `{"values": [-2, -1, 0, 1, 2]}`,
				want:            "cd587478cd18e449dcb50fc6bb2e7322d51c259f8eae35218a42e608f8817c54",
				skipEquivalence: true, // No equivalent JSON: JSON does not have an "integer" type. All numbers are floats.
			},
		} {
			t.Run(name, func(t *testing.T) {
				h := hasher{}

				got := getHash(t, func() ([]byte, error) {
					return h.hashMessage(unmarshalJson(t, tc.md, tc.json))
				})

				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("protohash (-want +got):\n%s", diff)
				}
				if !tc.skipEquivalence {
					if diff := cmp.Diff(tc.want, jsonHash(t, tc.json)); diff != "" {
						t.Errorf("jsonhash (-want +got):\n%s", diff)
					}
				}
			})
		}
	})

}

func unmarshalJson(t *testing.T, md protoreflect.MessageDescriptor, json string) protoreflect.Message {
	msg := dynamicpb.NewMessage(md)
	if err := protojson.Unmarshal([]byte(json), msg); err != nil {
		t.Fatal(err)
	}
	return msg
}

func fdByPath(t *testing.T, files *protoregistry.Files, filename string) protoreflect.FileDescriptor {
	fd, err := files.FindFileByPath(filename)
	if err != nil {
		t.Fatal(err)
	}
	return fd
}

func mdByPath(t *testing.T, files *protoregistry.Files, filename, name string) protoreflect.MessageDescriptor {
	fd := fdByPath(t, files, filename)
	md := fd.Messages().ByName(protoreflect.Name(name))
	if md == nil {
		t.Fatal(filename, "| message descriptor not found:", name)
	}
	return md
}

func getHash(t *testing.T, fn func() ([]byte, error)) string {
	hash, err := fn()
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("%x", hash)
}

func objectHash(t *testing.T, value interface{}) string {
	objh, err := objecthash.ObjectHash(value)
	if err != nil {
		t.Fatal(err)
	}
	hash := fmt.Sprintf("%x", objh)
	if err != nil {
		t.Fatal(err)
	}
	return hash
}

func jsonHash(t *testing.T, value string) string {
	objh, err := objecthash.CommonJSONHash(value)
	if err != nil {
		t.Fatal(err)
	}
	hash := fmt.Sprintf("%x", objh)
	if err != nil {
		t.Fatal(err)
	}
	return hash
}

func unmarshalFileDescriptorSet(t *testing.T, data []byte) *descriptorpb.FileDescriptorSet {
	var dpb descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(data, &dpb); err != nil {
		t.Fatalf("unmarshaling protoset file: %v", err)
	}
	return &dpb
}

func unmarshalProtoRegistryFiles(t *testing.T, data []byte) *protoregistry.Files {
	descriptor := unmarshalFileDescriptorSet(t, data)
	files, err := protodesc.NewFiles(descriptor)
	if err != nil {
		t.Fatal(err)
	}
	return files
}

type stringList []string

// Len reports the number of entries in the List.
// Get, Set, and Truncate panic with out of bound indexes.
func (ss stringList) Len() int {
	return len(ss)
}

// Get retrieves the value at the given index.
// It never returns an invalid value.
func (ss stringList) Get(i int) protoreflect.Value {
	return protoreflect.ValueOf(ss[i])
}

// Set stores a value for the given index.
// When setting a composite type, it is unspecified whether the set
// value aliases the source's memory in any way.
//
// Set is a mutating operation and unsafe for concurrent use.
func (ss stringList) Set(i int, v protoreflect.Value) {
	ss[i] = v.String()
}

// Append appends the provided value to the end of the list.
// When appending a composite type, it is unspecified whether the appended
// value aliases the source's memory in any way.
//
// Append is a mutating operation and unsafe for concurrent use.
func (ss stringList) Append(protoreflect.Value) {
	log.Panicln("not implemented")
}

// AppendMutable appends a new, empty, mutable message value to the end
// of the list and returns it.
// It panics if the list does not contain a message type.
func (ss stringList) AppendMutable() protoreflect.Value {
	log.Panicln("not implemented")
	return protoreflect.Value{}
}

// Truncate truncates the list to a smaller length.
//
// Truncate is a mutating operation and unsafe for concurrent use.
func (ss stringList) Truncate(len int) {
	log.Panicln("not implemented")
}

// NewElement returns a new value for a list element.
// For enums, this returns the first enum value.
// For other scalars, this returns the zero value.
// For messages, this returns a new, empty, mutable value.
func (ss stringList) NewElement() protoreflect.Value {
	log.Panicln("not implemented")
	return protoreflect.Value{}
}

// IsValid reports whether the list is valid.
//
// An invalid list is an empty, read-only value.
//
// Validity is not part of the protobuf data model, and may not
// be preserved in marshaling or other operations.
func (ss stringList) IsValid() bool {
	return true
}
