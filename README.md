[![Go Reference](https://img.shields.io/badge/go-reference-%2300ADD8?style=flat-square)](https://pkg.go.dev/github.com/backdround/go-fstree)
[![Tests](https://img.shields.io/github/workflow/status/backdround/go-fstree/tests?label=tests&style=flat-square)](https://github.com/backdround/go-fstree/actions)
[![Codecov](https://img.shields.io/codecov/c/github/backdround/go-fstree?style=flat-square)](https://app.codecov.io/gh/backdround/go-fstree/)
[![Go Report](https://goreportcard.com/badge/github.com/backdround/go-fstree?style=flat-square)](https://goreportcard.com/report/github.com/backdround/go-fstree)

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
	"log"
	"strings"

	"github.com/backdround/go-fstree"
)

var fstreeYaml =`
configs:
  config1.txt:
    type: file
    data: "format: txt"
  config2.txt:
    type: file
    data: "format: yaml"
pkg:
  pkg1:
    type: link
    path: ../../pkg1
`

func main() {
	fstreeYaml = strings.ReplaceAll(fstreeYaml, "\t", "  ")
	
	// Creates filesystem tree in ./project
	err := fstree.MakeOverOSFS("./project", fstreeYaml)
	
	if err != nil {
		log.Fatal(err)
	}
}
```

### FS Entries

#### Directory
```yaml
depth0:
  # No type field is setted
  depth1:
    # No type field is setted
    depth2:
      # No type field is setted
```
creates `ROOTPATH/depth0/depth1/depath2` directory

#### File
```yaml
file1.txt:
  # type is required
  type: file
  # data is optinal
  data: some file data
```
creates file `ROOTPATH/file1.txt` with data `some file data`

#### Link
```yaml
link1:
  # type is reqired
  type: link
  # path is required
  path: ./some/destination
```
creates link `ROOTPATH/link1` with destination `./some/destination`
