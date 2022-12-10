package types

type DirectoryEntry struct {
	Name    string
	Entries []any
}

type FileEntry struct {
	Name string
	Data string
}

type LinkEntry struct {
	Name string
	Path string
}
