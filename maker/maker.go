// Package maker makes filesystem tree by passed filesystem entries
package maker

import (
	"bytes"
	"fmt"
	"path"

	"github.com/backdround/go-fstree/entries"
)

type Maker struct {
	Fs FS
}

// makeFile creates a file in the workDirectory. It skips if file with the
// same data exists. Gives a error if by the filepath something exists.
func (m Maker) makeFile(workDirectory string, file entries.FileEntry) error {
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
func (m Maker) makeLink(workDirectory string, link entries.LinkEntry) error {
	linkPath := path.Join(workDirectory, link.Name)

	if !m.Fs.IsExist(linkPath) {
		return m.Fs.Symlink(link.Path, linkPath)
	}

	if !m.Fs.IsLink(linkPath) {
		return fmt.Errorf("filepath %q already exists", linkPath)
	}

	existingLinkDestination, err := m.Fs.Readlink(linkPath)
	if err != nil {
		return err
	}

	matched, err := path.Match(existingLinkDestination, link.Path)
	if err != nil {
		panic(err)
	}

	if !matched {
		return fmt.Errorf("link %q already exists", linkPath)
	}

	return nil
}

// makeDirectory creates directory in workDirectory
func (m Maker) MakeDirectory(workDirectory string,
	directory entries.DirectoryEntry) error {
	dirPath := path.Join(workDirectory, directory.Name)

	// Creates current directory
	if !m.Fs.IsDirectory(dirPath) {
		if m.Fs.IsExist(dirPath) {
			return fmt.Errorf("filepath %q already exists", dirPath)
		}

		err := m.Fs.Mkdir(dirPath)
		if err != nil {
			return err
		}
	}

	// Creates directory entries
	for _, entry := range directory.Entries {
		switch entry.(type) {
		case entries.FileEntry:
			fileEntry := entry.(entries.FileEntry)
			err := m.makeFile(dirPath, fileEntry)
			if err != nil {
				return err
			}
		case entries.LinkEntry:
			linkEntry := entry.(entries.LinkEntry)
			err := m.makeLink(dirPath, linkEntry)
			if err != nil {
				return err
			}
		case entries.DirectoryEntry:
			directoryEntry := entry.(entries.DirectoryEntry)
			err := m.MakeDirectory(dirPath, directoryEntry)
			if err != nil {
				return err
			}
		default:
			panic("unknown entry type")
		}
	}

	return nil
}
