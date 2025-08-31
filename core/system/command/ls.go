package command

import (
	"fmt"
	"io/fs" // For fs.FileInfo
	"os"
	"strconv" // For human-readable size
	"strings"
	
)

// Ls lists the contents of a directory with options.
func Ls(path string, showAll, longFormat, humanReadable bool) error {
	if path == "" {
		path = "." // Default to current directory
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("gagal membaca direktori %s: %w", path, err)
	}

	for _, file := range files {
		if !showAll && strings.HasPrefix(file.Name(), ".") {
			continue // Skip hidden files unless showAll is true
		}

		if longFormat {
			info, err := file.Info()
			if err != nil {
				fmt.Printf("Error getting info for %s: %v\n", file.Name(), err)
				continue
			}
			fmt.Printf("%s %s %s %s %s %s %s\n",
				formatPermissions(info.Mode()),
				"1", // Placeholder for number of links
				"owner", // Placeholder for owner
				"group", // Placeholder for group
				formatSize(info.Size(), humanReadable),
				info.ModTime().Format("Jan _2 15:04"),
				file.Name(),
			)
		} else {
			fmt.Println(file.Name())
		}
	}
	return nil
}

// Helper function to format file permissions (simplified)
func formatPermissions(mode fs.FileMode) string {
	perm := "-"
	if mode.IsDir() {
		perm = "d"
	}
	// Simplified permissions for rwx
	if mode&0400 != 0 {
		perm += "r"
	} else {
		perm += "-"
	}
	if mode&0200 != 0 {
		perm += "w"
	} else {
		perm += "-"
	}
	if mode&0100 != 0 {
		perm += "x"
	} else {
		perm += "-"
	}
	// Add placeholders for group and others
	perm += "rwx" // Simplified for example
	perm += "rwx" // Simplified for example
	return perm
}

// Helper function to format file size
func formatSize(size int64, humanReadable bool) string {
	if !humanReadable {
		return strconv.FormatInt(size, 10)
	}

	const (
		_ = iota
		KB int64 = 1 << (10 * iota)
		MB
		GB
		TB
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.1fT", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.1fG", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.1fM", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.1fK", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%dB", size)
	}
}