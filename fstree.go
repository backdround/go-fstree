package fstree

import (
	"github.com/backdround/go-fstree/config"
	"github.com/backdround/go-fstree/fstreemaker"
	"github.com/backdround/go-fstree/types"
)

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
