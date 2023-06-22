[![CI](https://github.com/stackb/protoreflecthash/actions/workflows/ci.yaml/badge.svg)](https://github.com/stackb/protoreflecthash/actions/workflows/ci.yaml)

<table border="0">
  <tr>
    <td><img src="https://user-images.githubusercontent.com/50580/141900696-bfb2d42d-5d2c-46f8-bd9f-06515969f6a2.png" height="120"/></td>
    <td><img src="https://camo.githubusercontent.com/e71e893edbc4626f5c28acde484c375e96109e64b835db9b60946ff092a7a87b/68747470733a2f2f6769746c61622e636f6d2f6d30336765656b2f6e6f64652d6f626a6563742d686173682f7261772f6d61737465722f6c6f676f2e737667" height="120"/></td>
  </tr>
  <tr>
    <td>protobuf</td>
    <td>objecthash</td>
  </tr>
</table>

# protoreflecthash

protoreflecthash is a re-implementation of
<https://github.com/deepmind/objecthash-proto>.  That repo is now archived and
was never updated for [protobuf-apiv2](https://go.dev/blog/protobuf-apiv2).

# Usage

```
go get github.com/stackb/protoreflecthash@latest
```

```go
package main

import (
    "github.com/stackb/protoreflecthash"
)

func main() {
    msg := mustGetProtoMessageSomewhere()

    options := []protoreflect.Option{
      // protoreflect.MessageFullnameIdentifier(),
      // protoreflect.FieldNamesAsKeys(),
    }
    hasher := protoreflect.NewHasher(options...)

    hash, err := hasher.HashProto(msg.ProtoReflect())
    if err != nil {
        panic(err.Error())
    }

    fmt.Printf("%x\n", hash)
}
```

# Background

`protoreflecthash` computes the hash value for a protobuf message by taking a
sha256 of the sum the individual component hashes of the message.  Special care
is taken to account for various semantics of the protobuf format.

This implementation passes all functional unit tests from the original library
[deepmind/objecthash-proto](https://github.com/deepmind/objecthash-proto)
(excluding badness detection).

Open questions remain about the handling of protobufs with extension fields and
the google.protobuf.Any type.

This package is currently experimental; hash values for messages may change
without warning until v1.
