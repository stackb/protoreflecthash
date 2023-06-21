package protoreflecthash

import (
	_ "embed"
	"encoding/json"
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

	pb2_latest "github.com/stackb/protoreflecthash/test_protos/generated/latest/proto2"
	pb3_latest "github.com/stackb/protoreflecthash/test_protos/generated/latest/proto3"
)

//go:embed testdata/protoset.pb
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
	for name, tc := range map[string]struct {
		value        proto.Message
		mapFieldName string
		obj          interface{}
		json         string
		want         string
	}{
		"IntMaps.int_to_string": {
			value:        &pb3_latest.IntMaps{IntToString: map[int64]string{0: "ZERO"}},
			mapFieldName: "int_to_string",
			obj:          map[int64]string{0: "ZERO"},
			// json:         `{0:"ZERO"}`, // can't use json representation in this case
			want: "8cda73a524d09ce6fa10b071cacd4c725521b660ee4a546b6ebdbf139370e9b9",
		},
		"StringMaps.string_to_bool": {
			value:        &pb3_latest.StringMaps{StringToBool: map[string]bool{"true": true}},
			mapFieldName: "string_to_bool",
			obj:          map[string]bool{"true": true},
			json:         `{"true":true}`,
			want:         "d84d7d0593f90628672ccc4fbc89e31c51a847f45f39d98b95ea032c8de25e64",
		},
		"StringMaps.string_to_string": {
			value:        &pb3_latest.StringMaps{StringToString: map[string]string{"foo": "bar"}},
			mapFieldName: "string_to_string",
			obj:          map[string]string{"foo": "bar"},
			json:         `{"foo":"bar"}`,
			want:         "7ef5237c3027d6c58100afadf37796b3d351025cf28038280147d42fdc53b960",
		},
		"StringMaps.string_to_string_k123": {
			value:        &pb3_latest.StringMaps{StringToString: map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"}},
			mapFieldName: "string_to_string",
			obj:          map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"},
			json:         `{"k1":"v1","k2":"v2","k3":"v3"}`,
			want:         "ddd65f1f7568269a30df7cafc26044537dc2f02a1a0d830da61762fc3e687057",
		},
		"StringMaps.string_to_string_k213": {
			value:        &pb3_latest.StringMaps{StringToString: map[string]string{"k2": "v2", "k1": "v1", "k3": "v3"}},
			mapFieldName: "string_to_string",
			obj:          map[string]string{"k2": "v2", "k1": "v1", "k3": "v3"},
			json:         `{"k1":"v1","k2":"v2","k3":"v3"}`,
			want:         "ddd65f1f7568269a30df7cafc26044537dc2f02a1a0d830da61762fc3e687057",
		},
		//
	} {
		t.Run(name, func(t *testing.T) {
			h := hasher{}

			msg := tc.value.ProtoReflect()
			fd := msg.Descriptor().Fields().ByName(protoreflect.Name(tc.mapFieldName))

			got := getHash(t, func() ([]byte, error) {
				return h.hashMap(fd.MapKey(), fd.MapValue(), msg.Get(fd).Map())
			})

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("protohash (-want +got):\n%s", diff)
			}
			if tc.json != "" {
				if diff := cmp.Diff(tc.want, jsonHash(t, tc.json)); diff != "" {
					t.Errorf("jsonhash (-want +got):\n%s", diff)
				}
			}
			if tc.obj != nil {
				if diff := cmp.Diff(tc.want, objectHash(t, tc.obj)); diff != "" {
					t.Errorf("objecthash (-want +got):\n%s", diff)
				}
				if tc.json != "" {
					if diff := cmp.Diff(tc.json, jsonString(t, tc.obj)); diff != "" {
						t.Errorf("jsonstring (-want +got):\n%s", diff)
					}
				}
			}
		})
	}
}

