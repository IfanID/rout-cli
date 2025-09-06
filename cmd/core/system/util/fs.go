package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// GetUserPaths mendapatkan path home dan basis pengguna.
func GetUserPaths() (homeDir, baseDir string, err error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", "", fmt.Errorf("gagal mendapatkan info pengguna: %w", err)
	}
	homeDir = currentUser.HomeDir
	baseDir = filepath.Join(homeDir, "storage", "shared")
	return homeDir, baseDir, nil
}

// GetConfigPath membangun path absolut ke file directory.json.
func GetConfigPath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("gagal mendapatkan info pengguna: %w", err)
	}
	return filepath.Join(currentUser.HomeDir, ".rout", "directory.json"), nil
}

// ReadConfig membaca dan mem-parsing file konfigurasi JSON.
func ReadConfig(path string) (map[string]string, error) {
	config := make(map[string]string)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return config, nil // Kembalikan config kosong jika file tidak ada
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file config: %w", err)
	}
	if len(file) > 0 {
		// Jangan hiraukan error unmarshal jika file kosong atau rusak
		_ = json.Unmarshal(file, &config)
	}
	return config, nil
}

// SaveConfig menulis map konfigurasi ke file JSON.
func SaveConfig(path string, config map[string]string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("gagal marshal JSON: %w", err)
	}
	return ioutil.WriteFile(path, data, 0644)
}

// GetDisplayPath mengubah path absolut menjadi format yang lebih ramah pengguna.
func GetDisplayPath(fullPath, baseDir, homeDir string) string {
	// Coba buat relatif ke baseDir dulu, ini paling relevan untuk 'storage/shared'
	displayPath, err := filepath.Rel(baseDir, fullPath)
	if err == nil {
		return displayPath
	}
	// Jika gagal (misal path di luar baseDir), gunakan format tilde
	return GetTildePath(fullPath, homeDir)
}

// GetTildePath mengubah path absolut menjadi format ~.
func GetTildePath(fullPath, homeDir string) string {
	if strings.HasPrefix(fullPath, homeDir) {
		return strings.Replace(fullPath, homeDir, "~", 1)
	}
	return fullPath
}
