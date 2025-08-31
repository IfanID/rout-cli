package core

import (
	"fmt"
	"os"
	"strings"
)

// PrintPrompt mencetak prompt CLI kustom.
// Simbol  (\uF07B) membutuhkan Nerd Font terinstal di terminal Anda.
func Prompt() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("ROut > [error] > ")
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("ROut %s [%s] > ", "\uF07B", currentDir) // Fallback to full path if home dir not found
		return
	}

	displayPath := currentDir
	if strings.HasPrefix(currentDir, homeDir) {
		displayPath = strings.Replace(currentDir, homeDir, "~", 1)
	}

	        fmt.Printf("ROut %s [%s] > ", "\uF07B", displayPath)
}
