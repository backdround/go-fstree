package config

import (
	"fmt"
	"strings"
	"testing"

	"github.com/backdround/go-fstree/v2/entries"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
)

func prepareYaml(data string) string {
	data = dedent.Dedent(data)
	data = strings.ReplaceAll(data, "\t", "  ")
	return data
}

func TestInvalidTypeField(t *testing.T) {
	yamlPattern := `
		directory:
			type: %v
	`
	yamlPattern = prepareYaml(yamlPattern)

	testCases := []struct {
		Name     string
		TypeData string
	}{
		{"InvalidTypeScalar1", "3"},
		{"InvalidTypeScalar2", "UnknownValue"},
		{"InvalidTypeDictionary", "{a: 3}"},
		{"InvalidTypeList", `["type"]`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			yaml := fmt.Sprintf(yamlPattern, testCase.TypeData)
			_, err := Parse(yaml)
			require.Error(t, err)
			require.Contains(t, err.Error(), "directory")
		})
	}
}

func TestInvalidDataStructure(t *testing.T) {
	yaml := `
		directory:
			- list
			- "is't"
			- permitted
	`
	yaml = prepareYaml(yaml)

	_, err := Parse(yaml)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unable to convert to dictionary")
}

func TestInvalidYaml(t *testing.T) {
	_, err := Parse("\t")
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot start any token")
}

func TestRootEntries(t *testing.T) {
	t.Run("File", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			yaml := `
				file.txt:
					type: file
			`
			yaml = prepareYaml(yaml)

			rootEntry, err := Parse(yaml)
			require.NoError(t, err)

			// Asserts entires
			require.Len(t, rootEntry.Entries, 1)
			require.IsType(t, entries.FileEntry{}, rootEntry.Entries[0])

			// Asserts file
			file := rootEntry.Entries[0].(entries.FileEntry)
			require.Equal(t, "file.txt", file.Name)
			require.Nil(t, file.Data)
		})

		t.Run("WithData", func(t *testing.T) {
			yaml := `
				file.txt:
					type: file
					data: some data
			`
			yaml = prepareYaml(yaml)

			rootEntry, err := Parse(yaml)
			require.NoError(t, err)

			// Asserts entires
			require.Len(t, rootEntry.Entries, 1)
			require.IsType(t, entries.FileEntry{}, rootEntry.Entries[0])

			// Asserts file
			file := rootEntry.Entries[0].(entries.FileEntry)
			require.Equal(t, "file.txt", file.Name)
			require.Equal(t, []byte("some data"), file.Data)
		})

		t.Run("ErrorInvalidDataType", func(t *testing.T) {
			yaml := `
				file.txt:
					type: file
					data:
						a: b
			`
			yaml = prepareYaml(yaml)

			_, err := Parse(yaml)
			require.Error(t, err)
			require.Contains(t, err.Error(), "file.txt")
		})

		t.Run("ErrorUnknownField", func(t *testing.T) {
			yaml := `
				file.txt:
					type: file
					path: "../../"
			`
			yaml = prepareYaml(yaml)

			_, err := Parse(yaml)
			require.Error(t, err)
			require.Contains(t, err.Error(), "file.txt")
		})
	})

	t.Run("Link", func(t *testing.T) {
		t.Run("Valid", func(t *testing.T) {
			yaml := `
				pkg1:
					type: link
					path: "../../pkg1"
			`
			yaml = prepareYaml(yaml)

			rootEntry, err := Parse(yaml)
			require.NoError(t, err)

			// Asserts enties
			require.Len(t, rootEntry.Entries, 1)
			require.IsType(t, entries.LinkEntry{}, rootEntry.Entries[0])

			// Asserts file
			link := rootEntry.Entries[0].(entries.LinkEntry)
			require.Equal(t, "pkg1", link.Name)
			require.Equal(t, "../../pkg1", link.Path)
		})

		t.Run("ErrorMissingPathValue", func(t *testing.T) {
			yaml := `
				link1:
					type: link
			`
			yaml = prepareYaml(yaml)

			_, err := Parse(yaml)
			require.Error(t, err)
			require.Contains(t, err.Error(), "link1")
		})

		t.Run("ErrorInvalidDataType", func(t *testing.T) {
			yaml := `
				link1:
					type: link
					path:
						a: b
			`
			yaml = prepareYaml(yaml)

			_, err := Parse(yaml)
			require.Error(t, err)
			require.Contains(t, err.Error(), "link1")
		})

		t.Run("ErrorUnknownField", func(t *testing.T) {
			yaml := `
				link1:
					type: link
					path: "../../"
					data: "some data"
			`
			yaml = prepareYaml(yaml)

			_, err := Parse(yaml)
			require.Error(t, err)
			require.Contains(t, err.Error(), "unknown property")
		})
	})

	t.Run("Directory", func(t *testing.T) {
		yaml := `
			"new-directory":
		`
		yaml = prepareYaml(yaml)

		rootEntry, err := Parse(yaml)
		require.NoError(t, err)

		// Asserts enties
		require.Len(t, rootEntry.Entries, 1)
		require.IsType(t, entries.DirectoryEntry{}, rootEntry.Entries[0])

		// Asserts file
		newDirectory := rootEntry.Entries[0].(entries.DirectoryEntry)
		require.Equal(t, "new-directory", newDirectory.Name)
		require.Len(t, newDirectory.Entries, 0)
	})

	t.Run("SeveralEntries", func(t *testing.T) {
		yaml := `
			file.txt:
				type: file
				data: "some data"
			new-directory:
		`
		yaml = prepareYaml(yaml)

		rootEntry, err := Parse(yaml)
		require.NoError(t, err)
		require.Len(t, rootEntry.Entries, 2)

		checkFile := func(entry interface{}) {
			file := entry.(entries.FileEntry)
			require.Equal(t, "file.txt", file.Name)
			require.Equal(t, []byte("some data"), file.Data)
		}

		checkDirectory := func(entry interface{}) {
			directory := entry.(entries.DirectoryEntry)
			require.Equal(t, "new-directory", directory.Name)
		}

		if _, ok := rootEntry.Entries[0].(entries.FileEntry); ok {
			checkFile(rootEntry.Entries[0])
			checkDirectory(rootEntry.Entries[1])
		} else {
			checkDirectory(rootEntry.Entries[0])
			checkFile(rootEntry.Entries[1])
		}
	})
}

func TestParseInSubDirectory(t *testing.T) {
	t.Run("File", func(t *testing.T) {
		yaml := `
			new-directory:
				file.txt:
					type: file
					data: "some data"
		`
		yaml = prepareYaml(yaml)

		rootEntry, err := Parse(yaml)
		require.NoError(t, err)

		// Checks directory
		require.Len(t, rootEntry.Entries, 1)
		require.IsType(t, entries.DirectoryEntry{}, rootEntry.Entries[0])
		newDirectory := rootEntry.Entries[0].(entries.DirectoryEntry)
		require.Equal(t, "new-directory", newDirectory.Name)
		require.Len(t, newDirectory.Entries, 1)

		// Checks file inside new directory
		require.IsType(t, entries.FileEntry{}, newDirectory.Entries[0])
		file := newDirectory.Entries[0].(entries.FileEntry)
		require.Equal(t, "file.txt", file.Name)
		require.Equal(t, []byte("some data"), file.Data)
	})

	t.Run("InvalidLink", func(t *testing.T) {
		yaml := `
			new-directory:
				file.txt:
					type: link
		`
		yaml = prepareYaml(yaml)

		_, err := Parse(yaml)
		require.Error(t, err)
		require.Contains(t, err.Error(), "new-directory/file.txt")
	})
}
