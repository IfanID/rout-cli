package command

import (
	"fmt"
	"rout/core/system/util"
	"strings"
)

// CommandFunc defines the type for command functions.
type CommandFunc func(args []string) error

// CommandsMap stores all registered commands.
var CommandsMap = make(map[string]CommandFunc)

func init() {
	// Register commands
	CommandsMap["ls"] = handleLs
	CommandsMap["cd"] = handleCd
	CommandsMap["pwd"] = handlePwd
	CommandsMap["touch"] = handleTouch
	CommandsMap["mkdir"] = handleMkdir
	CommandsMap["rm"] = handleRm
	CommandsMap["cp"] = handleCp
	CommandsMap["mv"] = handleMv
	CommandsMap["help"] = handleHelp
	CommandsMap["clear"] = handleClear
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

	if cmdFunc, ok := CommandsMap[commandName]; ok {
		if err := cmdFunc(commandArgs); err != nil {
			util.TypeOut(fmt.Sprintf("Error: %v", err))
		}
	} else {
		util.TypeOut(fmt.Sprintf("Perintah tidak dikenal: %s", commandName))
	}
}