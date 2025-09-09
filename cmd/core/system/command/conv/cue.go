package conv

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"rout/cmd/API"
)

const (
	CueText = "Ifan.3 V2S CoreX"
)

func AddCue(filename string) error {
	tmpfile := filename + ".tmp"
	lines, err := ReadLines(filename)
	if err != nil {
		return fmt.Errorf("gagal membaca file untuk addCue: %w", err)
	}

	cueContent := CueText
	var finalTimestamp string
	var cueStart time.Duration
	var cueEnd time.Duration

	// --- Remove existing CUEs first ---
	PrintLog(2, "🧹", "Membersihkan CUE lama...", ColorGray)
	cleanLines := []string{}
	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], cueContent) {
			// Cue found, skip this line and look for its timestamp and preceding blank line to skip as well.
			if i > 0 && regexp.MustCompile(`-->`).MatchString(lines[i-1]) {
				// Timestamp is on the previous line. Remove it from cleanLines if it was added.
				if len(cleanLines) > 0 && cleanLines[len(cleanLines)-1] == lines[i-1] {
					cleanLines = cleanLines[:len(cleanLines)-1]
					// Check for a blank line before the timestamp
					if len(cleanLines) > 0 && strings.TrimSpace(cleanLines[len(cleanLines)-1]) == "" {
						cleanLines = cleanLines[:len(cleanLines)-1]
					}
				}
			}
			continue // Skip the CUE content line
		}
		cleanLines = append(cleanLines, lines[i])
	}
	lines = cleanLines

	// --- Determine timestamp (Smart or Fallback) ---
	useOverwriteMethod := false
	apiKey := API.GetGeminiAPIKey()
	firstDialogueStart, firstDialogueFound := FindFirstDialogueStartTime(lines)

	if apiKey == "" {
		PrintLog(2, "⚠️", "GEMINI_API_KEY tidak ditemukan, menggunakan metode timpa.", ColorYellow)
		useOverwriteMethod = true
	} else {
		done := make(chan bool)
		go RunSpinner(done, "Menganalisis penempatan CUE cerdas...")
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		vttHeader := API.ExtractVTTHeader(lines, 10*time.Second)
		geminiTimestamp, reason, err := API.AnalyzeVTTForCUE(ctx, apiKey, vttHeader)
		done <- true
		close(done)

		if err != nil {
			PrintLog(2, "⚠️", "Analisis cerdas gagal, menggunakan metode timpa.", ColorYellow)
			useOverwriteMethod = true
		} else {
			PrintLog(2, "ℹ️", "Analisis Gemini Selesai. Hasil:", ColorBlue)
			PrintKeyValue(2, "Timestamp mentah", geminiTimestamp, ColorCyan)
			PrintKeyValue(2, "Alasan", reason, ColorCyan)

			tsParts := strings.Split(geminiTimestamp, " --> ")
			if len(tsParts) != 2 {
				useOverwriteMethod = true
			} else {
				var err1, err2 error
				cueStart, err1 = ParseVTTDuration(tsParts[0])
				cueEnd, err2 = ParseVTTDuration(tsParts[1])
				if err1 != nil || err2 != nil {
					useOverwriteMethod = true
				} else {
					// Validate Gemini's timestamp based on the 4-second rule
					isIdeal := true
					if firstDialogueFound {
						if firstDialogueStart > 4*time.Second && firstDialogueStart < 60*time.Second {
							// GAP RULE: CUE should NOT overlap
							if cueEnd >= firstDialogueStart {
										PrintLog(2, "⚠️", "Timestamp Gemini tidak ideal (melanggar aturan jeda 4s-1m).", ColorYellow)
										isIdeal = false
									}
						} else {
						// OVERWRITE RULE: CUE SHOULD overlap (for <4s and >1m starts)
							if cueEnd < firstDialogueStart {
										 PrintLog(2, "⚠️", "Timestamp Gemini tidak ideal (tidak menimpa dialog pertama).", ColorYellow)
										isIdeal = false
							}
						}
					}
					if !isIdeal {
						useOverwriteMethod = true
					}
				}
			}
		}
	}

	if useOverwriteMethod {
		PrintLog(2, "🔩", "Menggunakan metode timpa untuk penempatan CUE.", ColorBlue)
		cueStart = 3 * time.Millisecond
		if firstDialogueFound && firstDialogueStart > 4*time.Second && firstDialogueStart < 60*time.Second {
			PrintLog(3, "💡", "Dialog pertama antara 4s-1m. CUE ditempatkan di jeda.", ColorCyan)
			cueEnd = 2*time.Second + 500*time.Millisecond
		} else {
			var subtitleStartTimes []time.Duration
			// Use strings.Contains for robustness as regex was failing
			for i, line := range lines {
				if strings.Contains(line, "-->") {
					if i > 2 {
						parts := strings.Split(line, " --> ")
						if subStart, err := ParseVTTDuration(parts[0]); err == nil {
							subtitleStartTimes = append(subtitleStartTimes, subStart)
						}
					}
				}
			}
			defaultCueDuration := 2*time.Second + 500*time.Millisecond
			if len(subtitleStartTimes) >= 2 {
				cueEnd = subtitleStartTimes[1] - time.Millisecond
				PrintLog(3, "💡", fmt.Sprintf("Waktu CUE disesuaikan otomatis agar berhenti sebelum subtitle kedua (%s).", FormatVTTDuration(cueEnd)), ColorCyan)
			} else if len(subtitleStartTimes) == 1 {
				cueEnd = subtitleStartTimes[0] - time.Millisecond
				PrintLog(3, "💡", "Hanya ada satu subtitle, CUE disesuaikan agar tidak tumpang tindih.", ColorCyan)
			} else {
				cueEnd = defaultCueDuration
				PrintLog(3, "💡", "Tidak ada subtitle, menggunakan durasi CUE default.", ColorCyan)
			}
			if cueEnd <= cueStart {
				cueEnd = defaultCueDuration
			}
		}
	}

	finalTimestamp = FormatVTTDuration(cueStart) + " --> " + FormatVTTDuration(cueEnd)
	PrintKeyValue(2, "Timestamp Final", finalTimestamp, ColorGreen)

	// --- Deletion Logic ---
	PrintLog(2, "🗑️", "Menghapus subtitle yang tumpang tindih...", ColorGray)
	tempLines := []string{}
	i := 0
	for i < len(lines) {
		line := lines[i]
		// Use strings.Contains for robustness as regex was failing
		if strings.Contains(line, "-->") {
			parts := strings.Split(line, " --> ")
			subStart, err := ParseVTTDuration(parts[0])
			if err == nil && subStart < cueEnd {
				blockEnd := i
				for j := i + 1; j < len(lines); j++ {
					if strings.TrimSpace(lines[j]) == "" {
						blockEnd = j
						break
					}
					blockEnd = j
				}
				for k := i + 1; k <= blockEnd; k++ {
					if strings.TrimSpace(lines[k]) != "" {
						PrintKeyValue(3, "Menghapus", lines[k], ColorRed)
					}
				}
				i = blockEnd + 1
				continue
			}
		}
		tempLines = append(tempLines, line)
		i++
	}
	lines = tempLines

	// --- Insert the new CUE block ---
	PrintLog(2, "➕", "Menambahkan CUE baru di awal file...", ColorBlue)
	cueBlock := []string{"", finalTimestamp, cueContent, ""}
	headerEndIndex := 0
	for i, line := range lines {
		if i > 0 && strings.TrimSpace(line) == "" {
			headerEndIndex = i
			break
		}
	}
	if headerEndIndex == 0 && len(lines) > 0 {
		headerEndIndex = 1
	} else if headerEndIndex > 0 {
		headerEndIndex++
	}
	lines = InsertStringSlice(lines, headerEndIndex, cueBlock...)

	fmt.Println()
	PrintLog(2, "📄", "Pratinjau Hasil (15 baris pertama):", ColorBlue)
    	previewIndent := strings.Repeat("  ", 3)
	for i := 0; i < 15 && i < len(lines); i++ {
		fmt.Println(previewIndent + ColorGray + lines[i] + ColorReset)
	}
	if len(lines) > 15 {
		fmt.Println(previewIndent + ColorGray + "[...]")
	}
	fmt.Println()

	outFile, err := os.Create(tmpfile)
	if err != nil {
		return fmt.Errorf("gagal membuat file sementara: %w", err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("gagal menulis ke file sementara: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("gagal flush writer: %w", err)
	}
	outFile.Close()

	return os.Rename(tmpfile, filename)
}

func FindFirstDialogueStartTime(lines []string) (time.Duration, bool) {
	// Use strings.Contains for robustness as regex was failing
	for i, line := range lines {
		if strings.Contains(line, "-->") {
			if i > 0 && strings.Contains(lines[i-1], CueText) {
				continue
			}
			parts := strings.Split(line, " --> ")
			startTime, err := ParseVTTDuration(parts[0])
			if err == nil {
				return startTime, true
			}
		}
	}
	return 0, false
}

func VerifyCue(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), CueText) {
			return true
		}
	}
	return false
}
