package manajemen_file

import (
	"fmt"
	"os"
	
)

var oldPwd string // To store the previous working directory

// Cd changes the current working directory of the program.
func Cd(path string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("gagal mendapatkan direktori saat ini: %w", err)
	}
	oldPwd = currentDir // Store current directory as oldPwd

	if path == "" || path == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("gagal mendapatkan direktori home: %w", err)
		}
		path = homeDir
	} else if path == "-" {
		if oldPwd == "" {
			return fmt.Errorf("tidak ada direktori sebelumnya")
		}
		path = oldPwd
	}

	err = os.Chdir(path)
	if err != nil {
		return fmt.Errorf("gagal mengubah direktori ke %s: %w", path, err)
	}
	return nil
}