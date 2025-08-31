package manajemen_file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"rout/core/system/util"
)

// Cp copies a file or directory from source to destination.
func Cp(src, dst string, recursive, force, interactive bool) error {
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("cp: %w", err)
	}

	if sourceInfo.IsDir() {
		if !recursive {
			return fmt.Errorf("cp: -r tidak dispesifikasikan; mengabaikan direktori '%s'", src)
		}
		return copyDirectory(src, dst, force, interactive)
	}

	return copyFile(src, dst, force, interactive)
}

func copyFile(src, dst string, force, interactive bool) error {
	// Adjust destination if it's a directory
	destInfo, err := os.Stat(dst)
	if err == nil && destInfo.IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}

	// Check if destination file exists
	if _, err := os.Stat(dst); err == nil {
		if force {
			// Overwrite without asking
		} else if interactive {
			prompt := fmt.Sprintf("cp: timpa '%s'?", dst)
			confirmed, err := util.Confirm(prompt)
			if err != nil {
				return err
			}
			if !confirmed {
				return nil // User said no, do nothing.
			}
		} else {
			return fmt.Errorf("cp: file tujuan '%s' sudah ada. Gunakan -f atau -i.", dst)
		}
	}

	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// Preserve permissions
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, info.Mode())
}

func copyDirectory(src, dst string, force, interactive bool) error {
	// Create destination directory
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		fileInfo, err := os.Stat(srcPath)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			err = copyDirectory(srcPath, dstPath, force, interactive)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcPath, dstPath, force, interactive)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
