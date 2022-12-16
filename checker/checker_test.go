package checker

import (
	"os"
	"path"
	"testing"

	"github.com/backdround/go-fstree/osfs"
	"github.com/stretchr/testify/require"

	"github.com/backdround/go-fstree/entries"
)

////////////////////////////////////////////////////////////
// Utility functions

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

////////////////////////////////////////////////////////////
// Checked filetree preparating functions

func createRoot() (rootPath string, clean func()) {
	rootPath, err := os.MkdirTemp("", "go-fstree-checker-test-*.d")
	assertNoError(err)

	clean = func() {
		err := os.RemoveAll(rootPath)
		assertNoError(err)
	}

	return
}

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
// Perform test functions

func performCheck(rootPath string, internalEntries ...entries.Entry) (
	difference *Difference, err error) {

	expectedTree := entries.DirectoryEntry{
		Name:    "./",
		Entries: internalEntries,
	}

	checker := Checker{osfs.OsFS{}}
	return checker.Check(rootPath, expectedTree)
}

////////////////////////////////////////////////////////////
// Asserts functions

func requireDifferent(t *testing.T, difference *Difference, err error) {
	t.Helper()
	require.NoError(t, err)
	require.NotNil(t, difference)
}

func requireDifferentPath(t *testing.T, expectedPath string,
	differencePath string) {
	t.Helper()
	match, err := path.Match(expectedPath, differencePath)
	require.NoError(t, err, "Unable to match the given paths:",
		expectedPath, differencePath)
	require.True(t, match, "Expect that paths be the same",
		expectedPath, differencePath)
}

func requireTheSame(t *testing.T, difference *Difference, err error) {
	t.Helper()
	require.NoError(t, err)
	require.Nil(t, difference)
}

////////////////////////////////////////////////////////////
// Test cases

func TestUnexpectedEntry(t *testing.T) {
	rootPath, clean := createRoot()
	defer clean()
	filePath := createFile(rootPath, "file.txt", "")

	difference, err := performCheck(rootPath)

	requireDifferent(t, difference, err)
	requireDifferentPath(t, filePath, difference.Path)
}

func TestFile(t *testing.T) {
	t.Run("WithTheSameDataExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()
		createFile(rootPath, "file.txt", "some data")

		difference, err := performCheck(rootPath, entries.FileEntry{
			Name: "file.txt",
			Data: []byte("some data"),
		})

		requireTheSame(t, difference, err)
	})

	t.Run("WithAnotherDataExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()
		filePath := createFile(rootPath, "file.txt", "another data")

		difference, err := performCheck(rootPath, entries.FileEntry{
			Name: "file.txt",
			Data: []byte("some data"),
		})

		requireDifferent(t, difference, err)
		requireDifferentPath(t, filePath, difference.Path)
	})

	t.Run("DontCheckData", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()
		createFile(rootPath, "file.txt", "some data")

		difference, err := performCheck(rootPath, entries.FileEntry{
			Name: "file.txt",
		})

		requireTheSame(t, difference, err)
	})

	t.Run("CheckEmptyData", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()
		filePath := createFile(rootPath, "file.txt", "some data")

		difference, err := performCheck(rootPath, entries.FileEntry{
			Name: "file.txt",
			Data: []byte{},
		})

		requireDifferent(t, difference, err)
		requireDifferentPath(t, filePath, difference.Path)
	})

	t.Run("DoesntExist", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		difference, err := performCheck(rootPath, entries.FileEntry{
			Name: "file.txt",
			Data: []byte("some data"),
		})

		requireDifferent(t, difference, err)
		expectedDifferencePath := path.Join(rootPath, "file.txt")
		requireDifferentPath(t, expectedDifferencePath, difference.Path)
	})
}

func TestLink(t *testing.T) {
	t.Run("WithTheSamePathExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()
		createLink(rootPath, "link1", "./file.txt")

		difference, err := performCheck(rootPath, entries.LinkEntry{
			Name: "link1",
			Path: "./file.txt",
		})

		requireTheSame(t, difference, err)
	})

	t.Run("WithAnotherPathExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()
		linkPath := createLink(rootPath, "link1", "./another-file.txt")

		difference, err := performCheck(rootPath, entries.LinkEntry{
			Name: "link1",
			Path: "./file.txt",
		})

		requireDifferent(t, difference, err)
		requireDifferentPath(t, linkPath, difference.Path)
	})

	t.Run("DoesntExist", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		difference, err := performCheck(rootPath, entries.LinkEntry{
			Name: "link1",
			Path: "./file.txt",
		})

		requireDifferent(t, difference, err)
		expectedDifferencePath := path.Join(rootPath, "link1")
		requireDifferentPath(t, expectedDifferencePath, difference.Path)
	})
}

func TestDirectory(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()
		createDirectory(rootPath, "some-directory")

		difference, err := performCheck(rootPath, entries.DirectoryEntry{
			Name: "some-directory",
		})

		requireTheSame(t, difference, err)
	})

	t.Run("DoesntExist", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		difference, err := performCheck(rootPath, entries.DirectoryEntry{
			Name: "some-directory",
		})

		requireDifferent(t, difference, err)
		expectedDifferencePath := path.Join(rootPath, "some-directory")
		requireDifferentPath(t, expectedDifferencePath, difference.Path)
	})
}

func TestSubdirectory(t *testing.T) {
	t.Run("Same", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		createDirectory(rootPath, "sub-directory")
		subdirectoryPath := path.Join(rootPath, "sub-directory")
		createFile(subdirectoryPath, "file.txt", "some data")
		createLink(subdirectoryPath, "link1", "./file.txt")

		difference, err := performCheck(rootPath, entries.DirectoryEntry{
			Name: "sub-directory",
			Entries: []entries.Entry{
				entries.FileEntry{
					Name: "file.txt",
					Data: []byte("some data"),
				},
				entries.LinkEntry{
					Name: "link1",
					Path: "./file.txt",
				},
			},
		})

		requireTheSame(t, difference, err)
	})

	t.Run("Another", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		createDirectory(rootPath, "sub-directory")
		subdirectoryPath := path.Join(rootPath, "sub-directory")
		filePath := createFile(subdirectoryPath, "file.txt", "another data")

		difference, err := performCheck(rootPath, entries.DirectoryEntry{
			Name: "sub-directory",
			Entries: []entries.Entry{
				entries.FileEntry{
					Name: "file.txt",
					Data: []byte("some data"),
				},
			},
		})

		requireDifferent(t, difference, err)
		requireDifferentPath(t, filePath, difference.Path)
	})
}
