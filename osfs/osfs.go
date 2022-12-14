// Package osfs describes OsFS type that implements work with
// filesystem by os package.
package osfs

import "os"

type OsFS struct{}

func (OsFS) IsExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func (OsFS) IsFile(path string) bool {
	fileInfo, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false
	}

	return fileInfo.Mode().IsRegular()
}

func (OsFS) IsLink(path string) bool {
	pathInfo, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false
	}

	return (pathInfo.Mode() & os.ModeSymlink) == os.ModeSymlink
}

func (OsFS) IsDirectory(path string) bool {
	fileInfo, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false
	}

	return fileInfo.IsDir()
}

func (OsFS) ReadDir(path string) (entryPaths []string, err error) {
	dirEntries, err := os.ReadDir(path)
	entryPaths = []string{}

	for _, entry := range dirEntries {
		entryPaths = append(entryPaths, entry.Name())
	}

	return
}

func (OsFS) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (OsFS) Readlink(path string) (string, error) {
	return os.Readlink(path)
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
