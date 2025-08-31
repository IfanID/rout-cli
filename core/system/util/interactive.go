package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirm displays a prompt and waits for a yes/no response.
func Confirm(prompt string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		TypeOut(fmt.Sprintf("%s [y/n]:", prompt))
		input, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}
		input = strings.ToLower(strings.TrimSpace(input))
		if input == "y" || input == "yes" {
			return true, nil
		}
		if input == "n" || input == "no" {
			return false, nil
		}
		// If input is something else, the loop continues.
	}
}