package manajemen_file

import (
	"fmt"
	"os"
)

// Mkdir creates a new directory.
// If createParents is true, it will create parent directories as needed.
func Mkdir(path string, createParents bool) error {
	if path == "" {
		return fmt.Errorf("mkdir: nama direktori tidak boleh kosong")
	}

	if createParents {
		err := os.MkdirAll(path, 0755) // Using 0755 for default permissions
		if err != nil {
			return fmt.Errorf("gagal membuat direktori %s: %w", path, err)
		}
	} else {
		err := os.Mkdir(path, 0755) // Using 0755 for default permissions
		if err != nil {
			return fmt.Errorf("gagal membuat direktori %s: %w", path, err)
		}
	}

	return nil
}
