package command

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"rout/cmd/API"
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
	outputFile := filepath.Join(outDir, strings.TrimSuffix(inputFile, ".ts")+".vtt")

	fmt.Printf("▶️  Mengonversi %s ke format VTT...\n", inputFile)
	cmd := exec.Command("ffmpeg", "-y", "-i", inputFile, outputFile)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error menjalankan ffmpeg untuk %s: %v\n", inputFile, err)
		fmt.Fprintf(os.Stderr, "FFmpeg output:\n%s\n", stderr.String())
		return
	}
	fmt.Printf("✅ Konversi selesai: %s\n", outputFile)

	fmt.Println("🔍 Menambahkan atau memperbarui CUE...")
	err = addCue(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error menambahkan CUE ke %s: %v\n", outputFile, err)
		return
	}

	// Jeda singkat untuk memastikan sistem file selesai menulis sebelum verifikasi
	time.Sleep(100 * time.Millisecond)

	if !verifyCue(outputFile) {
		fmt.Fprintf(os.Stderr, "⚠️  Peringatan: Gagal memverifikasi CUE di file %s.\n", outputFile)
	} else {
		fmt.Println("✅ CUE berhasil ditambahkan atau diperbarui")
	}

	fmt.Printf("🎉 Semua proses selesai untuk: %s\n", outputFile)
}

