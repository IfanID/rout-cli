package command

import (
	"os"
	"os/exec"
)

func handleClear(args []string) error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
