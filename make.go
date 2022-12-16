package fstree

import (
	"github.com/backdround/go-fstree/config"
	"github.com/backdround/go-fstree/maker"
	"github.com/backdround/go-fstree/osfs"
)

// MakerFS describes required interface for making filetree.
// In the most cases it copies os package signatures.
type MakerFS interface {
	IsExist(path string) bool
	IsFile(path string) bool
	IsLink(path string) bool
	IsDirectory(path string) bool

	ReadFile(path string) ([]byte, error)
	Readlink(path string) (string, error)
	WriteFile(path string, data []byte) error
	Symlink(oldPath, newPath string) error
	Mkdir(path string) error
}

// Make makes filesystem tree in rootPath from yamlData.
// For example:
//
//	configs:
//	  "config1.txt":
//	    type: file
//	    data: "format: txt"
//	pkg:
//	  pkg1:
//	    type: link
//	    path: "../../pkg1"
//
// The function creates:
//   - ./configs/config1.txt (file with data "format: txt")
//   - ./pkg/pkg1 (link points to "../../pkg1")
func Make(fs MakerFS, rootPath string, yamlData string) error {
	var err error
	// Parses config
	directoryEntry, err := config.Parse(yamlData)
	if err != nil {
		return err
	}

	// Creates fs tree
	maker := maker.Maker{
		Fs: fs,
	}
	err = maker.Make(rootPath, *directoryEntry)
	return err
}

// MakeOverOSFS makes the same thing as Make, but uses the
// real filesystem
func MakeOverOSFS(rootPath string, yamlData string) error {
	fs := osfs.OsFS{}
	return Make(fs, rootPath, yamlData)
}
