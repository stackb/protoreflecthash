package protoreflecthash

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func Bytes(msg protoreflect.Message) ([]byte, error) {
	h := hasher{}
	hash, err := h.hashMessage(msg)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func String(msg protoreflect.Message) (string, error) {
	hash, err := Bytes(msg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash), nil
}
