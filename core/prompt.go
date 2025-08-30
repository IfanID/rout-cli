package core

import (
	"fmt"
)

// PrintPrompt mencetak prompt CLI kustom.
// Simbol  (\uF07B) membutuhkan Nerd Font terinstal di terminal Anda.
func Prompt() {
	fmt.Print("ROut \uF07B [~] > ")
}


