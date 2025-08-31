package manajemen_file

import (
	"fmt"
	"os"
	"path/filepath"
	"rout/core/system/util"
)

// Mv moves or renames a file or directory from source to destination.
func Mv(src, dst string, force, interactive bool) error {
	// Check if destination is a directory and adjust path
	destInfo, err := os.Stat(dst)
	if err == nil && destInfo.IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}

	// Check if the final destination already exists
	_, err = os.Stat(dst)
	if err == nil {
		// Destination exists
		if force {
			// Overwrite without asking
		} else if interactive {
			prompt := fmt.Sprintf("mv: timpa '%s'?", dst)
			confirmed, err := util.Confirm(prompt)
			if err != nil {
				return err
			}
			if !confirmed {
				return nil // User said no, so we do nothing and return
			}
		} else {
			return fmt.Errorf("mv: tidak dapat memindahkan '%s' ke '%s': File sudah ada", src, dst)
		}
	}

	// Perform the rename/move
	err = os.Rename(src, dst)
	if err != nil {
		// A more complete implementation would copy-then-delete for cross-device moves.
		return fmt.Errorf("mv: %w", err)
	}
	return nil
}