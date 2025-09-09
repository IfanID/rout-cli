package command

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"rout/cmd/core/system/util"
)

const configKeyRman = "lokasiRman"

var RmanCmd = &cobra.Command{
	Use:   "rman",
	Short: "Mengatur atau pindah ke direktori 'rman'.",
	Long: `Perintah ini membantu Anda mengelola dan berpindah ke direktori yang sering digunakan yang disebut 'rman'.

Penggunaan:
  rout rman          // Pindah ke direktori 'rman' yang tersimpan.
  rout rman -ganti   // Paksa untuk mengatur ulang lokasi direktori 'rman'.
  rout rman -lokasi  // Tampilkan lokasi direktori 'rman' yang saat ini disimpan.`, 
	Run: runRmanCommand,
}



func runRmanCommand(cmd *cobra.Command, args []string) {
	homeDir, baseDir, err := util.GetUserPaths()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error kritis:", err)
		return
	}

	configPath, _ := util.GetConfigPath()
	config, err := util.ReadConfig(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error membaca konfigurasi:", err)
		return
	}
	targetDir, exists := config[configKeyRman]

	if exists && strings.HasPrefix(targetDir, "~") {
		targetDir = strings.Replace(targetDir, "~", homeDir, 1)
	}

	if ShowLocation {
		handleShowLocationRman(targetDir, exists, baseDir, homeDir)
		return
	}

	if !exists || targetDir == "" || ForceChange {
		newTarget, err := promptForNewLocationRman(targetDir, exists, baseDir, homeDir)
		if err != nil {
			return
		}
		targetDir = newTarget

		pathToSave := util.GetTildePath(targetDir, homeDir)
		config[configKeyRman] = pathToSave

		if err := util.SaveConfig(configPath, config);
 err != nil {
			fmt.Fprintln(os.Stderr, "Gagal menyimpan konfigurasi:", err)
			return
		}
	}

	fmt.Println(targetDir)
}

func handleShowLocationRman(targetDir string, exists bool, baseDir, homeDir string) {
	if exists && targetDir != "" {
		displayPath := util.GetDisplayPath(targetDir, baseDir, homeDir)
		fmt.Fprintln(os.Stderr, "Lokasi 'rman' saat ini:", displayPath)
	} else {
		fmt.Fprintln(os.Stderr, "Lokasi 'rman' belum diatur.")
	}
}

func promptForNewLocationRman(currentTarget string, exists bool, baseDir, homeDir string) (string, error) {
	if exists && currentTarget != "" {
		displayPath := util.GetDisplayPath(currentTarget, baseDir, homeDir)
		fmt.Fprintln(os.Stderr, "Lokasi saat ini:", displayPath)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintf(os.Stderr, "Nama folder tujuan: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("gagal membaca input: %w", err)
		}

		cleanedInput := cleanInputRman(input)
		if cleanedInput == "" {
			continue
		}

		fullPath := resolveFullPathRman(cleanedInput, baseDir)

		if _, err := os.Stat(fullPath); err == nil {
			displayPath := util.GetDisplayPath(fullPath, baseDir, homeDir)
			fmt.Fprintln(os.Stderr, "Path baru berhasil disimpan:", displayPath)
			return fullPath, nil
		}

		displayPathForError := util.GetDisplayPath(fullPath, baseDir, homeDir)
		fmt.Fprintln(os.Stderr, "--> Direktori tidak ditemukan di '", displayPathForError, "'. Silakan coba lagi.")
	}
}

func cleanInputRman(input string) string {
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

func resolveFullPathRman(input, baseDir string) string {
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
