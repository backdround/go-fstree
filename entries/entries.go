package entries

type EntryType = int

type Entry interface {
	GetName() string
}

type DirectoryEntry struct {
	Name    string
	Entries []Entry
}

func (e DirectoryEntry) GetName() string {
	return e.Name
}

type FileEntry struct {
	Name string
	Data []byte
}

func (e FileEntry) GetName() string {
	return e.Name
}

type LinkEntry struct {
	Name string
	Path string
}

func (e LinkEntry) GetName() string {
	return e.Name
}
