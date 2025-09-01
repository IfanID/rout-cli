package core

import (
	"fmt" // Tambahkan kembali import ini
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

// PrintPrompt mencetak prompt CLI kustom.
// Simbol  (\uF07B) membutuhkan Nerd Font terinstal di terminal Anda.

// Define a slice of colors to choose from
var colors = []color.Attribute{
	color.FgRed,
	color.FgGreen,
	color.FgYellow,
	color.FgBlue,
	color.FgMagenta,
	color.FgCyan,
	color.FgWhite,
	color.FgHiRed,
	color.FgHiGreen,
	color.FgHiYellow,
	color.FgHiBlue,
	color.FgHiMagenta,
	color.FgHiCyan,
	color.FgHiWhite,
}

// Initialize random seed once
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Helper function to get a random color
func getRandomColor() color.Attribute {
	return colors[rand.Intn(len(colors))]
}

func Prompt() string { // Mengembalikan string
	currentDir, err := os.Getwd()
	if err != nil {
		return "[error] > " // Mengembalikan string error
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to full path if home dir not found
		return fmt.Sprintf("%s %s [%s] %s %s",
			color.New(color.FgCyan).Sprint("ROut"), // Fixed to Cyan
			color.New(getRandomColor()).Sprint("\uF07B"),
			color.New(color.FgWhite).Sprint(currentDir), // Fixed to White
			color.New(color.FgCyan).Sprint("»"), // Fixed to Cyan
			color.New(getRandomColor()).Sprint(""))
	}


	displayPath := currentDir
	if strings.HasPrefix(currentDir, homeDir) {
		displayPath = strings.Replace(currentDir, homeDir, "~", 1)
	}

	// Use random colors for prompt elements
	routColor := color.New(color.FgCyan).SprintFunc() // Fixed to Cyan
	iconColor := color.New(getRandomColor()).SprintFunc()
	pathColor := color.New(color.FgWhite).SprintFunc() // Fixed to White
	promptColor := color.New(color.FgCyan).SprintFunc() // Fixed to Cyan

	return fmt.Sprintf("%s %s [%s] %s %s", // Mengembalikan string
		routColor("ROut"),
		iconColor("\uF07B"),
		pathColor(displayPath),
		color.New(color.FgCyan).Sprint("»"), // Fixed to Cyan
		promptColor(""))
}
