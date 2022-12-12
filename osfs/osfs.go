package osfs

import "os"

type OsFS struct{}

func (OsFS) IsExist(path string) bool {
	_, err := os.Lstat(path)
	return os.IsNotExist(err)
}

func (OsFS) IsFile(path string) bool {
	fileInfo, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false
	}

	return fileInfo.Mode().IsRegular()
}

func (OsFS) IsDirectory(path string) bool {
	fileInfo, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false
	}

	return fileInfo.IsDir()
}

func (OsFS) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (OsFS) WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func (OsFS) Symlink(oldPath, newPath string) error {
	return os.Symlink(oldPath, newPath)
}

func (OsFS) Mkdir(path string) error {
	return os.Mkdir(path, 0755)
}
