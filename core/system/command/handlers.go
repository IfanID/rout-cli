package command

import (
	"fmt"
	"strings"
)

// handleLs adapts Ls to CommandFunc signature
func handleLs(args []string) error {
	var path string
	showAll := false
	longFormat := false
	humanReadable := false

	// Parse flags
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			for _, char := range arg[1:] {
				switch char {
				case 'a':
					showAll = true
				case 'l':
					longFormat = true
				case 'h':
					humanReadable = true
				default:
					return fmt.Errorf("ls: opsi tidak valid -- '%c'", char)
				}
			}
		} else {
			// Assume the first non-flag argument is the path
			if path == "" {
				path = arg
			} else {
				return fmt.Errorf("ls: terlalu banyak argumen")
			}
		}
	}

	return Ls(path, showAll, longFormat, humanReadable)
}

// handleCd adapts Cd to CommandFunc signature
func handleCd(args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	}
	return Cd(path)
}

// handlePwd adapts Pwd to CommandFunc signature
func handlePwd(args []string) error {
	dir, err := Pwd()
	if err != nil {
		return err
	}
	fmt.Println(dir)
	return nil
}
