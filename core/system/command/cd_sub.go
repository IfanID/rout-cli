package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const configKey = "lokasiCdSub"

// ====================================================================
// # Bagian 1: Fungsi Utama & Alur Kontrol
// ====================================================================

func main() {
	// 1. Setup: Dapatkan path penting dan parse flag dari argumen
	homeDir, baseDir, err := getUserPaths()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error kritis:", err)
		return
	}

	showLocation, forceChange := parseFlags(os.Args)

	// 2. Konfigurasi: Baca file konfigurasi yang ada
	configPath, _ := getConfigPath()
	config, err := readConfig(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error membaca konfigurasi:", err)
		return
	}
	targetDir, exists := config[configKey]

	// --- PENTING: Ubah path dari format ~ ke absolut setelah dibaca ---
	if exists && strings.HasPrefix(targetDir, "~") {
		targetDir = strings.Replace(targetDir, "~", homeDir, 1)
	}

	// 3. Aksi Berdasarkan Flag: Jalankan tugas spesifik jika ada flag
	if showLocation {
		handleShowLocation(targetDir, exists, baseDir, homeDir)
		return
	}

	// 4. Alur Utama: Tentukan apakah perlu meminta lokasi baru
	if !exists || targetDir == "" || forceChange {
		newTarget, err := promptForNewLocation(targetDir, exists, baseDir, homeDir)
		if err != nil {
			return
		}
		targetDir = newTarget // Gunakan path absolut untuk sesi ini
		
		// --- PENTING: Ubah path ke format ~ sebelum disimpan ---
		pathToSave := getTildePath(targetDir, homeDir)
		config[configKey] = pathToSave

		if err := saveConfig(configPath, config); err != nil {
			fmt.Fprintln(os.Stderr, "Gagal menyimpan konfigurasi:", err)
			return
		}
	}

	// 5. Output: Cetak direktori final untuk dieksekusi oleh shell
	fmt.Println(targetDir)
}

// ====================================================================
// # Bagian 2: Fungsi-fungsi Bantuan Logika
// ====================================================================

// parseFlags memproses argumen baris perintah dan mengembalikan flag yang relevan.
func parseFlags(args []string) (showLocation, forceChange bool) {
	if len(args) <= 1 {
		return false, false
	}
	for _, arg := range args[1:] {
		switch arg {
		case "-ganti":
			forceChange = true
		case "-lokasi":
			showLocation = true
		}
	}
	return showLocation, forceChange
}

// handleShowLocation menangani logika untuk flag -lokasi.
func handleShowLocation(targetDir string, exists bool, baseDir, homeDir string) {
	if exists && targetDir != "" {
		displayPath := getDisplayPath(targetDir, baseDir, homeDir)
		fmt.Fprintln(os.Stderr, "Lokasi 'sub' saat ini:", displayPath)
	} else {
		fmt.Fprintln(os.Stderr, "Lokasi 'sub' belum diatur.")
	}
}

// promptForNewLocation memulai loop interaktif untuk meminta lokasi baru dari pengguna.
func promptForNewLocation(currentTarget string, exists bool, baseDir, homeDir string) (string, error) {
	if exists && currentTarget != "" {
		displayPath := getDisplayPath(currentTarget, baseDir, homeDir)
		fmt.Fprintln(os.Stderr, "Lokasi saat ini:", displayPath)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintf(os.Stderr, "Nama folder tujuan: ")
		
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("gagal membaca input: %w", err)
		}

		cleanedInput := cleanInput(input)
		if cleanedInput == "" {
			continue
		}

		fullPath := resolveFullPath(cleanedInput, baseDir)

		if _, err := os.Stat(fullPath); err == nil {
			displayPath := getDisplayPath(fullPath, baseDir, homeDir)
			fmt.Fprintln(os.Stderr, "Path baru berhasil disimpan:", displayPath)
			return fullPath, nil // Sukses
		}

		// Jika gagal, cetak pesan error
		displayPathForError := getDisplayPath(fullPath, baseDir, homeDir)
		fmt.Fprintln(os.Stderr, "--> Direktori tidak ditemukan di '", displayPathForError, "'. Silakan coba lagi.")
	}
}

// ====================================================================
// # Bagian 3: Fungsi-fungsi Utilitas Murni
// ====================================================================

// getUserPaths mendapatkan path home dan basis pengguna.
func getUserPaths() (homeDir, baseDir string, err error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", "", fmt.Errorf("gagal mendapatkan info pengguna: %w", err)
	}
	homeDir = currentUser.HomeDir
	baseDir = filepath.Join(homeDir, "storage", "shared")
	return homeDir, baseDir, nil
}

// getConfigPath membangun path absolut ke file directory.json.
func getConfigPath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("gagal mendapatkan info pengguna: %w", err)
	}
	return filepath.Join(currentUser.HomeDir, ".rout", "directory.json"), nil
}

// readConfig membaca dan mem-parsing file konfigurasi JSON.
func readConfig(path string) (map[string]string, error) {
	config := make(map[string]string)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return config, nil
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file config: %w", err)
	}
	if len(file) > 0 {
		_ = json.Unmarshal(file, &config)
	}
	return config, nil
}

// saveConfig menulis map konfigurasi ke file JSON.
func saveConfig(path string, config map[string]string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("gagal marshal JSON: %w", err)
	}
	return ioutil.WriteFile(path, data, 0644)
}

// cleanInput membersihkan string input dari pengguna dengan sangat teliti.
func cleanInput(input string) string {
	cleaned := strings.TrimSpace(input)
	if !filepath.IsAbs(cleaned) && !strings.HasPrefix(cleaned, "/storage/emulated/0/") {
		parts := strings.Split(cleaned, "/")
		for i, part := range parts {
			parts[i] = strings.TrimSpace(part)
		}
		return strings.Join(parts, "/")
	}
	return cleaned
}

// resolveFullPath menentukan path absolut dari input pengguna.
func resolveFullPath(input, baseDir string) string {
	androidRoot := "/storage/emulated/0/"
	if strings.HasPrefix(input, androidRoot) {
		relativePath := strings.TrimPrefix(input, androidRoot)
		return filepath.Join(baseDir, relativePath)
	}
	if filepath.IsAbs(input) {
		return input
	}
	return filepath.Join(baseDir, input)
}

// getDisplayPath mengubah path absolut menjadi format yang lebih ramah pengguna.
func getDisplayPath(fullPath, baseDir, homeDir string) string {
	displayPath, err := filepath.Rel(baseDir, fullPath)
	if err == nil {
		return displayPath
	}
	return getTildePath(fullPath, homeDir)
}

// getTildePath mengubah path absolut menjadi format ~.
func getTildePath(fullPath, homeDir string) string {
	if strings.HasPrefix(fullPath, homeDir) {
		return strings.Replace(fullPath, homeDir, "~", 1)
	}
	return fullPath
}