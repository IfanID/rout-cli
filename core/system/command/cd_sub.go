package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Dapatkan home directory dari environment variable, ini cara paling andal.
	home := os.Getenv("HOME")
	if home == "" {
		// Jika karena alasan aneh $HOME tidak ada, jangan cetak apa-apa.
		return
	}

	// Gabungkan dengan path tujuan untuk membuat path absolut.
	targetDir := filepath.Join(home, "storage", "shared", "SubConverter")
	fmt.Println(targetDir)
}