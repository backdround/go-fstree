package maker

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

////////////////////////////////////////////////////////////
// Perform test functions

func performMake(rootPath string, internalEntries ...entries.Entry) error {
	newTree := entries.DirectoryEntry{
		Name:    "./",
		Entries: internalEntries,
	}

	maker := Maker{osfs.OsFS{}}
	return maker.Make(rootPath, newTree)
}

////////////////////////////////////////////////////////////
// Asserts functions

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

func TestEmptyYamlData(t *testing.T) {
	err := performMake("")
	require.Error(t, err)
}

func TestEmptyRootPath(t *testing.T) {
	rootPath, clean := createRoot()
	defer clean()

	err := performMake(rootPath)
	require.NoError(t, err)
}

func TestRootDoesntExist(t *testing.T) {
	rootPath, clean := createRoot()
	clean()

	err := performMake(rootPath)
	require.NoError(t, err)
	requireDirectory(t, rootPath, ".")

	clean()
}

func TestFile(t *testing.T) {
	t.Run("SuccessOnNewFile", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		// Tests
		err := performMake(rootPath,
			entries.FileEntry{
				Name: "file.txt",
				Data: []byte("some data"),
			},
		)

		// Asserts
		require.NoError(t, err)
		requireFile(t, rootPath, "file.txt", "some data")
	})

	t.Run("ErrorOnAnotherFileAlreadyExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		// Creates a file with another data
		existingFilePath := path.Join(rootPath, "file.txt")
		err := os.WriteFile(existingFilePath, []byte("another data"), 0644)
		assertNoError(err)

		// Tests
		err = performMake(rootPath,
			entries.FileEntry{
				Name: "file.txt",
				Data: []byte("some data"),
			},
		)

		// Asserts
		require.Error(t, err)
		require.Contains(t, err.Error(), "already exists")
	})

	t.Run("SkipOnSameFileAlreadyExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		// Creates a file with the same data
		existingFilePath := path.Join(rootPath, "file.txt")
		err := os.WriteFile(existingFilePath, []byte("some data"), 0644)
		assertNoError(err)

		// Tests
		err = performMake(rootPath,
			entries.FileEntry{
				Name: "file.txt",
				Data: []byte("some data"),
			},
		)

		// Asserts
		require.NoError(t, err)
		requireFile(t, rootPath, "file.txt", "some data")
	})
}

func TestLink(t *testing.T) {
	t.Run("SuccessOnNewlink", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		// Tests
		err := performMake(rootPath,
			entries.LinkEntry{
				Name: "link1",
				Path: "./file.txt",
			},
		)

		// Asserts
		require.NoError(t, err)
		requireLink(t, rootPath, "link1", "./file.txt")
	})

	t.Run("ErrorOnFileAlreadyExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		// Creates a file
		existingFilePath := path.Join(rootPath, "link1")
		err := os.WriteFile(existingFilePath, []byte("data"), 0644)
		assertNoError(err)

		// Tests
		err = performMake(rootPath,
			entries.LinkEntry{
				Name: "link1",
				Path: "./file.txt",
			},
		)

		// Asserts
		require.Error(t, err)
		require.Contains(t, err.Error(), "already exists")
	})

	t.Run("ErrorOnAnotherLinkAlreadyExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		// Creates a link with another destination
		existingLinkPath := path.Join(rootPath, "link1")
		err := os.Symlink("./another-file.txt", existingLinkPath)
		assertNoError(err)

		// Tests
		err = performMake(rootPath,
			entries.LinkEntry{
				Name: "link1",
				Path: "./file.txt",
			},
		)

		// Asserts
		require.Error(t, err)
		require.Contains(t, err.Error(), "already exists")
	})

	t.Run("SkipOnSameLinkExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		// Creates a link with the same destination
		existingLinkPath := path.Join(rootPath, "link1")
		err := os.Symlink("./file.txt", existingLinkPath)
		assertNoError(err)

		// Tests
		err = performMake(rootPath,
			entries.LinkEntry{
				Name: "link1",
				Path: "./file.txt",
			},
		)

		// Asserts
		require.NoError(t, err)
	})
}

func TestDirectory(t *testing.T) {
	t.Run("SuccessOnNewDirectory", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		// Tests
		err := performMake(rootPath,
			entries.DirectoryEntry{
				Name:    "new-directory",
				Entries: []entries.Entry{},
			},
		)

		// Asserts
		require.NoError(t, err)
	})

	t.Run("ErrorOnFileAlreadyExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		// Creates a file with the same path
		existingFilePath := path.Join(rootPath, "new-directory")
		err := os.WriteFile(existingFilePath, []byte("data"), 0644)
		assertNoError(err)

		// Tests
		err = performMake(rootPath,
			entries.DirectoryEntry{
				Name:    "new-directory",
				Entries: []entries.Entry{},
			},
		)

		// Asserts
		require.Error(t, err)
		require.Contains(t, err.Error(), "already exists")
	})

	t.Run("SkipOnAlreadyExists", func(t *testing.T) {
		rootPath, clean := createRoot()
		defer clean()

		// Creates a directory with the same path
		existingDirectoryPath := path.Join(rootPath, "new-directory")
		err := os.Mkdir(existingDirectoryPath, 0755)
		assertNoError(err)

		// Tests
		err = performMake(rootPath,
			entries.DirectoryEntry{
				Name:    "new-directory",
				Entries: []entries.Entry{},
			},
		)

		// Asserts
		require.NoError(t, err)
	})
}

func TestSubDirectory(t *testing.T) {
	rootPath, clean := createRoot()
	defer clean()

	// Tests
	err := performMake(rootPath,
		entries.DirectoryEntry{
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
				entries.DirectoryEntry{
					Name:    "more-sub-directory",
					Entries: []entries.Entry{},
				},
			},
		},
	)

	// Asserts
	require.NoError(t, err)
	requireDirectory(t, rootPath, "sub-directory")

	subdirectoryPath := path.Join(rootPath, "sub-directory")
	requireFile(t, subdirectoryPath, "file.txt", "some data")
	requireLink(t, subdirectoryPath, "link1", "./file.txt")
	requireDirectory(t, subdirectoryPath, "more-sub-directory")
}
