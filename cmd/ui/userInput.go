package ui

import (
	"strings"

	"github.com/chzyer/readline" // Import readline

	"rout/core"
	"rout/core/system/command"
)

func HandleUserInput() {
	// Buat completer dasar
	completer := readline.NewPrefixCompleter(
		readline.PcItem("ls"),
		readline.PcItem("cd"),
		readline.PcItem("pwd"),
		readline.PcItem("touch"),
		readline.PcItem("mkdir"),
		readline.PcItem("rm"),
		readline.PcItem("cp"),		readline.PcItem("mv"),
		readline.PcItem("help"),
		readline.PcItem("clear"),
		readline.PcItem("exit"),
		readline.PcItem("quit"),
	)

	// Inisialisasi readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       core.Prompt(), // Gunakan fungsi Prompt dari core
		AutoComplete: completer,
		HistoryFile:  "/tmp/rout_readline_history.tmp", // File history sementara
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil { // EOF or Ctrl+D
			break
		}

		trimmedInput := strings.TrimSpace(line)

		if trimmedInput == "exit" || trimmedInput == "quit" {
			core.Logout()
			break
		}

		// Abaikan baris kosong atau baris yang kemungkinan adalah item daftar dari hasil tempel.
		if trimmedInput == "" || strings.HasPrefix(trimmedInput, "*") || strings.HasPrefix(trimmedInput, "-") {
			continue
		}

		command.RegisterCommands(trimmedInput) // Panggil fungsi pendaftaran perintah
		rl.SetPrompt(core.Prompt()) // Perbarui prompt secara dinamis
	}
}
