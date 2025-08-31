package system

import (
	"io/ioutil"
	"os"
)

// CreateDirectory creates a new directory at the specified path.
// It's equivalent to 'mkdir -p'.
func CreateDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

// ReadDirectory lists the files and subdirectories within a given path.
func ReadDirectory(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(path)
}
