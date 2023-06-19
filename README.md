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

This package is currently experimental; the hash values for messages will likely
change without warning until v1.

# Usage

```
go get github.com/stackb/protoreflecthash
```

```go
package main

import (
    "github.com/stackb/protoreflecthash"
)

func main() {
    msg := mustGetProtoMessageSomewhere()
    hex, err := protoreflecthash.String(msg)
    if err != nil {
        log.Fatal(err)
    }
    println(hex)
}
```
