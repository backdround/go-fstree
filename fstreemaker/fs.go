package fstreemaker

type FS interface {
	IsExist(path string) bool
	IsFile(path string) bool
	IsLink(path string) bool
	IsDirectory(path string) bool

	ReadFile(path string) ([]byte, error)
	Readlink(path string) (string, error)
	WriteFile(path string, data []byte) error
	Symlink(oldPath, newPath string) error
	Mkdir(path string) error
}
