package conv

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	UiWidth = 80

	// ANSI Color Codes
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[90m"
)

// PrintHeader prints a styled header box.
func PrintHeader(title string, color string) {
	// Truncate title if it's too long
	if len(title) > UiWidth-4 {
		title = title[:UiWidth-7] + "..."
	}

	padding := (UiWidth - 2 - len(title)) / 2
	rightPadding := UiWidth - 2 - len(title) - padding

	border := "╭" + strings.Repeat("─", UiWidth-2) + "╮"
	titleLine := fmt.Sprintf("│%s%s%s│", strings.Repeat(" ", padding), title, strings.Repeat(" ", rightPadding))
	bottomBorder := "╰" + strings.Repeat("─", UiWidth-2) + "╯"

	fmt.Println()
	fmt.Println(color + border + ColorReset)
	fmt.Println(color + titleLine + ColorReset)
	fmt.Println(color + bottomBorder + ColorReset)
}

// PrintLog prints a formatted log line with indentation.
func PrintLog(level int, emoji string, message string, color string) {
	indent := strings.Repeat("  ", level)
	fmt.Printf("%s[%s] %s%s%s\n", indent, emoji, color, message, ColorReset)
}

// PrintSubLog is a convenience function for sub-level logging.
func PrintSubLog(level int, emoji string, message string, color string) {
	indent := strings.Repeat("  ", level)
	fmt.Printf("%s %s %s%s%s\n", indent, emoji, color, message, ColorReset)
}

// PrintKeyValue prints a key-value pair with aligned formatting.
func PrintKeyValue(level int, key, value, color string) {
	indent := strings.Repeat("  ", level)
	keyPart := fmt.Sprintf("%s  - %-18s: ", indent, key)
	valuePrefix := strings.Repeat(" ", len(keyPart))

	maxWidth := UiWidth - len(keyPart)
	wrappedValue := WrapText(value, valuePrefix, maxWidth)

	fmt.Printf("%s%s%s%s\n", keyPart, color, wrappedValue, ColorReset)
}

// WrapText membungkus teks panjang ke beberapa baris dengan prefix tertentu.
func WrapText(text string, prefix string, maxWidth int) string {
	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return ""
	}

	var wrapped strings.Builder
	line := words[0]
	for _, word := range words[1:] {
		if len(line+" "+word) > maxWidth {
			wrapped.WriteString(line + "\n" + prefix)
			line = word
		} else {
			line += " " + word
		}
	}
	wrapped.WriteString(line)
	return wrapped.String()
}

// PrintFilePreview prints the first few lines of a file.
func PrintFilePreview(level int, title, filePath string, maxLines int, color string) {
	PrintLog(level, "📄", title, color)
	file, err := os.Open(filePath)
	if err != nil {
		PrintSubLog(level+1, "❌", fmt.Sprintf("Gagal membaca file untuk pratinjau: %v", err), ColorRed)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	previewIndent := strings.Repeat("  ", level+1)
	lineCount := 0
	for scanner.Scan() {
		if lineCount >= maxLines {
			fmt.Println(previewIndent + ColorGray + "[...]" + ColorReset)
			break
		}
		fmt.Println(previewIndent + ColorGray + scanner.Text() + ColorReset)
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		PrintSubLog(level+1, "❌", fmt.Sprintf("Error saat memindai file: %v", err), ColorRed)
	}
	fmt.Println()
}

func RunSpinner(done chan bool, message string) {
	spinnerChars := []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
	i := 0
	for {
		select {
		case <-done:
			fmt.Printf("\r%s\r", strings.Repeat(" ", len(message)+5))
			return
		default:
			fmt.Printf("\r %s[%s] %c%s ", ColorBlue, message, spinnerChars[i], ColorReset)
			i = (i + 1) % len(spinnerChars)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func ShowFzfMenu(prompt string, options []string) (string, error) {
	input := strings.Join(options, "\n")
	cmd := exec.Command("fzf", "--prompt", prompt, "--height", "15%", "--layout", "reverse", "--border")
	cmd.Stdin = strings.NewReader(input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 130 {
			if len(options) > 0 {
				return options[0], nil
			}
			return "", fmt.Errorf("fzf dibatalkan dan tidak ada pilihan default")
		}
		return "", fmt.Errorf("gagal menjalankan fzf: %w\nStderr: %s", err, stderr.String())
	}

	selected := strings.TrimSpace(stdout.String())
	if selected == "" && len(options) > 0 {
		return options[0], nil
	}
	return selected, nil
}
