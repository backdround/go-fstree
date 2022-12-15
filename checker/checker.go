// Package checker makes compliance check with filesystem tree structure.
package checker

import (
	"bytes"
	"path"

	"github.com/backdround/go-fstree/entries"
)

type Checker struct {
	Fs FS
}

// Check makes compliance check with filesystem tree structure.
func (c Checker) Check(rootPath string, expectedTree entries.DirectoryEntry) (
	difference *Difference, err error) {
	return c.checkDir(rootPath, expectedTree)
}

func (c Checker) checkDir(currentPath string,
	expectedDir entries.DirectoryEntry) (difference *Difference, err error) {

	directoryPath := path.Join(currentPath, expectedDir.Name)

	// Checks that the directory exists
	if !c.Fs.IsDirectory(directoryPath) {
		difference = &Difference{
			Path:        directoryPath,
			Expectation: "directory exists",
		}

		if c.Fs.IsExist(directoryPath) {
			difference.Real = "path isn't a directory"
		} else {
			difference.Real = "directory doesn't exist"
		}

		return difference, nil
	}

	// Checks that all existing entries are expected
	existingEntryNames, err := c.Fs.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	checkIfPathIsExpected := func(entryPath string) bool {
		for _, expectedEntry := range expectedDir.Entries {
			var expectedEntryName string

			switch expectedEntry.(type) {
			case entries.FileEntry:
				expectedEntryName = expectedEntry.(entries.FileEntry).Name
			case entries.LinkEntry:
				expectedEntryName = expectedEntry.(entries.LinkEntry).Name
			case entries.DirectoryEntry:
				expectedEntryName = expectedEntry.(entries.DirectoryEntry).Name
			default:
				panic("unknown entry type")
			}

			expectedEntryPath := path.Join(currentPath, expectedEntryName)
			match, _ := path.Match(expectedEntryPath, entryPath)
			if match {
				return true
			}
		}

		return false
	}

	for _, existingEntryName := range existingEntryNames {
		existingEntryPath := path.Join(currentPath, existingEntryName)
		if !checkIfPathIsExpected(existingEntryPath) {
			difference = &Difference{
				Path:        existingEntryPath,
				Expectation: "path doesn't exist",
				Real: "path exists",
			}

			return difference, nil
		}
	}

	// Checks entries
	for _, expectedEntry := range expectedDir.Entries {
		var diff *Difference
		var err error

		subdirectoryPath := path.Join(currentPath, expectedDir.Name)

		switch expectedEntry.(type) {
		case entries.FileEntry:
			expectedFileEntry := expectedEntry.(entries.FileEntry)
			diff, err = c.checkFile(subdirectoryPath, expectedFileEntry)
		case entries.LinkEntry:
			expectedLinkEntry := expectedEntry.(entries.LinkEntry)
			diff, err = c.checkLink(subdirectoryPath, expectedLinkEntry)
		case entries.DirectoryEntry:
			expectedDirectoryEntry := expectedEntry.(entries.DirectoryEntry)
			diff, err = c.checkDir(subdirectoryPath, expectedDirectoryEntry)
		default:
			panic("unknown entry type")
		}

		if diff != nil || err != nil {
			return diff, err
		}
	}

	return nil, nil
}

func (c Checker) checkFile(currentPath string, expectedFile entries.FileEntry) (
	difference *Difference, err error) {

	// Checks that a file exists
	filePath := path.Join(currentPath, expectedFile.Name)
	if !c.Fs.IsFile(filePath) {
		difference = &Difference{
			Path:        filePath,
			Expectation: "file exists",
		}

		if c.Fs.IsExist(filePath) {
			difference.Real = "path isn't a file"
		} else {
			difference.Real = "file doesn't exist"
		}

		return difference, nil
	}

	// Checks the file data equality
	realData, err := c.Fs.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if expectedFile.Data == nil {
		return nil, nil
	}

	if !bytes.Equal(realData, expectedFile.Data) {
		difference = &Difference{
			Path:        filePath,
			Expectation: "file data is equal to expected data",
			Real:        "file data isn't equal ot expected data",
		}
		return difference, nil
	}

	// This check passed successfully
	return nil, nil
}

func (c Checker) checkLink(currentPath string, expectedLink entries.LinkEntry) (
	difference *Difference, err error) {

	linkPath := path.Join(currentPath, expectedLink.Name)

	// Checks that a link exists
	if !c.Fs.IsLink(linkPath) {
		difference = &Difference{
			Path:        linkPath,
			Expectation: "link exists",
		}

		if c.Fs.IsExist(linkPath) {
			difference.Real = "path isn't a link"
		} else {
			difference.Real = "link doesn't exist"
		}

		return difference, nil
	}

	// Checks the link destination
	linkDestination, err := c.Fs.Readlink(linkPath)
	if err != nil {
		return nil, err
	}

	match, err := path.Match(linkDestination, expectedLink.Path)
	if err != nil {
		return nil, err
	}

	if !match {
		difference = &Difference{
			Path:        linkPath,
			Expectation: "link points to" + expectedLink.Path,
			Real:        "link points to" + linkDestination,
		}
		return difference, nil
	}

	// This check passed successfully
	return nil, nil
}
