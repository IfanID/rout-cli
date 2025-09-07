package rcli

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

// --- Fungsi Word Wrap --- 
func wrapText(text string, lineWidth int) string {
	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return ""
	}

	var wrappedText strings.Builder
	wrappedText.WriteString(words[0])
	currentLineLength := len(words[0])

	for _, word := range words[1:] {
		if currentLineLength+1+len(word) > lineWidth {
			wrappedText.WriteString("\n")
			currentLineLength = 0
		} else {
			wrappedText.WriteString(" ")
			currentLineLength++
		}
		wrappedText.WriteString(word)
		currentLineLength += len(word)
	}

	return wrappedText.String()
}

// --- Fungsi Animasi --- 
func showWaitingAnimation(wg *sync.WaitGroup, done chan bool) {
	defer wg.Done()
	ticker := time.NewTicker(120 * time.Millisecond)
	defer ticker.Stop()

	states := []string{
		".    ", " .   ", "  .  ", "   . ", "    .", "   . ", "  .  ", " .   ",
	}
	i := 0

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			frame := states[i%len(states)]
			fmt.Printf("\rAI: %s  ", frame)
			i++
		}
	}
}

// --- Fungsi Panggilan AI --- 
func callGemini(ctx context.Context, model *genai.GenerativeModel, prompt string) {
	var wg sync.WaitGroup
	done := make(chan bool)

	wg.Add(1)
	go showWaitingAnimation(&wg, done)

	systemInstruction := "Gunakan Bahasa Indonesia, campur gaya formal dan non-formal. Apapun yang terjadi, selalu gunakan Bahasa Indonesia."
	fullPrompt := []genai.Part{genai.Text(systemInstruction), genai.Text(prompt)}
	resp, err := model.GenerateContent(ctx, fullPrompt...)

	done <- true
	wg.Wait()

	fmt.Printf("\r%s\r", strings.Repeat(" ", 20))

	if err != nil {
		log.Printf("Error dari API: %v\n\n", err)
		return
	}

	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		fmt.Print("AI: ")
		for _, part := range resp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				// Terapkan word wrap di sini
				fmt.Println(wrapText(string(txt), 80))
			}
		}
	} else {
		fmt.Println("AI: Tidak ada respons dari model.")
	}
	fmt.Println()
}

// --- Perintah Utama Rcli --- 
var RcliCmd = &cobra.Command{
	Use:   "rcli",
	Short: "Memulai sesi chat interaktif dengan AI Gemini",
	Long:  `Memulai sesi chat interaktif dengan AI Gemini, menggunakan model gemini-2.5-flash.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error memuat file .env: %v", err)
		}
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			fmt.Println("GEMINI_API_KEY tidak ditemukan.")
			return
		}

		ctx := context.Background()
		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		model := client.GenerativeModel("gemini-2.5-flash")

		fmt.Println("Masuk ke mode chat interaktif (Model: gemini-2.5-flash). Ketik 'exit' atau 'keluar' untuk berhenti.")
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Anda: ")
			if !scanner.Scan() {
				break
			}
			input := scanner.Text()
			if input == "exit" || input == "keluar" {
				break
			}
			if input != "" {
				callGemini(ctx, model, input)
			}
		}
	},
}