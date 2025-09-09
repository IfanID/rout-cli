package API

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// AnalyzeVTTForCUE menganalisis konten VTT dan menentukan penempatan CUE terbaik
func AnalyzeVTTForCUE(ctx context.Context, apiKey string, vttContent string) (string, string, error) {
	// Membuat klien Gemini
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", "", fmt.Errorf("gagal membuat klien Gemini: %w", err)
	}
	defer client.Close()

	// Mengatur model
	model := client.GenerativeModel("gemini-2.5-flash")
	model.SetTemperature(0.3)

	// Membuat prompt untuk analisis
	prompt := fmt.Sprintf(`
Anda adalah ahli subtitle. Tugas Anda adalah menentukan timestamp yang ideal untuk sebuah CUE branding ("Ifan.3 V2S CoreX") di awal file WebVTT.

Aturan Penting:
1.  **Kondisi Jeda (Gap):** Periksa waktu mulai dialog pertama. Jika dialog pertama dimulai **setelah** '00:00:04.000', **JANGAN** menimpa dialog tersebut. Tempatkan CUE di dalam jeda kosong tersebut (misal: '00:00:00.003 --> 00:00:02.503').
2.  **Kondisi Timpa (Overwrite):** Jika dialog pertama dimulai **sebelum** '00:00:04.000', maka timestamp CUE **HARUS** tumpang tindih (overlap) dengan dialog pertama tersebut untuk menggantikannya.
3.  **Durasi Ideal:** Buat durasi CUE sekitar 2 hingga 2.5 detik.
4.  **Penempatan:** Saat menimpa, mulai CUE di '00:00:00.003'. Sesuaikan waktu berakhirnya ('end time') agar tumpang tindih dengan dialog pertama, namun usahakan tidak menyentuh dialog kedua jika memungkinkan.

Analisis konten WebVTT berikut dan berikan timestamp berdasarkan aturan di atas.

Berikan respons dalam format berikut:
TIMESTAMP: [timestamp yang direkomendasikan]
ALASAN: [penjelasan singkat mengapa timestamp ini dipilih, sebutkan dialog mana yang ditimpa atau jika ditempatkan di jeda]

Konten WebVTT:
%s
`, vttContent)

	// Mengirim permintaan ke Gemini
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", "", fmt.Errorf("gagal menghasilkan konten: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", "", fmt.Errorf("tidak ada respons dari Gemini")
	}

	// Mengekstrak teks respons
	responseText := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			responseText += string(text)
		}
	}

	// Mengekstrak timestamp dan alasan dari respons
	timestamp, reason := parseGeminiResponse(responseText)
	if timestamp == "" {
		return "", "", fmt.Errorf("gagal mengekstrak informasi dari respons Gemini")
	}

	return timestamp, reason, nil
}

// parseGeminiResponse mengekstrak timestamp dan alasan dari respons Gemini
func parseGeminiResponse(response string) (string, string) {
	var timestamp, reason string

	scanner := bufio.NewScanner(strings.NewReader(response))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "TIMESTAMP:") {
			timestamp = strings.TrimSpace(strings.TrimPrefix(line, "TIMESTAMP:"))
		} else if strings.HasPrefix(line, "ALASAN:") {
			reason = strings.TrimSpace(strings.TrimPrefix(line, "ALASAN:"))
		}
	}

	return timestamp, reason
}

// ExtractVTTHeader mengekstrak bagian awal file VTT hingga durasi tertentu
func ExtractVTTHeader(vttContent []string, duration time.Duration) string {
	if len(vttContent) == 0 {
		return ""
	}

	var result []string
	
	// Menambahkan header WEBVTT
	result = append(result, vttContent[0])
	
	// Menghitung waktu akhir dalam detik
	endTime := duration.Seconds()
	
	// Memproses setiap baris untuk mengekstrak konten dalam rentang waktu
	for i := 1; i < len(vttContent); i++ {
		line := vttContent[i]
		
		// Memeriksa apakah baris ini adalah baris timestamp
		if strings.Contains(line, "-->") {
			// Mem-parse waktu awal
			parts := strings.Split(line, " --> ")
			if len(parts) >= 1 {
				startTimeStr := strings.TrimSpace(parts[0])
				startTime := parseVTTTime(startTimeStr)
				
				// Jika kita telah melewati jendela waktu, berhenti
				if startTime > endTime {
					break
				}
			}
		}
		
		result = append(result, line)
	}
	
	return strings.Join(result, "\n")
}

// parseVTTTime mengonversi string waktu VTT menjadi detik
func parseVTTTime(timeStr string) float64 {
	// Menangani format seperti "00:00.003" (mm:ss.mmm)
	if len(strings.Split(timeStr, ":")) == 2 {
		parts := strings.Split(timeStr, ":")
		if len(parts) == 2 {
			var minutes, seconds float64
			fmt.Sscanf(parts[0], "%f", &minutes)
			fmt.Sscanf(parts[1], "%f", &seconds)
			return minutes*60 + seconds
		}
	}
	
	// Menangani format seperti "00:00:00.003" (hh:mm:ss.mmm)
	if len(strings.Split(timeStr, ":")) == 3 {
		parts := strings.Split(timeStr, ":")
		if len(parts) == 3 {
			var hours, minutes, seconds float64
			fmt.Sscanf(parts[0], "%f", &hours)
			fmt.Sscanf(parts[1], "%f", &minutes)
			fmt.Sscanf(parts[2], "%f", &seconds)
			return hours*3600 + minutes*60 + seconds
		}
	}
	
	return 0.0
}

// GetGeminiAPIKey mendapatkan API key dari berbagai sumber secara berurutan
func GetGeminiAPIKey() string {
	// 1. Mencoba mendapatkan API key dari environment variable sistem
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey != "" {
		return apiKey
	}

	// 2. Mencoba membaca dari file .env di direktori kerja saat ini atau parent-nya
	paths := []string{".env", "../.env", "../../.env", ".env.local", "../.env.local", "../../.env.local"}
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			if apiKey = os.Getenv("GEMINI_API_KEY"); apiKey != "" {
				return apiKey
			}
		}
	}

	// 3. Jika masih tidak ditemukan, coba cari .env di direktori executable itu sendiri
	execPath, err := os.Executable()
	if err == nil {
		projectEnvPath := filepath.Join(filepath.Dir(execPath), ".env")
		if err := godotenv.Load(projectEnvPath); err == nil {
			if apiKey = os.Getenv("GEMINI_API_KEY"); apiKey != "" {
				return apiKey
			}
		}
	}

	// Jika tidak ditemukan di semua path, kembalikan string kosong
	return ""
}