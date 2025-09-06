package command

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"rout/cmd/core/system/util"
)

var ConvCmd = &cobra.Command{
	Use:   "conv [file.ts|all]",
	Short: "Konversi file .ts ke .vtt dengan CUE.",
	Long: `Alat ini mengonversi file video .ts menjadi file subtitle .vtt.
Ini juga menambahkan atau memperbarui blok CUE kustom di awal file .vtt.

Penggunaan:
  rout conv               // Menampilkan informasi penggunaan dan lokasi 'sub'.
  rout conv <namafile.ts>  // Mengonversi satu file
  rout conv all             // Mengonversi semua file .ts di direktori saat ini`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if checkSubDir() {
				configPath, _ := util.GetConfigPath()
				config, err := util.ReadConfig(configPath)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error membaca konfigurasi:", err)
					return
				}
			lokasiCdSub, exists := config["lokasiCdSub"]
			if exists && lokasiCdSub != "" {
				fmt.Printf("Perintah 'conv' siap digunakan di direktori: %s\n", lokasiCdSub)
				fmt.Println("Gunakan 'conv all' untuk mengonversi semua file .ts, atau 'conv <namafile.ts>' untuk satu file.")
			} else {
				fmt.Fprintln(os.Stderr, "Lokasi 'sub' belum diatur. Jalankan 'rout sub' terlebih dahulu.")
			}
		}
	} else {
		handleConversion(args[0])
	}
	},
}

func checkSubDir() bool {
	configPath, _ := util.GetConfigPath()
	config, err := util.ReadConfig(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error membaca konfigurasi:", err)
		return false
	}
	lokasiCdSub, exists := config["lokasiCdSub"]
	if !exists || lokasiCdSub == "" {
		fmt.Fprintln(os.Stderr, "Lokasi 'sub' belum diatur. Jalankan 'rout sub' terlebih dahulu.")
		return false
	}

	homeDir, _, err := util.GetUserPaths()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error kritis:", err)
		return false
	}
	if strings.HasPrefix(lokasiCdSub, "~") {
		lokasiCdSub = strings.Replace(lokasiCdSub, "~", homeDir, 1)
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error mendapatkan direktori saat ini:", err)
		return false
	}

	if wd != lokasiCdSub {
		fmt.Fprintln(os.Stderr, "Perintah 'conv' hanya bisa dijalankan di dalam direktori 'sub' yang telah diatur.")
		fmt.Fprintln(os.Stderr, "Silakan jalankan 'sub' untuk pindah ke direktori yang benar.")
		return false
	}
	return true
}

func handleConversion(arg string) {
	if !checkSubDir() {
		return
	}

	configPath, _ := util.GetConfigPath()
	config, err := util.ReadConfig(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error membaca konfigurasi:", err)
		return
	}
	lokasiCdSub, exists := config["lokasiCdSub"]
	if !exists || lokasiCdSub == "" {
		// This case should ideally be caught by checkSubDir, but as a fallback
		fmt.Fprintln(os.Stderr, "Lokasi 'sub' belum diatur. Jalankan 'rout sub' terlebih dahulu.")
		return
	}

	homeDir, _, err := util.GetUserPaths()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error kritis:", err)
		return
	}
	if strings.HasPrefix(lokasiCdSub, "~") {
		lokasiCdSub = strings.Replace(lokasiCdSub, "~", homeDir, 1)
	}

	outputBaseDir := filepath.Join(lokasiCdSub, "ROutConv")
	if _, err := os.Stat(outputBaseDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputBaseDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Gagal membuat direktori output %s: %v\n", outputBaseDir, err)
			return
		}
	}

	if arg == "all" {
		fmt.Println("🔎 Mengonversi semua file .ts...")
		files, err := filepath.Glob("*.ts")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error mencari file .ts: %v\n", err)
			return
		}
		if len(files) == 0 {
			fmt.Println("Tidak ada file .ts yang ditemukan.")
			return
		}
		for _, f := range files {
			runFfmpegAndCue(f, outputBaseDir)
		}
	} else {
		runFfmpegAndCue(arg, outputBaseDir)
		}
}

func runFfmpegAndCue(inputFile string, outDir string) {
	if !strings.HasSuffix(inputFile, ".ts") {
		fmt.Fprintf(os.Stderr, "Error: File input '%s' harus berekstensi .ts\n", inputFile)
		return
	}
	outputFile := filepath.Join(outDir, strings.TrimSuffix(inputFile, ".ts") + ".vtt")

	fmt.Printf("Mengonversi %s...\n", inputFile)
	cmd := exec.Command("ffmpeg", "-y", "-i", inputFile, outputFile)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error menjalankan ffmpeg untuk %s: %v\n", inputFile, err)
		fmt.Fprintf(os.Stderr, "FFmpeg output:\n%s\n", stderr.String())
		return
	}

	err = addCue(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error menambahkan CUE ke %s: %v\n", outputFile, err)
		return
	}

	if !verifyCue(outputFile) {
		fmt.Fprintf(os.Stderr, "⚠️  Peringatan: Gagal memverifikasi CUE di file %s.\n", outputFile)
	}

	fmt.Printf("✅ Selesai: %s\n", outputFile)
}

func addCue(filename string) error {
	tmpfile := filename + ".tmp"
	lines, err := readLines(filename)
	if err != nil {
		return fmt.Errorf("gagal membaca file untuk addCue: %w", err)
	}

	cueContent := "Ifan.3 V2S CoreX"
	newTimestamp := "00:00.003 --> 00:02.500"
	cueBlock := []string{newTimestamp, cueContent, ""}

	cueFoundIndex := -1
	for i, line := range lines {
		if strings.Contains(line, cueContent) {
			cueFoundIndex = i
			break
		}
	}

	if cueFoundIndex != -1 {
		if cueFoundIndex > 0 {
			lines[cueFoundIndex-1] = newTimestamp
		}
		lines[cueFoundIndex] = cueContent
		if cueFoundIndex+1 < len(lines) && lines[cueFoundIndex+1] != "" {
			lines = insertStringSlice(lines, cueFoundIndex+1, "")
		} else if cueFoundIndex+1 == len(lines) {
			lines = append(lines, "")
		}
	} else {
		timestampRegex, err := regexp.Compile(`^\d{2}:\d{2}:\d{2}\.\d{3}\s-->\s\d{2}:\d{2}:\d{2}\.\d{3}`)
		if err != nil {
			return fmt.Errorf("gagal compile regex: %w", err)
		}

		firstTimestampIndex := -1
		for i, line := range lines {
			if timestampRegex.MatchString(line) {
				firstTimestampIndex = i
				break
			}
		}

		if firstTimestampIndex != -1 {
			lines = insertStringSlice(lines, firstTimestampIndex, cueBlock...)
		} else {
			lines = append(lines, cueBlock...)
		}
	}

	outFile, err := os.Create(tmpfile)
	if err != nil {
		return fmt.Errorf("gagal membuat file sementara: %w", err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	for _, line := range lines {
		writer.WriteString(line + "\n")
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("gagal flush writer: %w", err)
	}

	return os.Rename(tmpfile, filename)
}

func verifyCue(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Ifan.3 V2S CoreX") {
			return true
		}
	}
	return false
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func insertStringSlice(slice []string, index int, values ...string) []string {
	if index < 0 || index > len(slice) {
		return slice
	}
	result := make([]string, len(slice)+len(values))
	copy(result[:index], slice[:index])
	copy(result[index:index+len(values)], values)
	copy(result[index+len(values):], slice[index:])
	return result
}