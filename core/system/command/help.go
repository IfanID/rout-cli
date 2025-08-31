package command

import (
	"fmt"
	"rout/core/system/util"
	"sort"
)

// Help displays a list of available commands.
func Help() {
	util.TypeOut("Daftar perintah yang tersedia:")

	// Get command names and sort them alphabetically
	var commandNames []string
	for name := range CommandsMap { // Accessing the exported map
		commandNames = append(commandNames, name)
	}
	sort.Strings(commandNames)

	for _, name := range commandNames {
		// For now, just print the name. We can add descriptions later.
		util.TypeOut(fmt.Sprintf("  %s", name))
	}
	util.TypeOut("\nKetik 'exit' atau 'quit' untuk keluar.")
}