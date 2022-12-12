package types

type FS interface {
	IsExist(path string) bool
	IsFile(path string) bool
	IsDirectory(path string) bool

	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	Symlink(oldPath, newPath string) error
	Mkdir(path string) error
}
