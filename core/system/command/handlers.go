package command

import (
	"fmt"
	"rout/core/system/util"
	"strings"
	"rout/core/system/command/manajemen_file"
)

// handleLs adapts Ls to CommandFunc signature
func handleLs(args []string) error {
	var path string
	showAll := false
	longFormat := false
	humanReadable := false
	sortByTime := false
	sortBySize := false
	reverseOrder := false

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
				case 't':
					sortByTime = true
				case 'S':
					sortBySize = true
				case 'r':
					reverseOrder = true
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

	return manajemen_file.Ls(path, showAll, longFormat, humanReadable, sortByTime, sortBySize, reverseOrder)
}

// handleCd adapts Cd to CommandFunc signature
func handleCd(args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	}
	return manajemen_file.Cd(path)
}

// handlePwd adapts Pwd to CommandFunc signature
func handlePwd(args []string) error {
	dir, err := manajemen_file.Pwd()
	if err != nil {
		return err
	}
	util.TypeOut(dir)
	return nil
}

// handleTouch adapts Touch to CommandFunc signature
func handleTouch(args []string) error {
	return manajemen_file.Touch(args)
}

// handleMkdir adapts Mkdir to CommandFunc signature
func handleMkdir(args []string) error {
	createParents := false
	paths := []string{}

	for _, arg := range args {
		if arg == "-p" {
			createParents = true
		} else {
			paths = append(paths, arg)
		}
	}

	if len(paths) == 0 {
		return fmt.Errorf("mkdir: nama direktori harus disertakan")
	}

    // Standard mkdir can create multiple directories at once
    for _, path := range paths {
        err := manajemen_file.Mkdir(path, createParents)
        if err != nil {
            return err // Stop on first error
        }
    }
	return nil
}

// handleRm adapts Rm to CommandFunc signature
func handleRm(args []string) error {
	recursive := false
	force := false
	paths := []string{}

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
            // It's a flag, check for 'r' and 'f'
            if strings.Contains(arg, "r") {
                recursive = true
            }
            if strings.Contains(arg, "f") {
                force = true
            }
		} else {
			paths = append(paths, arg)
		}
	}

	if len(paths) == 0 {
		return fmt.Errorf("rm: operand file hilang")
	}

    for _, path := range paths {
        err := manajemen_file.Rm(path, recursive, force)
        if err != nil {
            return err
        }
    }
	return nil
}

// handleCp adapts Cp to CommandFunc signature
func handleCp(args []string) error {
	recursive := false
	force := false
	interactive := false
	paths := []string{}

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			if strings.Contains(arg, "r") {
				recursive = true
			}
			if strings.Contains(arg, "f") {
				force = true
			}
			if strings.Contains(arg, "i") {
				interactive = true
			}
		} else {
			paths = append(paths, arg)
		}
	}

	if len(paths) < 2 {
		if len(paths) == 1 {
			return fmt.Errorf("cp: operand file tujuan hilang setelah '%s'", paths[0])
		}
		return fmt.Errorf("cp: operand file hilang")
	}

	if len(paths) > 2 {
		return fmt.Errorf("cp: dukungan untuk banyak sumber belum diimplementasikan")
	}

	src := paths[0]
	dst := paths[1]

	// -f overrides -i
	if force {
		interactive = false
	}

	return manajemen_file.Cp(src, dst, recursive, force, interactive)
}

// handleMv adapts Mv to CommandFunc signature
func handleMv(args []string) error {
	force := false
	interactive := false
	paths := []string{}

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			if strings.Contains(arg, "f") {
				force = true
			}
			if strings.Contains(arg, "i") {
				interactive = true
			}
		} else {
			paths = append(paths, arg)
		}
	}

	if len(paths) < 2 {
		if len(paths) == 1 {
			return fmt.Errorf("mv: operand file tujuan hilang setelah '%s'", paths[0])
		}
		return fmt.Errorf("mv: operand file hilang")
	}

	if len(paths) > 2 {
		return fmt.Errorf("mv: dukungan untuk banyak sumber belum diimplementasikan")
	}

	src := paths[0]
	dst := paths[1]

	// -f overrides -i
	if force {
		interactive = false
	}

	return manajemen_file.Mv(src, dst, force, interactive)
}

// handleHelp adapts Help to CommandFunc signature
func handleHelp(args []string) error {
	Help()
	return nil
}
