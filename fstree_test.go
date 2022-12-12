package fstree_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"

	"github.com/backdround/go-fstree"
)

////////////////////////////////////////////////////////////
// Utility functions

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func prepareYaml(data string) string {
	data = dedent.Dedent(data)
	data = strings.ReplaceAll(data, "\t", "  ")
	return data
}

func createRoot() (rootPath string, clean func()) {
	rootPath, err := os.MkdirTemp("", "go-fstree-test-*.d")
	assertNoError(err)

	clean = func() {
		err := os.RemoveAll(rootPath)
		assertNoError(err)
	}

	return
}

////////////////////////////////////////////////////////////
// Tests

func TestEmptyRoot(t *testing.T) {
	yamlData :=`file.txt: {type: "file", data: "data"}`

	err := fstree.MakeOverOSFS("", yamlData)
	require.Error(t, err)
}

func TestFileCreation(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		yamlData :=`
			file.txt:
				type: file
		`
		yamlData = prepareYaml(yamlData)

		root, clean := createRoot()
		defer clean()

		err := fstree.MakeOverOSFS(root, yamlData)
		require.NoError(t, err)

		filePath := path.Join(root, "file.txt")
		data, err := os.ReadFile(filePath)
		require.NoError(t, err)
		require.Empty(t, data)
	})

	t.Run("WithData", func(t *testing.T) {
		yamlData :=`
			file.txt:
				type: file
				data: some data
		`
		yamlData = prepareYaml(yamlData)

		root, clean := createRoot()
		defer clean()

		err := fstree.MakeOverOSFS(root, yamlData)
		require.NoError(t, err)

		filePath := path.Join(root, "file.txt")
		data, err := os.ReadFile(filePath)
		require.NoError(t, err)
		require.Equal(t, "some data", string(data))
	})

	t.Run("InvalidData", func(t *testing.T) {
		yamlData :=`
			file.txt:
				type: file
				data:
					var: dictionary value
		`
		yamlData = prepareYaml(yamlData)

		root, clean := createRoot()
		defer clean()

		err := fstree.MakeOverOSFS(root, yamlData)
		require.Error(t, err)
	})
}

func TestLinkCreation(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		yamlData :=`
			link:
				type: link
				path: "./destination"
		`
		yamlData = prepareYaml(yamlData)

		root, clean := createRoot()
		defer clean()

		err := fstree.MakeOverOSFS(root, yamlData)
		require.NoError(t, err)

		linkPath := path.Join(root, "link")
		linkDestination, err := os.Readlink(linkPath)
		require.NoError(t, err)
		require.Equal(t, "./destination", linkDestination)
	})

	t.Run("Invalid", func(t *testing.T) {
		yamlData :=`
			link:
				type: link
		`
		yamlData = prepareYaml(yamlData)

		root, clean := createRoot()
		defer clean()

		err := fstree.MakeOverOSFS(root, yamlData)
		require.Error(t, err)
	})
}

func TestDirectoryCreation(t *testing.T) {
	yamlData :=`directory:`

	root, clean := createRoot()
	defer clean()

	err := fstree.MakeOverOSFS(root, yamlData)
	require.NoError(t, err)

	directoryPath := path.Join(root, "directory")
	pathInfo, err := os.Lstat(directoryPath)
	require.NoError(t, err)
	require.True(t, pathInfo.IsDir())
}

func TestSubdirectory(t *testing.T) {
	yamlData :=`
		new-directory:
			file.txt:
				type: file
			link:
				type: link
				path: ./file.txt
			subdirectory:
	`
	yamlData = prepareYaml(yamlData)

	root, clean := createRoot()
	defer clean()

	err := fstree.MakeOverOSFS(root, yamlData)
	require.NoError(t, err)

	// Checks file
	filePath := path.Join(root, "new-directory", "file.txt")
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	require.Empty(t, data)

	// Checks link
	linkPath := path.Join(root, "new-directory", "link")
	linkDestination, err := os.Readlink(linkPath)
	require.NoError(t, err)
	require.Equal(t, "./file.txt", linkDestination)

	// Checks subdirectory
	subdirectoryPath := path.Join(root, "new-directory", "subdirectory")
	directoryInfo, err := os.Lstat(subdirectoryPath)
	require.NoError(t, err)
	require.True(t, directoryInfo.IsDir())
}
