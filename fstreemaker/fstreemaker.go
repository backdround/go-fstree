package fstreemaker

import (
	"bytes"
	"fmt"
	"path"

	"github.com/backdround/go-fstree/types"
)

type Maker struct {
	Fs types.FS
}

// makeFile creates a file in the workDirectory. It skips if file with the
// same data exists. Gives a error if by the filepath something exists.
func (m Maker) makeFile(workDirectory string, file types.FileEntry) error {
	filePath := path.Join(workDirectory, file.Name)

	if m.Fs.IsFile(filePath) {
		data, err := m.Fs.ReadFile(filePath)
		if err != nil {
			return err
		}
		if !bytes.Equal(file.Data, data) {
			return fmt.Errorf("file %q already exists", filePath)
		}
		return nil
	}

	if m.Fs.IsExist(filePath) {
		return fmt.Errorf("filepath %q already exists", filePath)
	}

	return m.Fs.WriteFile(filePath, file.Data)
}

// makeLink creates link in workDirectory. Gives a error if by the
// filepath something exists.
func (m Maker) makeLink(workDirectory string, link types.LinkEntry) error {
	linkPath := path.Join(workDirectory, link.Name)

	if m.Fs.IsExist(linkPath) {
		return fmt.Errorf("filepath %q already exists", linkPath)
	}

	return m.Fs.Symlink(link.Path, linkPath)
}

// makeDirectory creates directory in workDirectory
func (m Maker) MakeDirectory(workDirectory string,
	directory types.DirectoryEntry) error {
	dirPath := path.Join(workDirectory, directory.Name)

	// Creates current directory
	if !m.Fs.IsDirectory(dirPath) {
		if m.Fs.IsExist(dirPath) {
			return fmt.Errorf("filepath %q already exists", dirPath)
		}

		m.Fs.Mkdir(dirPath)
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
