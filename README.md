# FSTree

FSTree makes a filesystem tree from a given yaml definition. This
package is intended to use in tests.

### Installation

```bash
go get github.com/backdround/go-fstree
```

### Example

```go
package main

import (
  "strings"
  "github.com/backdround/go-fstree"
)

var fstreeYaml =`
configs:
  "config1.txt":
    type: file
    data: "format: txt"
  "config2.txt":
    type: file
    data: "format: yaml"
pkg:
  pkg1:
    type: link
    path: "../../pkg1"
`

func main() {
  fstreeYaml = strings.ReplaceAll(fstreeYaml, "\t", "  ")

  // Creates filesystem tree in ./project
  err = fstree.Make("./project", fstreeYaml)

  if err != nil {
    log.Fatal(err)
  }
}
```
