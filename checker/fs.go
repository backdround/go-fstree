package checker

type FS interface {
	IsExist(path string) bool
	IsFile(path string) bool
	IsLink(path string) bool
	IsDirectory(path string) bool

	Abs(path string) (string, error)
	ReadDir(path string) ([]string, error)
	ReadFile(path string) ([]byte, error)
	Readlink(path string) (string, error)
}
