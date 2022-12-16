package fstree_test

import (
	"os"
	"path"
	"testing"

	"github.com/backdround/go-fstree"
	"github.com/stretchr/testify/require"
)

////////////////////////////////////////////////////////////
// Assertion functions

func requireFile(t *testing.T, basePath string, name string, data string) {
	t.Helper()

	filePath := path.Join(basePath, name)
	realData, err := os.ReadFile(filePath)
	require.NoError(t, err)
	require.Equal(t, data, string(realData))
}

func requireLink(t *testing.T, basePath string, name string,
	destination string) {
	t.Helper()

	linkPath := path.Join(basePath, name)
	realDestination, err := os.Readlink(linkPath)
	require.NoError(t, err)

	match, err := path.Match(destination, realDestination)
	require.NoError(t, err)
	require.True(t, match)
}

func requireDirectory(t *testing.T, basePath string, name string) {
	t.Helper()

	directoryPath := path.Join(basePath, name)
	info, err := os.Lstat(directoryPath)
	require.NoError(t, err)
	require.True(t, info.IsDir())
}

////////////////////////////////////////////////////////////
// Tests

func TestMake(t *testing.T) {
	t.Run("ErrorOnEmptyRoot", func(t *testing.T) {
		yamlData := `file.txt: {type: "file", data: "data"}`

		err := fstree.MakeOverOSFS("", yamlData)
		require.Error(t, err)
	})

	t.Run("CreateRootOnNonexistence", func(t *testing.T) {
		root, clean := createRoot()
		clean()

		yamlData := `file.txt: {type: "file", data: "data"}`
		err := fstree.MakeOverOSFS(root, yamlData)

		require.NoError(t, err)
		requireDirectory(t, root, ".")

		clean()
	})

	t.Run("Complex", func(t *testing.T) {
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

		// Tests
		err := fstree.MakeOverOSFS(root, yamlData)
		require.NoError(t, err)

		// Checks fs tree
		subdirectoryPath := path.Join(root, "new-directory")
		requireFile(t, subdirectoryPath, "file.txt", "")
		requireLink(t, subdirectoryPath, "link1", "./file.txt")
		requireDirectory(t, subdirectoryPath, "subdirectory")
	})

	t.Run("Idempotency", func(t *testing.T) {
		root, clean := createRoot()
		defer clean()

		yamlData := prepareYaml(`
			new-directory:
				file.txt:
					type: file
				link:
					type: link
					path: ./file.txt
				subdirectory:
		`)

		err := fstree.MakeOverOSFS(root, yamlData)
		require.NoError(t, err)
		err = fstree.MakeOverOSFS(root, yamlData)
		require.NoError(t, err)
	})
}
