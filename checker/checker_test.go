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

func createFile(basePath string, name string, data string) {
	filePath := path.Join(basePath, name)
	err := os.WriteFile(filePath, []byte(data), 0644)
	assertNoError(err)
}

func createLink(basePath string, name string, destination string) {
	linkPath := path.Join(basePath, name)
	err := os.Symlink(destination, linkPath)
	assertNoError(err)
}

func createDirectory(basePath string, name string) {
	directoryPath := path.Join(basePath, name)
	err := os.Mkdir(directoryPath, 0755)
	assertNoError(err)
}

////////////////////////////////////////////////////////////
// Perform test functions

func performCheck(rootPath string, internalEntries ...any) (difference string,
	err error) {

	expectedTree := entries.DirectoryEntry{
		Name: "./",
		Entries: internalEntries,
	}

	checker := Checker{osfs.OsFS{}}
	return checker.Check(rootPath, expectedTree)
}

////////////////////////////////////////////////////////////
// Asserts functions

func requireDifferent(t *testing.T, difference string, err error) {
	t.Helper()
	require.NoError(t, err)
	require.NotEmpty(t, difference)
}

func requireTheSame(t *testing.T, difference string, err error) {
	t.Helper()
	require.NoError(t, err)
	require.Empty(t, difference)
}

////////////////////////////////////////////////////////////
// Test cases

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
		createFile(rootPath, "file.txt", "another data")

		difference, err := performCheck(rootPath, entries.FileEntry{
			Name: "file.txt",
			Data: []byte("some data"),
		})

		requireDifferent(t, difference, err)
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
		createFile(rootPath, "file.txt", "some data")

		difference, err := performCheck(rootPath, entries.FileEntry{
			Name: "file.txt",
			Data: []byte{},
		})

		requireDifferent(t, difference, err)
	})


	t.Run("DoesntExist", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		difference, err := performCheck(rootPath, entries.FileEntry{
			Name: "file.txt",
			Data: []byte("some data"),
		})

		requireDifferent(t, difference, err)
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
		createLink(rootPath, "link1", "./another-file.txt")

		difference, err := performCheck(rootPath, entries.LinkEntry{
			Name: "link1",
			Path: "./file.txt",
		})

		requireDifferent(t, difference, err)
	})

	t.Run("DoesntExist", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		difference, err := performCheck(rootPath, entries.LinkEntry{
			Name: "link1",
			Path: "./file.txt",
		})

		requireDifferent(t, difference, err)
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
			Entries: []any{
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
		createFile(subdirectoryPath, "file.txt", "another data")

		difference, err := performCheck(rootPath, entries.DirectoryEntry{
			Name: "sub-directory",
			Entries: []any{
				entries.FileEntry{
					Name: "file.txt",
					Data: []byte("some data"),
				},
			},
		})

		requireDifferent(t, difference, err)
	})
}
