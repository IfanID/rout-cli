package ui

import (
	"os"
	"os/user"
	"path/filepath"
)

// setupZshEnvironment memastikan lingkungan Zsh kustom untuk rout sudah siap.
// Fungsi ini membuat direktori ~/.rout dan file .zshrc di dalamnya jika belum ada.
// Mengembalikan path ke ZDOTDIR kustom dan error jika terjadi.
func setupZshEnvironment() (string, error) {
	// Dapatkan direktori home pengguna
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	homeDir := currentUser.HomeDir

	// Buat path ZDOTDIR kustom
	customZdotdir := filepath.Join(homeDir, ".rout")

	// Pastikan direktori ZDOTDIR kustom ada
	if _, err := os.Stat(customZdotdir); os.IsNotExist(err) {
		if err := os.MkdirAll(customZdotdir, 0755); err != nil {
			return "", err
		}
	}

	// Periksa apakah .zshrc ada di direktori kustom, jika tidak, salin dari template
	zshrcPath := filepath.Join(customZdotdir, ".zshrc")
	if _, err := os.Stat(zshrcPath); os.IsNotExist(err) {
		// Cari path executable untuk menemukan template
		exePath, err := os.Executable()
		if err != nil {
			// Fallback untuk keamanan
			fallbackContent := []byte("# Fallback .zshrc, tidak dapat menemukan path executable\n")
			_ = os.WriteFile(zshrcPath, fallbackContent, 0644)
		} else {
			exeDir := filepath.Dir(exePath)
			templatePath := filepath.Join(exeDir, "config", "zshrc")

			// Baca template
			templateContent, err := os.ReadFile(templatePath)
			if err != nil {
				// Jika template tidak ditemukan, buat .zshrc minimal untuk menghindari error
				fallbackContent := []byte("# Fallback .zshrc, template tidak ditemukan\n")
				_ = os.WriteFile(zshrcPath, fallbackContent, 0644)
			} else {
				// Tulis konten template ke .zshrc baru
				err = os.WriteFile(zshrcPath, templateContent, 0644)
				if err != nil {
					return "", err
				}
			}
		}
	}
	return customZdotdir, nil
}
