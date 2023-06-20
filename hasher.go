package protoreflecthash

import (
	"bytes"
	"fmt"
	"sort"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ProtoHasher interface {
	HashProto(msg protoreflect.Message) ([]byte, error)
}

func NewHasher() ProtoHasher {
	return &hasher{}
}

type hasher struct {
}

// HashProto implements MessageHasher
func (h *hasher) HashProto(msg protoreflect.Message) ([]byte, error) {
	// Check if the value is nil.
	if msg == nil {
		return h.hashNil()
	}

	// Make sure the proto itself is actually valid (ie. can be marshalled).
	// If this fails, it probably means there are unset required fields or invalid
	// values.
	if _, err := proto.Marshal(msg.Interface()); err != nil {
		return nil, err
	}

	return h.hashMessage(msg)
}

func (h *hasher) hashMessage(msg protoreflect.Message) ([]byte, error) {
	md := msg.Descriptor()

	// TOOD(pcj): what is the correct handling of placeholder types?
	if md.IsPlaceholder() {
		return nil, nil
	}

	var buf bytes.Buffer

	fhash, err := h.hashFields(msg, md.Fields())
	if err != nil {
		return nil, fmt.Errorf("hashing fields: %w", err)
	}
	buf.Write(fhash)

	ohash, err := h.hashOneofs(msg, md.Oneofs())
	if err != nil {
		return nil, fmt.Errorf("hashing fields: %w", err)
	}
	buf.Write(ohash)

	identifier := mapIdentifier
	// if hasher.messageIdentifier != "" {
	// 	identifier = hasher.messageIdentifier
	// }
	return hash(identifier, buf.Bytes())
}

func (h *hasher) hashOneofs(msg protoreflect.Message, oneofs protoreflect.OneofDescriptors) ([]byte, error) {
	type oneOfHash struct {
		number int
		data   []byte
	}

	hashes := make([]oneOfHash, oneofs.Len())
	for i := 0; i < oneofs.Len(); i++ {
		od := oneofs.Get(i)
		fields := od.Fields()
		data, err := h.hashFields(msg, fields)
		if err != nil {
			return nil, fmt.Errorf("hashing field %d (%s): %w", i, od.FullName(), err)
		}
		hashes[i] = oneOfHash{number: i, data: data}
	}

	var buf bytes.Buffer
	for _, hash := range hashes {
		buf.Write(hash.data)
	}

	return buf.Bytes(), nil
}

func (h *hasher) hashFields(msg protoreflect.Message, fields protoreflect.FieldDescriptors) ([]byte, error) {
	type fieldHash struct {
		number int
		data   []byte
	}

	hashes := make([]fieldHash, 0, fields.Len())
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		if !msg.Has(fd) {
			continue
		}
		value := msg.Get(fd)
		defaultValue := fd.Default()
		if value.Equal(defaultValue) {
			continue
		}
		data, err := h.hashField(fd, value)
		if err != nil {
			return nil, fmt.Errorf("hashing field %d (%s = %d): %w", i, fd.FullName(), fd.Number(), err)
		}
		hashes = append(hashes, fieldHash{number: int(fd.Number()), data: data})
	}
	sort.Slice(hashes, func(i, j int) bool {
		return hashes[i].number < hashes[j].number
	})

	var buf bytes.Buffer
	for _, hash := range hashes {
		buf.Write(hash.data)
	}
	return buf.Bytes(), nil
}

func (h *hasher) hashField(fd protoreflect.FieldDescriptor, value protoreflect.Value) ([]byte, error) {
	if fd.IsList() {
		return h.hashList(fd.Kind(), value.List())
	}
	if fd.IsMap() {
		return h.hashMap(fd.MapKey(), fd.MapValue(), value.Map())
	}
	return h.hashValue(fd.Kind(), value)
}

