package command

import (
	"fmt"
	"strings"
)

// CommandFunc defines the type for command functions.
type CommandFunc func(args []string) error

// commandsMap stores all registered commands.
var commandsMap = make(map[string]CommandFunc)

func init() {
	// Register commands
	commandsMap["ls"] = handleLs
	commandsMap["cd"] = handleCd
	commandsMap["pwd"] = handlePwd
}

// RegisterCommands parses the input and executes the corresponding command.
func RegisterCommands(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	commandName := parts[0]
	commandArgs := []string{}
	if len(parts) > 1 {
		commandArgs = parts[1:]
	}

	if cmdFunc, ok := commandsMap[commandName]; ok {
		if err := cmdFunc(commandArgs); err != nil {
			fmt.Println("Error:", err)
		}
	} else {
		fmt.Println("Perintah tidak dikenal:", commandName)
	}
}