func addCue(filename string) error {
	tmpfile := filename + ".tmp"
	lines, err := readLines(filename)
	if err != nil {
		return fmt.Errorf("gagal membaca file untuk addCue: %w", err)
	}

	cueContent := "Ifan.3 V2S CoreX"

	// --- Block Analisis Gemini & Penyesuaian Timestamp ---
	var finalTimestamp string

	apiKey := API.GetGeminiAPIKey()
	if apiKey == "" {
		return fmt.Errorf("GEMINI_API_KEY tidak ditemukan. Silakan atur environment variable atau buat file .env")
	}

	fmt.Println("🔍 Menghubungkan ke Gemini API untuk analisis CUE...")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // Timeout dinaikkan menjadi 60 detik
	defer cancel()

	vttHeader := API.ExtractVTTHeader(lines, 10*time.Second)
	geminiTimestamp, reason, err := API.AnalyzeVTTForCUE(ctx, apiKey, vttHeader)
	if err != nil {
		return fmt.Errorf("gagal menganalisis konten dengan Gemini API: %w", err)
	}

	fmt.Printf("🤖 Gemini menyarankan timestamp: %s\n", geminiTimestamp)
	fmt.Printf("📋 Alasan: %s\n", reason)
	
	fmt.Println("\n--- Info Analisis Detail ---")
	fmt.Println("Cuplikan Subtitle yang Dianalisis:")
	fmt.Println(vttHeader)
	fmt.Println("--------------------------")

	// Terapkan aturan baru
	tsParts := strings.Split(geminiTimestamp, " --> ")
	if len(tsParts) != 2 {
		return fmt.Errorf("timestamp dari Gemini tidak valid: %s", geminiTimestamp)
	}

	cueStart, err := parseVTTDuration(tsParts[0])
	if err != nil {
		return fmt.Errorf("gagal parse waktu mulai CUE: %w", err)
	}
	cueEnd, err := parseVTTDuration(tsParts[1])
	if err != nil {
		return fmt.Errorf("gagal parse waktu akhir CUE: %w", err)
	}

	fmt.Println("Proses Penyesuaian Timestamp:")
	
	// Aturan 1: Hindari Tumpang Tindih
	if firstDialogueStart, found := findFirstDialogueStartTime(lines); found {
		fmt.Printf("- Waktu mulai dialog pertama terdeteksi: %s\n", formatVTTDuration(firstDialogueStart))
		if cueEnd >= firstDialogueStart {
			fmt.Printf("  ℹ️  Waktu akhir CUE (%s) menabrak/melewati dialog pertama. Menyesuaikan...\n", formatVTTDuration(cueEnd))
			cueEnd = firstDialogueStart - time.Millisecond
		}
	}

	// Aturan 2: Durasi Minimal 2 Detik
	minDuration := 2 * time.Second
	if cueEnd-cueStart < minDuration {
		fmt.Printf("  ℹ️  Durasi CUE (%s) kurang dari 2 detik. Menyesuaikan...\n", (cueEnd-cueStart).String())
		cueEnd = cueStart + minDuration
	}

	finalTimestamp = formatVTTDuration(cueStart) + " --> " + formatVTTDuration(cueEnd)
	fmt.Println("--------------------------")
	fmt.Printf("✅ Timestamp final setelah penyesuaian: %s\n\n", finalTimestamp)
	// --- Akhir Block ---

	cueBlock := []string{finalTimestamp, cueContent, ""}

	cueFoundIndex := -1
	for i, line := range lines {
		if strings.Contains(line, cueContent) {
			cueFoundIndex = i
			break
		}
	}

	if cueFoundIndex != -1 {
		fmt.Println("🔄 Memperbarui CUE yang sudah ada...")
		if cueFoundIndex > 0 {
			lines[cueFoundIndex-1] = finalTimestamp
		}
		lines[cueFoundIndex] = cueContent
		if cueFoundIndex+1 < len(lines) && lines[cueFoundIndex+1] != "" {
			lines = insertStringSlice(lines, cueFoundIndex+1, "")
		} else if cueFoundIndex+1 == len(lines) {
			lines = append(lines, "")
		}
	} else {
		fmt.Println("➕ Menambahkan CUE baru...")
		timestampRegex, err := regexp.Compile(`^(\d{2}:)?\d{2}:\d{2}\.\d{3}\s-->\s(\d{2}:)?\d{2}:\d{2}\.\d{3}`)
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

	fmt.Println("\n--- Hasil Akhir Sebelum Disimpan ---")
	for i := 0; i < 15 && i < len(lines); i++ {
		fmt.Println(lines[i])
	}
	fmt.Println("---------------------------------")

	outFile, err := os.Create(tmpfile)
	if err != nil {
		return fmt.Errorf("gagal membuat file sementara: %w", err)
	}

	writer := bufio.NewWriter(outFile)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			outFile.Close() // Pastikan file ditutup sebelum return
			return fmt.Errorf("gagal menulis ke file sementara: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		outFile.Close() // Pastikan file ditutup sebelum return
		return fmt.Errorf("gagal flush writer: %w", err)
	}

	// Tutup file secara eksplisit sebelum me-rename
	if err := outFile.Close(); err != nil {
		return fmt.Errorf("gagal menutup file sementara: %w", err)
	}

	// Rename file sementara ke file asli
	return os.Rename(tmpfile, filename)
}

// parseVTTDuration mengubah VTT timestamp string (misal: 00:01:02.345) menjadi time.Duration
func parseVTTDuration(ts string) (time.Duration, error) {
	var h, m, s, ms int
	// Format bisa hh:mm:ss.ms atau mm:ss.ms
	parts := strings.Split(ts, ":")
	var err error
	if len(parts) == 3 { // hh:mm:ss.ms
		h, err = strconv.Atoi(parts[0])
		if err != nil { return 0, err }
		m, err = strconv.Atoi(parts[1])
		if err != nil { return 0, err }
		sMs := strings.Split(parts[2], ".")
		s, err = strconv.Atoi(sMs[0])
		if err != nil { return 0, err }
		ms, err = strconv.Atoi(sMs[1])
		if err != nil { return 0, err }
	} else if len(parts) == 2 { // mm:ss.ms
		m, err = strconv.Atoi(parts[0])
		if err != nil { return 0, err }
		sMs := strings.Split(parts[1], ".")
		s, err = strconv.Atoi(sMs[0])
		if err != nil { return 0, err }
		ms, err = strconv.Atoi(sMs[1])
		if err != nil { return 0, err }
	} else {
		return 0, fmt.Errorf("format timestamp tidak valid: %s", ts)
	}

	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second + time.Duration(ms)*time.Millisecond, nil
}

// formatVTTDuration mengubah time.Duration menjadi VTT timestamp string (hh:mm:ss.ms)
func formatVTTDuration(d time.Duration) string {
	d = d.Round(time.Millisecond)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	d -= s * time.Second
	ms := d / time.Millisecond
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}

// findFirstDialogueStartTime mencari waktu mulai subtitle dialog pertama
func findFirstDialogueStartTime(lines []string) (time.Duration, bool) {
	re := regexp.MustCompile(`^(\d{2}:)?\d{2}:\d{2}\.\d{3}\s-->\s(\d{2}:)?\d{2}:\d{2}\.\d{3}`)
	for i, line := range lines {
		if i > 0 && strings.Contains(lines[i-1], "Ifan.3 V2S CoreX") {
			continue // Lewati timestamp dari CUE yang sudah ada
		}
		if re.MatchString(line) {
			parts := strings.Split(line, " --> ")
			startTime, err := parseVTTDuration(parts[0])
			if err == nil {
				return startTime, true
			}
		}
	}
	return 0, false // Tidak ditemukan dialog
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