func (h *hasher) hashValue(kind protoreflect.Kind, value protoreflect.Value) ([]byte, error) {
	switch kind {
	case protoreflect.BoolKind:
		return h.hashBool(value.Bool())
	case protoreflect.EnumKind:
		return h.hashEnum(value.Enum())
	case protoreflect.Uint32Kind,
		protoreflect.Fixed32Kind,
		protoreflect.Fixed64Kind,
		protoreflect.Uint64Kind:
		return h.hashUint(value.Uint())
	case protoreflect.Int32Kind,
		protoreflect.Sint32Kind,
		protoreflect.Int64Kind,
		protoreflect.Sint64Kind,
		protoreflect.Sfixed32Kind,
		protoreflect.Sfixed64Kind:
		return h.hashInt(value.Int())
	case protoreflect.FloatKind,
		protoreflect.DoubleKind:
		return h.hashFloat(value.Float())
	case protoreflect.StringKind:
		return h.hashString(value.String())
	case protoreflect.BytesKind:
		return h.hashBytes(value.Bytes())
	case protoreflect.MessageKind:
		return h.hashMessage(value.Message())
	case protoreflect.GroupKind:
		return nil, fmt.Errorf("protoreflect.GroupKind: not implemented: %T", value)
	}
	return nil, fmt.Errorf("unexpected field kind: %v (%T)", kind, value)
}

func (h *hasher) hashNil() ([]byte, error) {
	return hashNil()
}

func (h *hasher) hashBool(value bool) ([]byte, error) {
	return hashBool(value)
}

func (h *hasher) hashEnum(value protoreflect.EnumNumber) ([]byte, error) {
	return hashInt64(int64(value))
}

func (h *hasher) hashInt(value int64) ([]byte, error) {
	return hashInt64(value)
}

func (h *hasher) hashUint(value uint64) ([]byte, error) {
	return hashUint64(value)
}

func (h *hasher) hashFloat(value float64) ([]byte, error) {
	return hashFloat(value)
}

func (h *hasher) hashString(value string) ([]byte, error) {
	return hashUnicode(value)
}

func (h *hasher) hashBytes(value []byte) ([]byte, error) {
	return hashBytes(value)
}

func (h *hasher) hashList(kind protoreflect.Kind, list protoreflect.List) ([]byte, error) {
	var buf bytes.Buffer

	for i := 0; i < list.Len(); i++ {
		value := list.Get(i)
		data, err := h.hashValue(kind, value)
		if err != nil {
			return nil, fmt.Errorf("hashing list item %d: %w", i, err)
		}
		buf.Write(data)
	}

	return hash(listIdentifier, buf.Bytes())
}

func (h *hasher) hashMap(kd, fd protoreflect.FieldDescriptor, m protoreflect.Map) ([]byte, error) {

	var mapHashEntries []hashMapEntry

	var errValue error
	var errKey protoreflect.MapKey
	m.Range(func(mk protoreflect.MapKey, v protoreflect.Value) bool {
		khash, err := h.hashField(kd, mk.Value())
		if err != nil {
			errKey = mk
			errValue = err
			return false
		}

		vhash, err := h.hashField(fd, v)
		if err != nil {
			errKey = mk
			errValue = err
			return false
		}

		mapHashEntries = append(mapHashEntries, hashMapEntry{
			khash: khash,
			vhash: vhash,
		})

		return true
	})
	if errValue != nil {
		return nil, fmt.Errorf("hashing map key %v: %w", errKey, errValue)
	}

	sort.Sort(byKHash(mapHashEntries))

	var buf bytes.Buffer
	for _, e := range mapHashEntries {
		buf.Write(e.khash[:])
		buf.Write(e.vhash[:])
	}

	return hash(mapIdentifier, buf.Bytes())
}

type hashMapEntry struct {
	khash []byte
	vhash []byte
}

type byKHash []hashMapEntry

func (h byKHash) Len() int {
	return len(h)
}

func (h byKHash) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h byKHash) Less(i, j int) bool {
	return bytes.Compare(h[i].khash[:], h[j].khash[:]) < 0
}
