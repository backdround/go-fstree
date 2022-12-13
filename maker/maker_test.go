package maker

import (
	"os"
	"path"
	"testing"

	"github.com/backdround/go-fstree/osfs"
	"github.com/stretchr/testify/require"

	"github.com/backdround/go-fstree/entries"
)

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func createRoot() (rootPath string, clean func()) {
	rootPath, err := os.MkdirTemp("", "go-fstreemaker-test-*.d")
	assertNoError(err)

	clean = func() {
		err := os.RemoveAll(rootPath)
		assertNoError(err)
	}

	return
}

func GetOsMaker() Maker {
	maker := Maker{
		Fs: osfs.OsFS{},
	}
	return maker
}

func TestFileCreationCornerCases(t *testing.T) {
	t.Run("ErrorOnAnotherFileAlreadyExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		rootEntry := entries.DirectoryEntry{
			Name: rootPath,
			Entries: []any{
				entries.FileEntry{
					Name: "file.txt",
					Data: []byte("some data"),
				},
			},
		}

		// Creates a file with another data
		existingFilePath := path.Join(rootPath, "file.txt")
		err := os.WriteFile(existingFilePath, []byte("another data"), 0644)
		assertNoError(err)

		// Tests
		err = GetOsMaker().Make(rootEntry)
		require.Error(t, err)
		require.Contains(t, err.Error(), "already exists")
	})

	t.Run("SkipOnSameFileAlreadyExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		rootEntry := entries.DirectoryEntry{
			Name: rootPath,
			Entries: []any{
				entries.FileEntry{
					Name: "file.txt",
					Data: []byte("some data"),
				},
			},
		}

		// Creates a file with the same data
		existingFilePath := path.Join(rootPath, "file.txt")
		err := os.WriteFile(existingFilePath, []byte("some data"), 0644)
		assertNoError(err)

		// Tests
		err = GetOsMaker().Make(rootEntry)
		require.NoError(t, err)
	})
}

func TestLinkCreationCornerCases(t *testing.T) {
	t.Run("ErrorOnFileAlreadyExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		rootEntry := entries.DirectoryEntry{
			Name: rootPath,
			Entries: []any{
				entries.LinkEntry{
					Name: "link1",
					Path: "./file.txt",
				},
			},
		}

		// Creates a file with the same path
		existingFilePath := path.Join(rootPath, "link1")
		err := os.WriteFile(existingFilePath, []byte("data"), 0644)
		assertNoError(err)

		// Tests
		err = GetOsMaker().Make(rootEntry)
		require.Error(t, err)
		require.Contains(t, err.Error(), "already exists")
	})

	t.Run("ErrorOnAnotherLinkAlreadyExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		rootEntry := entries.DirectoryEntry{
			Name: rootPath,
			Entries: []any{
				entries.LinkEntry{
					Name: "link1",
					Path: "./file.txt",
				},
			},
		}

		// Creates a link with another destination
		existingLinkPath := path.Join(rootPath, "link1")
		err := os.Symlink("./another-file.txt", existingLinkPath)
		assertNoError(err)

		// Tests
		err = GetOsMaker().Make(rootEntry)
		require.Error(t, err)
		require.Contains(t, err.Error(), "already exists")
	})

	t.Run("SkipOnSameLinkExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		rootEntry := entries.DirectoryEntry{
			Name: rootPath,
			Entries: []any{
				entries.LinkEntry{
					Name: "link1",
					Path: "./file.txt",
				},
			},
		}

		// Creates a link with the same destination
		existingLinkPath := path.Join(rootPath, "link1")
		err := os.Symlink("./file.txt", existingLinkPath)
		assertNoError(err)

		// Tests
		err = GetOsMaker().Make(rootEntry)
		require.NoError(t, err)
	})
}

func TestDirectoryCreationCornerCases(t *testing.T) {
	t.Run("ErrorOnFileAlreadyExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		rootEntry := entries.DirectoryEntry{
			Name: rootPath,
			Entries: []any{
				entries.DirectoryEntry{
					Name:    "new-directory",
					Entries: []any{},
				},
			},
		}

		// Creates a file with the same path
		existingFilePath := path.Join(rootPath, "new-directory")
		err := os.WriteFile(existingFilePath, []byte("data"), 0644)
		assertNoError(err)

		// Tests
		err = GetOsMaker().Make(rootEntry)
		require.Error(t, err)
		require.Contains(t, err.Error(), "already exists")
	})

	t.Run("SkipOnAlreadyExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		rootEntry := entries.DirectoryEntry{
			Name: rootPath,
			Entries: []any{
				entries.DirectoryEntry{
					Name:    "new-directory",
					Entries: []any{},
				},
			},
		}

		// Creates a directory with the same path
		existingDirectoryPath := path.Join(rootPath, "new-directory")
		err := os.Mkdir(existingDirectoryPath, 0755)
		assertNoError(err)

		// Tests
		err = GetOsMaker().Make(rootEntry)
		require.NoError(t, err)
	})
}
