package manajemen_file

import (
	"fmt"
	"os"
	"time"
)

// Touch creates new empty files. If a file already exists, it updates its timestamp.
func Touch(filenames []string) error {
	if len(filenames) == 0 {
		return fmt.Errorf("touch: nama file harus disertakan")
	}

	for _, filename := range filenames {
		// Check if the file exists
		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			// File does not exist, create it
			file, err := os.Create(filename)
			if err != nil {
				return fmt.Errorf("gagal membuat file %s: %w", filename, err)
			}
			file.Close()
		} else if err == nil {
			// File exists, update its modification time
			currentTime := time.Now().Local()
			err := os.Chtimes(filename, currentTime, currentTime)
			if err != nil {
				return fmt.Errorf("gagal memperbarui waktu modifikasi file %s: %w", filename, err)
			}
		} else {
			// Another error occurred
			return fmt.Errorf("gagal mengakses file %s: %w", filename, err)
		}
	}
	return nil
}
