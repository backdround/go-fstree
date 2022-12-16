package fstree

import (
	"github.com/backdround/go-fstree/checker"
	"github.com/backdround/go-fstree/config"
	"github.com/backdround/go-fstree/osfs"
)

// FS describes required interface for filesystem checking.
// In the most cases it copies os package signatures.
type CheckFS interface {
	IsExist(path string) bool
	IsFile(path string) bool
	IsLink(path string) bool
	IsDirectory(path string) bool

	ReadDir(path string) ([]string, error)
	ReadFile(path string) ([]byte, error)
	Readlink(path string) (string, error)
}

// Difference type describes specific difference between filesystem
// and check expectation
type Difference struct {
	Path        string
	Expectation string
	Real        string
}

// Check checks filesystem tree in rootPath by yamlData.
// For example:
//
//	configs:
//	  config1.txt:
//	    type: file
//	    data: some data
//	pkg:
//	  pkg1:
//	    type: link
//	    path: "../../pkg1"
//
// The function checks:
//   - that ./configs/config1.txt is a file with data "some data"
//   - that ./pkg/pkg1 is a link that points to "../../pkg1"
func Check(fs CheckFS, rootPath string, yamlData string) (*Difference, error) {
	// Parses config
	directoryEntry, err := config.Parse(yamlData)
	if err != nil {
		return nil, err
	}

	// Checks fs tree
	checker := checker.Checker{
		Fs: fs,
	}
	difference, err := checker.Check(rootPath, *directoryEntry)
	return (*Difference)(difference), err
}

// CheckOverOSFS makes the same thing as Chek, but uses the
// real filesystem
func CheckOverOSFS(rootPath string, yamlData string) (*Difference, error) {
	fs := osfs.OsFS{}
	return Check(fs, rootPath, yamlData)
}
