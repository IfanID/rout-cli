package conv

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func RunFfmpegAndCue(inputFile string, outDir string) {
	var err error
	var choice string

	PrintHeader(fmt.Sprintf("Proses File: %s", inputFile), ColorCyan)

	if !strings.HasSuffix(inputFile, ".ts") {
		PrintLog(1, "❌", fmt.Sprintf("Error: File input '%s' harus berekstensi .ts", inputFile), ColorRed)
		return
	}
	outputFile := filepath.Join(outDir, strings.TrimSuffix(inputFile, ".ts") + ".vtt")

	fmt.Println()
	PrintLog(1, "🔍", "Mengecek file output...", ColorBlue)
	if _, err = os.Stat(outputFile); err == nil {
		PrintLog(2, "⚠️", fmt.Sprintf("File output '%s' sudah ada.", filepath.Base(outputFile)), ColorYellow)
		options := []string{"1. Timpa File", "2. Lewati File Ini"}
		choice, err = ShowFzfMenu("Pilihan untuk file yang sudah ada: ", options)
		if err != nil {
			PrintLog(2, "❗", "Proses dibatalkan, melewati file.", ColorYellow)
			return
		}

		switch choice {
		case "1. Timpa File":
			PrintLog(2, "✅", "Memilih untuk menimpa file. Melanjutkan konversi...", ColorGreen)
		case "2. Lewati File Ini":
			PrintLog(2, "↪️", "Melewati file ini.", ColorYellow)
			PrintHeader(fmt.Sprintf("Selesai: %s", inputFile), ColorGreen)
			return
		default:
			PrintLog(2, "❗", "Pilihan tidak valid, melewati file.", ColorYellow)
			return
		}
	} else {
		PrintLog(2, "✅", "File output belum ada, melanjutkan proses.", ColorGreen)
	}

	fmt.Println()
	PrintLog(1, "▶️", "Menjalankan konversi FFMPEG...", ColorBlue)
	cmd := exec.Command("ffmpeg", "-y", "-i", inputFile, outputFile)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		PrintLog(2, "❌", fmt.Sprintf("Error menjalankan ffmpeg untuk %s: %v", inputFile, err), ColorRed)
		PrintSubLog(2, "📕", "FFmpeg output:\n"+stderr.String(), ColorGray)
		return
	}
	PrintLog(2, "✅", fmt.Sprintf("Konversi selesai: %s", filepath.Base(outputFile)), ColorGreen)

	// Menampilkan pratinjau setelah konversi ffmpeg
	PrintFilePreview(2, "Pratinjau Hasil Konversi FFMPEG (15 baris pertama):", outputFile, 15, ColorBlue)

	loopActive := true
	for loopActive {
		fmt.Println()
		PrintLog(1, "🤖", "Memulai analisis CUE dengan Gemini API...", ColorBlue)
		err = AddCue(outputFile)
		if err != nil {
			PrintLog(2, "❌", fmt.Sprintf("Error pada proses CUE: %v", err), ColorRed)
			return
		}

		time.Sleep(100 * time.Millisecond)

		if !VerifyCue(outputFile) {
			PrintLog(2, "⚠️", fmt.Sprintf("Peringatan: Gagal memverifikasi CUE di file %s.", filepath.Base(outputFile)), ColorYellow)
		} else {
			PrintLog(2, "✅", "CUE berhasil ditambahkan atau diperbarui.", ColorGreen)
		}

		fmt.Println()
		PrintLog(1, "❓", "Tindakan selanjutnya...", ColorBlue)
		options := []string{"1. Simpan", "2. Coba Lagi", "3. Hapus"}
		choice, err = ShowFzfMenu("Pilihan untuk file yang sudah ada: ", options)
		if err != nil {
			PrintLog(2, "❗", "Proses dibatalkan, menyimpan file secara default.", ColorYellow)
			choice = "1. Simpan"
		}

		switch choice {
		case "1. Simpan":
			PrintLog(2, "💾", "Menyimpan file dan menyelesaikan proses.", ColorGreen)
			loopActive = false
		case "2. Coba Lagi":
			PrintLog(2, "🔄", "Mengulang proses analisis CUE...", ColorBlue)
		case "3. Hapus":
			PrintLog(2, "🗑️", "Menghapus file...", ColorRed)
			err = os.Remove(outputFile)
			if err != nil {
				PrintSubLog(3, "❌", fmt.Sprintf("Gagal menghapus file %s: %v", filepath.Base(outputFile), err), ColorRed)
			} else {
				PrintSubLog(3, "✅", fmt.Sprintf("File %s berhasil dihapus.", filepath.Base(outputFile)), ColorGreen)
			}
			loopActive = false
		default:
			PrintLog(2, "❗", "Pilihan tidak valid, menyimpan file.", ColorYellow)
			loopActive = false
		}
	}
	PrintHeader(fmt.Sprintf("🎉 Selesai: %s", inputFile), ColorGreen)
}

func IsFzfInstalled() bool {
	_, err := exec.LookPath("fzf")
	return err == nil
}
