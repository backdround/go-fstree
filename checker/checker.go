package checker

import (
	"github.com/backdround/go-fstree/entries"
)

type Checker struct {
	Fs FS
}

func (c Checker) Check(rootPath string, expectedTree entries.DirectoryEntry) (
	difference string, err error) {
	return "", nil
}
