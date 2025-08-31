package manajemen_file

import (
	"fmt"
	"os"
)

// Rm removes files or directories.
func Rm(path string, recursive, force bool) error {
	if path == "" {
		return fmt.Errorf("rm: nama file atau direktori tidak boleh kosong")
	}

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		if force {
			return nil // With -f, ignore non-existent files
		}
		return fmt.Errorf("rm: tidak dapat menghapus '%s': File atau direktori tidak ada", path)
	}

	if info.IsDir() && !recursive {
		return fmt.Errorf("rm: tidak dapat menghapus '%s': Adalah sebuah direktori", path)
	}

	// The actual removal logic
	if recursive {
		err := os.RemoveAll(path)
		if err != nil {
			return fmt.Errorf("gagal menghapus '%s': %w", path, err)
		}
	} else {
		err := os.Remove(path)
		if err != nil {
			return fmt.Errorf("gagal menghapus '%s': %w", path, err)
		}
	}

	return nil
}
