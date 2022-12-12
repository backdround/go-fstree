package fstree

import (
	"github.com/backdround/go-fstree/config"
	"github.com/backdround/go-fstree/fstreemaker"
	"github.com/backdround/go-fstree/osfs"
	"github.com/backdround/go-fstree/types"
)

// Make makes filesystem tree in rootPath from yamlData.
// For example with yaml data:
/*
configs:
  "config1.txt":
    type: file
    data: "format: txt"
pkg:
  pkg1:
    type: link
    path: "../../pkg1"
*/
// The function creates:
// - ./configs/config1.txt (file with data "format: txt")
// - ./pkg/pkg1 (link points to "../../pkg1")
func Make(fs types.FS, rootPath string, yamlData string) error {
	var err error
	// Parses config
	directoryEntry, err := config.Parse(rootPath, yamlData)
	if err != nil {
		return err
	}

	// Creates root directory
	if !fs.IsDirectory(rootPath) {
		err := fs.Mkdir(rootPath)
		if err != nil {
			return err
		}
	}

	// Creates fs tree
	maker := fstreemaker.Maker{
		Fs: fs,
	}
	err = maker.MakeDirectory(rootPath, *directoryEntry)
	return err
}

// MakeOverOSFS makes the same thing as Make, but uses the
// real filesystem
func MakeOverOSFS(rootPath string, yamlData string) error {
	fs := osfs.OsFS{}
	return Make(fs, rootPath, yamlData)
}
