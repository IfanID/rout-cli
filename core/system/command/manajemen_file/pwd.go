package manajemen_file

import (
	"fmt"
	"os"
)

// Pwd prints the current working directory.
func Pwd() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("gagal mendapatkan direktori kerja saat ini: %w", err)
	}
	return dir, nil
}