func TestHashEnum(t *testing.T) {
	for name, tc := range map[string]struct {
		value protoreflect.EnumNumber
		want  string
	}{
		"zero": {
			value: 0,
			want:  "a4e167a76a05add8a8654c169b07b0447a916035aef602df103e8ae0fe2ff390",
		},
		"earth": {
			value: pb3_latest.PlanetV1_EARTH_V1.Number(),
			want:  "9a83c6cb1126d93de4a30715b28f1f4b26b983c57fb39e6d826d7e893ae4ee74",
		},
	} {
		t.Run(name, func(t *testing.T) {
			h := hasher{}

			got := getHash(t, func() ([]byte, error) {
				return h.hashEnum(tc.value)
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

func TestHashEmpty(t *testing.T) {
	emptyHash := "18ac3e7343f016890c510e93f935261169d9e3f565436429830faf0934f4f8e4"
	emptyJson := "{}"

	if diff := cmp.Diff(emptyHash, jsonHash(t, emptyJson)); diff != "" {
		t.Errorf("jsonhash (-want +got):\n%s", diff)
	}

	for _, tc := range []struct {
		msg proto.Message
	}{
		{&pb3_latest.Empty{}},
		// Empty repeated fields are ignored.
		{&pb3_latest.Repetitive{StringField: []string{}}},
		// Empty map fields are ignored.
		{&pb3_latest.StringMaps{StringToString: map[string]string{}}},
		// Proto3 scalar fields set to their default values are considered empty.
		{&pb3_latest.Simple{BoolField: false}},
		{&pb3_latest.Simple{BytesField: []byte{}}},
		{&pb3_latest.Simple{DoubleField: 0}},
		{&pb3_latest.Simple{DoubleField: 0.0}},
		{&pb3_latest.Simple{Fixed32Field: 0}},
		{&pb3_latest.Simple{Fixed64Field: 0}},
		{&pb3_latest.Simple{FloatField: 0}},
		{&pb3_latest.Simple{FloatField: 0.0}},
		{&pb3_latest.Simple{Int32Field: 0}},
		{&pb3_latest.Simple{Int64Field: 0}},
		{&pb3_latest.Simple{Sfixed32Field: 0}},
		{&pb3_latest.Simple{Sfixed64Field: 0}},
		{&pb3_latest.Simple{Sint32Field: 0}},
		{&pb3_latest.Simple{Sint64Field: 0}},
		{&pb3_latest.Simple{StringField: ""}},
		{&pb3_latest.Simple{Uint32Field: 0}},
		{&pb3_latest.Simple{Uint64Field: 0}},
	} {
		t.Run(fmt.Sprintf("%+v", tc.msg), func(t *testing.T) {
			h := hasher{}

			got := getHash(t, func() ([]byte, error) {
				return h.hashMessage(tc.msg.ProtoReflect())
			})

			if diff := cmp.Diff(emptyHash, got); diff != "" {
				t.Errorf("protohash (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHashIntegerFields(t *testing.T) {

	for name, tc := range map[string]struct {
		fieldNamesAsKeys bool
		protos           []proto.Message
		obj              interface{}
		json             string
		want             string
	}{
		"equivalence": {
			protos: []proto.Message{},
			// obj: map[string][]float64{"values": {-2, -1, 0, 1, 2}},
			// json: `{2: [-2, -1, 0, 1, 2]}`, skipping json as this is invalid json
			want: "08775d05cd028265e4956a95aef6c050a45652e9c59462da636a8460c5ed52f3",
		},
	} {
		t.Run(name, func(t *testing.T) {
			for _, msg := range tc.protos {
				t.Run(fmt.Sprintf("%+v", msg), func(t *testing.T) {
					h := hasher{fieldNamesAsKeys: tc.fieldNamesAsKeys}

					got := getHash(t, func() ([]byte, error) {
						return h.hashMessage(msg.ProtoReflect())
					})

					if diff := cmp.Diff(tc.want, got); diff != "" {
						t.Errorf("protohash (-want +got):\n%s", diff)
					}

					if tc.json != "" {
						if diff := cmp.Diff(tc.want, jsonHash(t, tc.json)); diff != "" {
							t.Errorf("jsonhash (-want +got):\n%s", diff)
						}
					}
					if tc.obj != nil {
						if diff := cmp.Diff(tc.want, objectHash(t, tc.obj)); diff != "" {
							t.Errorf("objecthash (-want +got):\n%s", diff)
						}
					}
				})
			}
		})
	}
}

func TestHashFloatFields(t *testing.T) {

	for name, tc := range map[string]struct {
		fieldNamesAsKeys bool
		protos           []proto.Message
		obj              interface{}
		json             string
		want             string
	}{
		"float fields (hashing key field numbers)": {
			protos: []proto.Message{
				&pb2_latest.DoubleMessage{Values: []float64{-2, -1, 0, 1, 2}},
				&pb3_latest.DoubleMessage{Values: []float64{-2, -1, 0, 1, 2}},
				&pb2_latest.FloatMessage{Values: []float32{-2, -1, 0, 1, 2}},
				&pb3_latest.FloatMessage{Values: []float32{-2, -1, 0, 1, 2}},
			},
			// obj: map[string][]float64{"values": {-2, -1, 0, 1, 2}},
			// json: `{2: [-2, -1, 0, 1, 2]}`, skipping json as this is invalid json
			want: "08775d05cd028265e4956a95aef6c050a45652e9c59462da636a8460c5ed52f3",
		},
		"float fields (hashing key field strings)": {
			fieldNamesAsKeys: true,
			protos: []proto.Message{
				&pb2_latest.DoubleMessage{Values: []float64{-2, -1, 0, 1, 2}},
				&pb3_latest.DoubleMessage{Values: []float64{-2, -1, 0, 1, 2}},
				&pb2_latest.FloatMessage{Values: []float32{-2, -1, 0, 1, 2}},
				&pb3_latest.FloatMessage{Values: []float32{-2, -1, 0, 1, 2}},
			},
			obj:  map[string][]float64{"values": {-2, -1, 0, 1, 2}},
			json: `{"values": [-2, -1, 0, 1, 2]}`,
			want: "586202dddb0e98bb8ce0b7289e29a9f7397b9b1996f3f8fe788f4cfb230b7ee8",
		},
		"float fields (fractions 32)": {
			fieldNamesAsKeys: true,
			protos: []proto.Message{
				&pb2_latest.DoubleMessage{Values: []float64{0.0078125, 7.888609052210118e-31}},
				&pb3_latest.DoubleMessage{Values: []float64{0.0078125, 7.888609052210118e-31}},
				&pb2_latest.FloatMessage{Values: []float32{0.0078125, 7.888609052210118e-31}},
				&pb3_latest.FloatMessage{Values: []float32{0.0078125, 7.888609052210118e-31}},
			},
			obj:  map[string][]float64{"values": {0.0078125, 7.888609052210118e-31}},
			json: `{"values": [0.0078125, 7.888609052210118e-31]}`,
			want: "7b7cba0ed312bc6611f0523e7c46ce9a2ed9ecb798eb80e1cdf93c95faf503c7",
		},
		"float fields (fractions 64)": {
			fieldNamesAsKeys: true,
			protos: []proto.Message{
				&pb2_latest.DoubleMessage{Values: []float64{-1.0, 1.5, 1000.000244140625, 1267650600228229401496703205376, 32.0, 13.0009765625}},
				&pb3_latest.DoubleMessage{Values: []float64{-1.0, 1.5, 1000.000244140625, 1267650600228229401496703205376, 32.0, 13.0009765625}},
				&pb2_latest.FloatMessage{Values: []float32{-1.0, 1.5, 1000.000244140625, 1267650600228229401496703205376, 32.0, 13.0009765625}},
				&pb3_latest.FloatMessage{Values: []float32{-1.0, 1.5, 1000.000244140625, 1267650600228229401496703205376, 32.0, 13.0009765625}},
			},
			json: `{"values": [-1.0, 1.5, 1000.000244140625, 1267650600228229401496703205376, 32.0, 13.0009765625]}`,
			want: "ac261ff3d8b933998e3fea278539eb40b15811dd835d224e0150dce4794168b7",
		},
		"float fields (Non-equivalence of Floats using different representations)": {
			fieldNamesAsKeys: true,
			protos: []proto.Message{
				&pb2_latest.FloatMessage{Value: proto.Float32(0.1)},
				&pb3_latest.FloatMessage{Value: 0.1},
				// A float64 "0.1" is not equal to a float32 "0.1".
				// However, float32 "0.1" is equal to float64 "1.0000000149011612e-1".
				&pb2_latest.DoubleMessage{Value: proto.Float64(1.0000000149011612e-1)},
				&pb3_latest.DoubleMessage{Value: 1.0000000149011612e-1},
			},
			obj:  map[string]float32{"value": 0.1},
			json: `{"value": 1.0000000149011612e-1}`,
			want: "7081ed6a1e7ad8e7f981a2894a3bd6d3b0b0033b69c03cce84b61dd063f4efaa",
		},
		"float fields (There's no float32 number that is equivalent to a float64 '0.1'.)": {
			fieldNamesAsKeys: true,
			protos: []proto.Message{
				&pb2_latest.DoubleMessage{Value: proto.Float64(0.1)},
				&pb3_latest.DoubleMessage{Value: 0.1},
			},
			obj:  map[string]float64{"value": 0.1},
			json: `{"value": 0.1}`,
			want: "e175fbe785bae88b598d3ecaad8a64d2a998e9f673173a226868f2ef312a5225",
		},
		"float fields (Non-equivalence of Floats using different representations - decimal)": {
			fieldNamesAsKeys: true,
			protos: []proto.Message{
				&pb2_latest.FloatMessage{Value: proto.Float32(1.2163543e+25)},
				&pb3_latest.FloatMessage{Value: 1.2163543e+25},
				// The decimal representation of the equivalent 64-bit float is different.
				&pb2_latest.DoubleMessage{Value: proto.Float64(1.2163543234531120e+25)},
				&pb3_latest.DoubleMessage{Value: 1.2163543234531120e+25},
			},
			obj:  map[string]float32{"value": 1.2163543e+25},
			json: `{"value": 1.2163543234531120e+25}`,
			want: "bbb17cf7312f2ba5b0002d781f16d1ab50c3d25dc044ed3428750826a1c68653",
		},
		"float fields (no float32 number that is equivalent to a float64 '1e+25')": {
			fieldNamesAsKeys: true,
			protos: []proto.Message{
				&pb2_latest.DoubleMessage{Value: proto.Float64(1e+25)},
				&pb3_latest.DoubleMessage{Value: 1e+25},
			},
			obj:  map[string]float64{"value": 1e+25},
			json: `{"value": 1e+25}`,
			want: "874beabbede24974a9f3f74e3448670e0c42c0aaba082f18b963b72253649362",
		},
		"float fields (proto2 unset)": {
			fieldNamesAsKeys: true,
			protos: []proto.Message{
				&pb2_latest.DoubleMessage{Value: proto.Float64(0)},
				&pb2_latest.FloatMessage{Value: proto.Float32(0)},
			},
			obj:  map[string]float64{"value": 0},
			json: `{"value":0}`,
			want: "94136b0850db069dfd7bee090fc7ede48aa7da53ae3cc8514140a493818c3b91",
		},
		"float fields (special NaN)": {
			fieldNamesAsKeys: true,
			protos: []proto.Message{
				&pb2_latest.DoubleMessage{Value: proto.Float64(math.NaN())},
				&pb3_latest.DoubleMessage{Value: math.NaN()},

				&pb2_latest.FloatMessage{Value: proto.Float32(float32(math.NaN()))},
				&pb3_latest.FloatMessage{Value: float32(math.NaN())},
			},
			obj: map[string]float64{"value": math.NaN()},
			// No equivalent JSON: JSON does not support special float values.
			// See: https://tools.ietf.org/html/rfc4627#section-2.4
			// json: `{"value": NaN}`,
			want: "16614de29b0823c41cabc993fa6c45da87e4e74c5d836edbcddcfaaf06ffafd1",
		},
		"float fields (special Inf(+))": {
			fieldNamesAsKeys: true,
			protos: []proto.Message{
				&pb2_latest.DoubleMessage{Value: proto.Float64(math.Inf(1))},
				&pb3_latest.DoubleMessage{Value: math.Inf(1)},

				&pb2_latest.FloatMessage{Value: proto.Float32(float32(math.Inf(1)))},
				&pb3_latest.FloatMessage{Value: float32(math.Inf(1))},
			},
			obj: map[string]float64{"value": math.Inf(1)},
			// No equivalent JSON: JSON does not support special float values.
			// See: https://tools.ietf.org/html/rfc4627#section-2.4
			// json: `{"value": Inf}`,
			want: "c58cd512e86204e99cb6c11d83bb3daaccdd946e66383004cb9b7f87f762935c",
		},
		"float fields (special Inf(-))": {
			fieldNamesAsKeys: true,
			protos: []proto.Message{
				&pb2_latest.DoubleMessage{Value: proto.Float64(math.Inf(-1))},
				&pb3_latest.DoubleMessage{Value: math.Inf(-1)},

				&pb2_latest.FloatMessage{Value: proto.Float32(float32(math.Inf(-1)))},
				&pb3_latest.FloatMessage{Value: float32(math.Inf(-1))},
			},
			obj: map[string]float64{"value": math.Inf(-1)},
			// No equivalent JSON: JSON does not support special float values.
			// See: https://tools.ietf.org/html/rfc4627#section-2.4
			// json: `{"value": Inf}`,
			want: "1a4ffd7e9dc1f915c5b3b821d9194ac7d6d2bdec947aa8c3b3b1e9017c651331",
		},
	} {
		t.Run(name, func(t *testing.T) {
			for _, msg := range tc.protos {
				t.Run(fmt.Sprintf("%+v", msg), func(t *testing.T) {
					h := hasher{fieldNamesAsKeys: tc.fieldNamesAsKeys}

					got := getHash(t, func() ([]byte, error) {
						return h.hashMessage(msg.ProtoReflect())
					})

					if diff := cmp.Diff(tc.want, got); diff != "" {
						t.Errorf("protohash (-want +got):\n%s", diff)
					}

					if tc.json != "" {
						if diff := cmp.Diff(tc.want, jsonHash(t, tc.json)); diff != "" {
							t.Errorf("jsonhash (-want +got):\n%s", diff)
						}
					}
					if tc.obj != nil {
						if diff := cmp.Diff(tc.want, objectHash(t, tc.obj)); diff != "" {
							t.Errorf("objecthash (-want +got):\n%s", diff)
						}
					}
				})
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
				want:            "ec28f92dbcce2dc9e38b48cd7725337ca7df40d729b8523a5b3512f7449e8156",
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
				want:            "08775d05cd028265e4956a95aef6c050a45652e9c59462da636a8460c5ed52f3",
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

func jsonString(t *testing.T, value interface{}) string {
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
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
