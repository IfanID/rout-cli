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
	// Siapkan konfigurasi starship
	if err := setupStarship(); err != nil {
		// Jika terjadi error, kita bisa log atau menanganinya di sini.
		// Untuk saat ini, kita tidak menghentikan proses utama.
	}

	// Siapkan konfigurasi directory.json
	if err := setupDirectoryConfig(); err != nil {
		// Jika terjadi error, kita bisa log atau menanganinya di sini.
		// Untuk saat ini, kita tidak menghentikan proses utama.
	}

	

	

	return customZdotdir, nil
}

// setupStarship memastikan file konfigurasi starship.toml ada di ~/.config.
// Jika tidak ada atau kosong, ia akan menyalinnya dari template proyek.
func setupStarship() error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	homeDir := currentUser.HomeDir

	// Path ke file konfigurasi starship pengguna
	configDir := filepath.Join(homeDir, ".config")
	starshipConfigPath := filepath.Join(configDir, "starship.toml")

	// Periksa apakah file sudah ada dan tidak kosong
	info, err := os.Stat(starshipConfigPath)
	if err == nil && info.Size() > 0 {
		// File sudah ada dan tidak kosong, tidak perlu melakukan apa-apa
		return nil
	}

	// Jika file tidak ada atau kosong, buat/timpa
	// Pastikan direktori ~/.config ada
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return err
		}
	}

	// Dapatkan path template dari direktori executable
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)
	templatePath := filepath.Join(exeDir, "config", "starship.toml")

	// Baca template
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}

	// Tulis template ke file konfigurasi pengguna
	return os.WriteFile(starshipConfigPath, templateContent, 0644)
}

// setupDirectoryConfig memastikan file konfigurasi directory.json ada di ~/.rout.
// Jika tidak ada, ia akan menyalinnya dari template proyek.
func setupDirectoryConfig() error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	homeDir := currentUser.HomeDir

	// Path ke direktori konfigurasi rout
	routConfigDir := filepath.Join(homeDir, ".rout")
	directoryConfigPath := filepath.Join(routConfigDir, "directory.json")

	// Periksa apakah file sudah ada. Jika ya, tidak perlu melakukan apa-apa.
	if _, err := os.Stat(directoryConfigPath); err == nil {
		return nil
	}

	// Dapatkan path template dari direktori executable
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)
	templatePath := filepath.Join(exeDir, "config", "directory.json")

	// Baca template
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}

	// Tulis template ke file konfigurasi pengguna
	return os.WriteFile(directoryConfigPath, templateContent, 0644)
}


