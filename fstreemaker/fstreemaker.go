package fstreemaker

import (
	"bytes"
	"fmt"
	"path"

	"github.com/backdround/go-fstree/types"
)

type Maker struct {
	fs types.FS
}

// makeFile creates a file in the workDirectory. It skips if file with the
// same data exists. Gives a error if by the filepath something exists.
func (m Maker) makeFile(workDirectory string, file types.FileEntry) error {
	filePath := path.Join(workDirectory, file.Name)

	if m.fs.IsFile(filePath) {
		data, err := m.fs.ReadFile(filePath)
		if err != nil {
			return err
		}
		if !bytes.Equal(file.Data, data) {
			return fmt.Errorf("file %q already exists", filePath)
		}
		return nil
	}

	if m.fs.IsExist(filePath) {
		return fmt.Errorf("filepath %q already exists", filePath)
	}

	return m.fs.WriteFile(filePath, file.Data)
}

// makeLink creates link in workDirectory. Gives a error if by the
// filepath something exists.
func (m Maker) makeLink(workDirectory string, link types.LinkEntry) error {
	linkPath := path.Join(workDirectory, link.Name)

	if m.fs.IsExist(linkPath) {
		return fmt.Errorf("filepath %q already exists", linkPath)
	}

	return m.fs.Symlink(link.Path, linkPath)
}

// makeDirectory creates directory in workDirectory
func (m Maker) MakeDirectory(workDirectory string,
	directory types.DirectoryEntry) error {
	dirPath := path.Join(workDirectory, directory.Name)

	// Creates current directory
	if !m.fs.IsDirectory(dirPath) {
		if m.fs.IsExist(dirPath) {
			return fmt.Errorf("filepath %q already exists", dirPath)
		}

		m.fs.Mkdir(dirPath)
	}

	// Creates directory entries
	for _, entry := range directory.Entries {
		switch entry.(type) {
		case types.FileEntry:
			fileEntry := entry.(types.FileEntry)
			m.makeFile(dirPath, fileEntry)
		case types.LinkEntry:
			linkEntry := entry.(types.LinkEntry)
			m.makeLink(dirPath, linkEntry)
		case types.DirectoryEntry:
			directoryEntry := entry.(types.DirectoryEntry)
			m.MakeDirectory(dirPath, directoryEntry)
		default:
			panic("unknown entry type")
		}
	}

	return nil
}
