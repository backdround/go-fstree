package fstree_test

import (
	"os"
	"path"
	"testing"

	"github.com/backdround/go-fstree"
	"github.com/stretchr/testify/require"
)

////////////////////////////////////////////////////////////
// Utility functions

func createFile(basePath string, name string, data string) string {
	filePath := path.Join(basePath, name)
	err := os.WriteFile(filePath, []byte(data), 0644)
	assertNoError(err)

	return filePath
}

func createLink(basePath string, name string, destination string) string {
	linkPath := path.Join(basePath, name)
	err := os.Symlink(destination, linkPath)
	assertNoError(err)

	return linkPath
}

func createDirectory(basePath string, name string) string {
	directoryPath := path.Join(basePath, name)
	err := os.Mkdir(directoryPath, 0755)
	assertNoError(err)

	return directoryPath
}

////////////////////////////////////////////////////////////
// Tests

func TestCheck(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		root, clean := createRoot()
		defer clean()

		yamlData := prepareYaml(`
			new-directory:
				file.txt:
					type: file
				link1:
					type: link
					path: ./file.txt
				subdirectory:
		`)

		// Creates filetree to check
		newDirectoryPath := createDirectory(root, "new-directory")
		createFile(newDirectoryPath, "file.txt", "")
		createLink(newDirectoryPath, "link1", "./file.txt")
		createDirectory(newDirectoryPath, "subdirectory")

		// Tests
		difference, err := fstree.CheckOverOSFS(root, yamlData)
		require.NoError(t, err)
		require.Nil(t, difference)
	})

	t.Run("Fail", func(t *testing.T) {
		root, clean := createRoot()
		defer clean()

		yamlData := prepareYaml(`
			new-directory:
				# This file will be missed
				file.txt:
					type: file
				link1:
					type: link
					path: ./file.txt
				subdirectory:
		`)

		// Creates filetree to check
		newDirectoryPath := createDirectory(root, "new-directory")
		createLink(newDirectoryPath, "link1", "./file.txt")
		createDirectory(newDirectoryPath, "subdirectory")

		// Tests
		difference, err := fstree.CheckOverOSFS(root, yamlData)
		require.NoError(t, err)
		require.NotNil(t, difference)

		missedFilePath := path.Join(newDirectoryPath, "file.txt")
		require.Equal(t, missedFilePath, difference.Path)
		require.Contains(t, difference.Real, "file doesn't exist")
	})
}
