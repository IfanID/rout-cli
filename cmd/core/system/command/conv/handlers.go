package conv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"rout/cmd/core/system/util"
)

func HandleConversion(arg string) {
	if !CheckSubDir() {
		return
	}

	if !IsFzfInstalled() {
		PrintLog(0, "❌", "Error: fzf tidak ditemukan. Harap instal fzf untuk menggunakan fitur interaktif.", ColorRed)
		PrintSubLog(1, "💡", "Untuk Termux, Anda bisa menginstal dengan: pkg install fzf", ColorCyan)
		return
	}

	configPath, _ := util.GetConfigPath()
	config, err := util.ReadConfig(configPath)
	if err != nil {
		PrintLog(0, "❌", "Error membaca konfigurasi: "+err.Error(), ColorRed)
		return
	}
	lokasiCdSub, exists := config["lokasiCdSub"]
	if !exists || lokasiCdSub == "" {
		PrintLog(0, "❌", "Lokasi 'sub' belum diatur. Jalankan 'rout sub' terlebih dahulu.", ColorRed)
		return
	}

	homeDir, _, err := util.GetUserPaths()
	if err != nil {
		PrintLog(0, "❌", "Error kritis: "+err.Error(), ColorRed)
		return
	}
	if strings.HasPrefix(lokasiCdSub, "~") {
		lokasiCdSub = strings.Replace(lokasiCdSub, "~", homeDir, 1)
	}

	outputBaseDir := filepath.Join(lokasiCdSub, "ROutConv")
	if _, err := os.Stat(outputBaseDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputBaseDir, 0755); err != nil {
			PrintLog(0, "❌", fmt.Sprintf("Gagal membuat direktori output %s: %v", outputBaseDir, err), ColorRed)
			return
		}
	}

	if arg == "all" {
		PrintHeader("Konversi Semua File .ts", ColorBlue)
		files, err := filepath.Glob("*.ts")
		if err != nil {
			PrintLog(0, "❌", fmt.Sprintf("Error mencari file .ts: %v", err), ColorRed)
			return
		}
		if len(files) == 0 {
			PrintLog(0, "ℹ️", "Tidak ada file .ts yang ditemukan.", ColorYellow)
			return
		}
		for _, f := range files {
			RunFfmpegAndCue(f, outputBaseDir)
		}
	} else {
		RunFfmpegAndCue(arg, outputBaseDir)
	}
}

func CheckSubDir() bool {
	configPath, _ := util.GetConfigPath()
	config, err := util.ReadConfig(configPath)
	if err != nil {
		PrintLog(0, "❌", "Error membaca konfigurasi: "+err.Error(), ColorRed)
		return false
	}
	lokasiCdSub, exists := config["lokasiCdSub"]
	if !exists || lokasiCdSub == "" {
		PrintLog(0, "⚠️", "Lokasi 'sub' belum diatur. Jalankan 'rout sub' terlebih dahulu.", ColorYellow)
		return false
	}

	homeDir, _, err := util.GetUserPaths()
	if err != nil {
		PrintLog(0, "❌", "Error kritis: "+err.Error(), ColorRed)
		return false
	}
	if strings.HasPrefix(lokasiCdSub, "~") {
		lokasiCdSub = strings.Replace(lokasiCdSub, "~", homeDir, 1)
	}

	wd, err := os.Getwd()
	if err != nil {
		PrintLog(0, "❌", "Error mendapatkan direktori saat ini: "+err.Error(), ColorRed)
		return false
	}

	if wd != lokasiCdSub {
		PrintLog(0, "⚠️", "Perintah 'conv' hanya bisa dijalankan di dalam direktori 'sub' yang telah diatur.", ColorYellow)
		PrintSubLog(1, "↪️", "Silakan jalankan 'sub' untuk pindah ke direktori yang benar.", ColorGray)
		return false
	}
	return true
}
