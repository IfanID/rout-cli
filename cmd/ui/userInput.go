package ui

import (
	"bufio"
	"os"
	"strings"

	"rout/core"
	"rout/core/system/command"
)

func HandleUserInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		core.Prompt() // Panggil fungsi prompt dari paket core

		input, _ := reader.ReadString('\n')
		trimmedInput := strings.TrimSpace(input)

		if trimmedInput == "exit" || trimmedInput == "quit" {
			core.Logout()
			break
		}

		// FIX: Abaikan baris kosong atau baris yang kemungkinan adalah item daftar dari hasil tempel.
		if trimmedInput == "" || strings.HasPrefix(trimmedInput, "*") || strings.HasPrefix(trimmedInput, "-") {
			continue
		}

		command.RegisterCommands(trimmedInput) // Panggil fungsi pendaftaran perintah
	}
}
